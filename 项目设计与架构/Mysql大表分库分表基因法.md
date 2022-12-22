# 引出题目- 关于数据库设计每天新增100W数据
背景: 订单表,每天新增100W数据, 三个字段 user_id, order_id, amount

要求: 1- 根据订单id查询到订单信息,

2- 根据用户id查询到所有订单的信息

设计方案: Mysql单表一般存储500W-1000W左右比较合适

# 索引表法
比如使用, 10个库+每个库100张表，平均每张表每天会产生100w的数据，这样每张表每个月就会产生3000w的数据，

在这个表中我们可以至保留一个月的数据，其余的数据归档至数据仓库，比如我们需要olap场景，就归档至`clickhouse、es`等；

## 对`order_id`进行分片键设计
hash(order_id) % 10 得到库, hash(order_id) % 100得到表

为了能通过`user_id`查询`order_id`，需要额外建一张`user_id->order_id`的索引表。`这个索引表也需要根据user_id作为分片键进行相应的分库分表`

这样我们通过order_id能够直接定位数据库的表进行查询，通过user_id再查order_id, 可以查到用户的所有订单信息,满足需求；

反思：上面的做法可以实现需求，但是通过user_id查询订单时，需要多进行一次查询，效率降低了一倍；并且索引表也需要进行分库分表，当然索引也可以考虑其他存储介质，如Hbase，或者增加缓存来提高索引效率；`如果需要某个用户订单列表的话，还需要在应用层做数据整合，很麻烦`

# 基因法
考虑把同一个userI-d的所有订单数据落到同一个库同一个表中

思想:

先对`user_id进行取余`，落库，然后在生成`order_id`时就不能随便生成了，需要从`user_id中提取基因`，在`生成order_id`的时候，`把这个基因放到order_id的生成过程中`，这样生成出来的`order_id通过取余等运算就能得到和user_id一致的结果了`

## 原理:
现在我们把焦点集中在了“基因”这个点上，我们先来看看一个数a对另外一个数b(数b为2^n)进行取余时，其实本质上最后的结果就是a这个数二进制的最后（n+1）位，举个例子：`9%4 = 1（1001 % 100 = 001）/ 10 % 4 = 2 （1010 % 100 = 010）`，那么我们在生成订单id 的时候，只要把order_id二进制的最后（n+1）位的二进制数设置为user_id的最后（n+1）位，那么我们对user_id/order_id取余都能得到相同的结果了。（原理：比n+1位高的值，都是b数的倍数，取余时直接归零，所以取余就是取二进制最后n+1位）


了解原理后，我们只需要重新合理设计分库分表的数量，让其都是 2^n，我们重设 16个库每个库64张表 ；

`hash(user_id）% 16 `定位库的位置， `hash(user_id）% 64` 定位表的位置

生成order_id ，对user_id提取一个基因 % 64 也就是二进制最后 7位，把这个最后7位二进制也作为order_id的二进制最后7位，这样就能保证order_id的路由结果与user_id完全一致；

**通过基因法，不管是通过order_id查询数据，还是通过user_id查询数据，都能准确定位到具体的表，效率高；**

# `golang`实现
## 这里生成id可以使用`uuid`或者`snowflake`
```golang

package main

import (
	"fmt"
	"github.com/satori/go.uuid"
	"math/big"
	"strconv"
	"strings"
	"github.com/bwmarrin/snowflake"
	"time"
)
// 可以使用uuid
func generUUID() {
	// 生成uuid -> 转为 整型
	id := uuid.NewV4()
	fmt.Println("ids: ", id.String())
	fmt.Println("ids: ", len(id.String()))

	var i big.Int
	i.SetString(strings.Replace(id.String(), "-", "", 4), 16)

	fmt.Println(i.String())
	fmt.Println(len(i.String()))
}

func main() {
	// 使用雪花算法
	node , err := snowflake.NewNode(64)	
	if err != nil {
		fmt.Println("failed to create node")
		return
	}
	// 基因法
	for i := 0; i <1000; i++ {
		// 生成userId 和 提取基因
		id := node.Generate()
		userId, _ := strconv.Atoi(id.String())
		userIdGene := fetchGene(userId, 64)
		fmt.Println("userIdGene: ", userIdGene)

		orderId := node.Generate().String()
		oid, _ := strconv.Atoi(orderId)
		fmt.Println("oid: ", oid)
		orderIdGene := generateWithGene(oid, userIdGene)

		// 余数一致
		fmt.Println("userId: ", userId, ",余数：", userId % 64)
		fmt.Println("orderId: ", orderIdGene, ",余数：", orderIdGene % 64)
		fmt.Println("====================")
	}
}

func fetchGene(id int, index int) string {
	i, _ := strconv.Atoi(converToBianry(id % index))
	return fmt.Sprintf("%07d", i)
}

func generateWithGene(id int, binarySuffix string) int64  {
	if id == 0 {return 0}
	s := converToBianry(id)
	substring := s[0 : len(s) - len(binarySuffix) + 1]	//  截取位置
	newBinaryString := substring + binarySuffix
	fmt.Println("newBinaryString: ", newBinaryString)
	res, _ := strconv.ParseInt(newBinaryString, 2, 64)
	return res
}


func converToBianry(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}
```
## 这里的雪花算法问题
并发低的情况下, snowflake作为分库分表key的话, 基本求余都会集中在前面几个, 这是因为12bits位的timestamp每次**毫秒都会重新从0开始**

解决办法,需要加源码, 把timestamp 起始改为 `random.NextInt(64)` 

参考: https://www.cnblogs.com/matengfei123/p/15872114.html