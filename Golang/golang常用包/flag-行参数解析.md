
- [`flag`包](#flag包)
- [`os.Args`](#osargs)
- [`flag`包使用](#flag包使用)
	- [`flag`参数类型](#flag参数类型)
	- [定义命令行`flag`参数](#定义命令行flag参数)
		- [`flag.Type()`](#flagtype)
		- [`flag.TypeVar()`](#flagtypevar)
	- [`flag.Parse()`--解析](#flagparse--解析)
	- [`flag`其他函数](#flag其他函数)
- [完整实例](#完整实例)

# `flag`包

内置`flag`包实现了命令行参数解析

# `os.Args`

简单获取命令行参数

```golang
package main

import (
	"fmt"
	"os"
)

//os.Args demo
func main() {
    //os.Args是一个[]string 字符串切片
    if len(os.Args) > 0 {
        for index, arg := range os.Args {
            fmt.Println("args[%d] = %v\n", index, arg)
        }
    }
}
```

编译  `go build -o "args_demo"` ,执行

```shell
$ ./args_demo a b c d
args[0]=./args_demo
args[1]=a
args[2]=b
args[3]=c
args[4]=d
```

它的第一个元素是执行文件名称。

# `flag`包使用

本文介绍了flag包的常用函数和基本用法，更详细的内容请查看[官方文档](https://studygolang.com/pkgdoc)。

## `flag`参数类型

`flag`包支持命令行参数类型有`bool、int、int64、uint、uint64、float float64、string、duration`

| `flag`参数   | 有效值                                                                                                       |
| ------------ | ------------------------------------------------------------------------------------------------------------ |
| 字符串`flag` | 合法字符串                                                                                                   |
| 整数         | 1234, 0664, 0x1234等,也可以是负数                                                                            |
| 浮点数       | 合法浮点数                                                                                                   |
| `bool`       | 1, 0, t, f, T, F, true, false, TRUE, FALSE, True, False。                                                    |
| 时间段       | 任何合法的时间段字符串。如`”300ms”、”-1.5h”、”2h45m”`。合法的单位有`”ns”、”us” /“µs”、”ms”、”s”、”m”、”h”`。 |

## 定义命令行`flag`参数

### `flag.Type()`

```golang
flag.Type(flag名, 默认值, 帮助信息) *Type

name := flag.String("name", "张三", "姓名")
age := flag.Int("age", 18, "年龄")
married := flag.Bool("married", false, "婚否")
delay := flag.Duration("d", 0, "时间间隔")

// 此时name、age、married、delay均为对应类型的指针。
```

### `flag.TypeVar()`

```golang
flag.TypeVar(Type指针, flag名, 默认值, 帮助信息)

var name string
var age int
var married bool
var delay time.Duration
flag.StringVar(&name, "name", "张三", "姓名")
flag.IntVar(&age, "age", 18, "年龄")
flag.BoolVar(&married, "married", false, "婚否")
flag.DurationVar(&delay, "d", 0, "时间间隔")
```

## `flag.Parse()`--解析

通过上面两种方法定义好命令行`flag`参数后,需要使用`flag.Parse()`对参数进行解析

支持命令行参数格式如下:

*   `-flag xxx `
*   `--flag xxx`
*   `-flag=xxx`
*   `--flag=xxx`

其中布尔类型,必须使用`=`等号形式指定

## `flag`其他函数

```golang
flag.Args()  ////返回命令行参数后的其他参数，以[]string类型
flag.NArg()  //返回命令行参数后的其他参数个数
flag.NFlag() //返回使用的命令行参数个数
```

# 完整实例

```golang
func main() {
	//定义命令行参数方式1
	var name string
	var age int
	var married bool
	var delay time.Duration
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")
	flag.BoolVar(&married, "married", false, "婚否")
	flag.DurationVar(&delay, "d", 0, "延迟的时间间隔")

	//解析命令行参数
	flag.Parse()
	fmt.Println(name, age, married, delay)
	//返回命令行参数后的其他参数
	fmt.Println(flag.Args())
	//返回命令行参数后的其他参数个数
	fmt.Println(flag.NArg())
	//返回使用的命令行参数个数
	fmt.Println(flag.NFlag())
	
//命令行参数使用提示
$ ./flag_demo -help
Usage of ./flag_demo:
  -age int
        年龄 (default 18)
  -d duration
        时间间隔
  -married
        婚否
  -name string
        姓名 (default "张三")
        
//正常使用命令行flag参数：
$ ./flag_demo -name 沙河娜扎 --age 28 -married=false -d=1h30m
沙河娜扎 28 false 1h30m0s
[]
0
4

//使用非flag命令行参数：
$ ./flag_demo a b c
张三 18 false 0s
[a b c]
3
0
```

