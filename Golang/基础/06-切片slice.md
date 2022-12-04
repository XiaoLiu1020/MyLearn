- [切片本质](#切片本质)
		- [切片无法直接比较](#切片无法直接比较)
		- [切片的赋值拷贝--拷贝前后两个变量共享底层数组](#切片的赋值拷贝--拷贝前后两个变量共享底层数组)
	- [`append()` 方法为切片添加元素](#append-方法为切片添加元素)
		- [切片扩容源码](#切片扩容源码)
	- [`copy()` 复制切片](#copy-复制切片)
	- [删除切片元素方法](#删除切片元素方法)


## `slice` 切片
使用数组 `array` 引起问题:

1. 数组长度固定,并且数组长度属于类型一部分
2. 只能接受一致的类型,比如一个函数, 传入参数定义为`[3]int`类型, 其他的都不能支持了
3. 数组中元素满了,不能继续往数组`a`中添加

### 概述
`slice` 是一个拥有**相同类型的可变长度的序列**, 是基于`array`数组类型的一层封装。

是一个引用类型, 内部结构包含 `地址`、`长度`和 `容量`。 经常用于快速操作一块数据集合。

#### 切片定义
```
var name []Type

// name:表示变量名
// Type: 切片中元素类型

// 数组定义
var name [...]Type{array}   //相比，数组[]必须有东西，slice切片多了{}

func main() {
	// 声明切片类型
	var a []string              //声明一个字符串切片
	var b = []int{}             //声明一个整型切片并初始化
	var c = []bool{false, true} //声明一个布尔切片并初始化
	var d = []bool{false, true} //声明一个布尔切片并初始化
	fmt.Println(a)              //[]
	fmt.Println(b)              //[]
	fmt.Println(c)              //[false true]
	fmt.Println(a == nil)       //true
	fmt.Println(b == nil)       //false
	fmt.Println(c == nil)       //false
	// fmt.Println(c == d)   //切片是引用类型，不支持直接比较，只能和nil比较
}
```

切片有自己的长度和容量, `len()` 求长度, `cap()` 求容量

#### 可以对数组再切片
**切片本身底层就是个数组**
```
func main() {
    //基于数组定义切片
    a := [5]int{55, 56, 57, 58, 59}
    b := a[1:4]                         //基于数组a创建切片
    fmt.Println(b)                      // [56, 57, 58]
    fmt.Printf("type of b:%T\n", b)     //type of b:[]int
}

//使用方式与Python无异
```

### 可以使用 `make()` 函数构造切片
**动态创建一个切片**,使用内置的`make()`函数
```
make([]Type, size, cap)
```
* `Type`：切片元素类型
* `size`： 切片种元素数量
* `cap`：切片的容量

```
func main(){
    a := make([]int, 2, 10)
    fmt.Println(a)          //[0,0]  值
    fmt.Println(len(a))     //2    长度     
    fmt.Println(cap(a))     //10  返回该切片容量
}
```

# 切片本质
切片就是**对底层数组封装**，包含三个信息：**底层数组指针**，**切片长度(len)**,**切片容量(cap)**

现有数组：`a := [8]int{0,1,2,3,4,5,6,7}`, 切片 `s1 := a[:5]`,示意图如下：

![切片](https://www.liwenzhou.com/images/Go/slice/slice_01.png)

切片`s2 := a[3:6]`，示意图如下：
![切片2](https://www.liwenzhou.com/images/Go/slice/slice_02.png)

### 切片无法直接比较
切片之间**不能使用**`==`操作符判断切片元素是否全部相等；

唯一合法操作是和 `nil` 比较，`nil`值的切片没有底层数组，长度，容量都为`0`

```{golang}
var s1 []int         //len(s1)=0;cap(s1)=0;s1==nil
s2 := []int{}        //len(s2)=0;cap(s2)=0;s2!=nil
s3 := make([]int, 0) //len(s3)=0;cap(s3)=0;s3!=nil
```
**所以判断一个切片是否是空的，要用** `len(s) == 0`判断，**不应该用** `s == nil`。

### 切片的赋值拷贝--拷贝前后两个变量共享底层数组
类似python，对一个切片的修改会影响另一个切片的内容，这点需要注意
```
func main() {
	s1 := make([]int, 3) //[0 0 0]
	s2 := s1             //将s1直接赋值给s2，s1和s2共用一个底层数组
	s2[0] = 100
	fmt.Println(s1) //[100 0 0]
	fmt.Println(s2) //[100 0 0]
}
```

## `append()` 方法为切片添加元素
`append()`可以为切片动态添加元素，每个切片都指向一个底层数组，这个数组能容纳一定容量元素。**当底层数组不能容纳新增元素时，切片就会自动按照策略进行 扩容，此时该切片指向底层数组就会更换**。`扩容` 操作经常发生在调用 `append()`函数时；
```
func main() {
	//append()添加元素和切片扩容
	var numSlice []int
	for i := 0; i < 10; i++ {
		numSlice = append(numSlice, i)
		fmt.Printf("%v  len:%d  cap:%d  ptr:%p\n", numSlice, len(numSlice), cap(numSlice), numSlice)
	}
}

// 输出
[0]  len:1  cap:1  ptr:0xc0000a8000
[0 1]  len:2  cap:2  ptr:0xc0000a8040
[0 1 2]  len:3  cap:4  ptr:0xc0000b2020
[0 1 2 3]  len:4  cap:4  ptr:0xc0000b2020
[0 1 2 3 4]  len:5  cap:8  ptr:0xc0000b6000
[0 1 2 3 4 5]  len:6  cap:8  ptr:0xc0000b6000
[0 1 2 3 4 5 6]  len:7  cap:8  ptr:0xc0000b6000
[0 1 2 3 4 5 6 7]  len:8  cap:8  ptr:0xc0000b6000
[0 1 2 3 4 5 6 7 8]  len:9  cap:16  ptr:0xc0000b8000
[0 1 2 3 4 5 6 7 8 9]  len:10  cap:16  ptr:0xc0000b8000
```
从结果看出：
1. `append()`将元素追加到最后
2. 扩容规则为之前容量的**两倍**

`append()`也支持追加多个：
```
var citySlice []string
// 追加一个元素
citySlice = append(citySlice, "北京")
// 追加多个元素
citySlice = append(citySlice, "上海", "广州", "深圳")
// 追加切片
a := []string{"成都", "重庆"}
citySlice = append(citySlice, a...)         //追加的是切片，需要加...
fmt.Println(citySlice) //[北京 上海 广州 深圳 成都 重庆]
```

### 切片扩容源码
切片扩容，还会根据**切片中元素类型不同而进行不同处理**

可以通过查看 `$GOROOT/src/runtime/slice.go` 源码，其中扩容相关代码如下：
```
newcap := old.cap
doublecap := newcap + newcap    //申请两倍容量
if cap > doublecap {            //如果新的还比不上原来的，就使用原来的old.cap
	newcap = cap
} else {
	if old.len < 1024 {         //如果old切片长度小于1024，新的就等于原来的两倍
		newcap = doublecap
	} else {
		// Check 0 < newcap to detect overflow
		// and prevent an infinite loop.
		for 0 < newcap && newcap < cap {    //如果大于等于1024，容量比以前增加1/4，直到最终新容量大于旧容量
			newcap += newcap / 4
		}
		// Set newcap to the requested cap when
		// the newcap calculation overflowed.
		if newcap <= 0 {            //如果最终容量(cap)计算值溢出，则最终容量(cap)就是新申请的容量(cap)
			newcap = cap
		}
	}
}
```

## `copy()` 复制切片
`copy()` 函数可以将切片数据复制到**另外一个切片空间中**，指向不是同一块内存地址了；
```{golang}
copy(destSlice, srcSlice, []Type)

```
* `srcSlice` 数据来源切片
* `destSlice` 目标切片
```
func main() {
	// copy()复制切片
	a := []int{1, 2, 3, 4, 5}
	c := make([]int, 5, 5)
	copy(c, a)     //使用copy()函数将切片a中的元素复制到切片c
	fmt.Println(a) //[1 2 3 4 5]
	fmt.Println(c) //[1 2 3 4 5]
	c[0] = 1000
	fmt.Println(a) //[1 2 3 4 5]
	fmt.Println(c) //[1000 2 3 4 5] // c指向的是另外一个切片空间
}
```

## 删除切片元素方法
go语言没有专用方法，可以使用切片特性删除：
```
func main() {
	// 从切片中删除元素
	a := []int{30, 31, 32, 33, 34, 35, 36, 37}
	// 要删除索引为2的元素
	a = append(a[:2], a[3:]...)             //这里有... 是因为追加的是一个切片
	fmt.Println(a) //[30 31 33 34 35 36 37]
}
```
删除索引为 `index`的元素, 方法：`a = append(a[:index], a[index+1:]...)