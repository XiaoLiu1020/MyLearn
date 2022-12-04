- [`func` 函数](#func-函数)
	- [函数定义:](#函数定义)
	- [函数调用：](#函数调用)
	- [参数](#参数)
		- [类型简写](#类型简写)
		- [可变参数--本质是切片实现的](#可变参数--本质是切片实现的)
	- [返回值](#返回值)
		- [多返回值](#多返回值)
		- [返回值命名](#返回值命名)
- [变量作用域](#变量作用域)
	- [全局变量](#全局变量)
	- [常用在`if，for循环，switch语句`上使用定义变量的方式](#常用在iffor循环switch语句上使用定义变量的方式)
- [函数类型与变量](#函数类型与变量)
	- [定义函数类型](#定义函数类型)
	- [函数类型变量](#函数类型变量)
		- [函数也可以作为参数](#函数也可以作为参数)
		- [函数作为返回值](#函数作为返回值)
	- [匿名函数和闭包](#匿名函数和闭包)
		- [匿名函数](#匿名函数)
		- [回调函数`(sort.SliceStable)`](#回调函数sortslicestable)
		- [闭包](#闭包)
	- [`defer` 语句](#defer-语句)
		- [`defer`执行时机](#defer执行时机)
		- [`defer` 面试题](#defer-面试题)
- [内置函数](#内置函数)
	- [`panic/recover`](#panicrecover)
	- [在宕机时触发延迟执行语句](#在宕机时触发延迟执行语句)
- [面试题](#面试题)


`Go`语言支持函数，匿名函数和闭包
# `func` 函数
## 函数定义:
```
func name(args)(return) {
    func_body
}
```
* `name` 函数名，第一个不能为数字，同一个包，函数名不能重名
* `args` 参数，由参数变量和参数变量类型组成，多个参数之间使用 `,` 分隔
* `return` 返回值，由返回值变量和其变量类型组成，也可以**只写返回值类型**，**多个返回值必须用()包裹**
* `func_body` 函数体
```
func intSum(x int, y int) int {
    return x + y
}
```
**函数参数和返回值都是可选的**

## 函数调用：
可以通过 `函数名()`方式调用
```
func main() {
    sayHello()
    ret := intSum(10, 20)       //不一定需要接收其返回值
    fmt.Println(ret)
}
```
## 参数
### 类型简写
```
func intSum(x, y int) int {      //类型一样，最后写上就行
    return x + y
}
```
### 可变参数--本质是切片实现的
可变参数通过在参数名后加 `...`来表示，**可变参数通常作为函数最后一个参数**
```
func intSum(x ...int) int {
    fmt.Println(x)          // x是一个切片
    sum := 
    for _, value := range x {
        sum = sum + v
    }
    return sum
}

//调用
ret1 := intSum2()
ret2 := intSum2(10)
ret3 := intSum2(10, 20)
ret4 := intSum2(10, 20, 30)
fmt.Println(ret1, ret2, ret3, ret4) //0 10 30 60
```
**固定参数搭配可变参数时，可变参数放到最后
```
func intSum3(x int ,y ...int) int {
    fmt.Println(x, y)
    sum := x
    for _, value := range y {
        sum = sum + value
    }
    return sum
}

//调用
ret5 := intSum3(100)
ret6 := intSum3(100, 10)
ret7 := intSum3(100, 10, 20)
ret8 := intSum3(100, 10, 20, 30)
fmt.Println(ret5, ret6, ret7, ret8) //100 110 130 160
```
## 返回值
### 多返回值
```
func calc(x, y int) (int, int) {
    sum := x + y
    sub := x - y
    return sum, sub
}
```
### 返回值命名
```
func calc(x, y int) (sum, sub int){
    sum := x + y
    sub := x - y
    return
}
```

# 变量作用域
## 全局变量
**全局变量定义在函数外部，在程序整个运行周期内都有效，函数中可以访问到全局变量**
```
package main

import "fmt"

// 定义全局变量

var num int64 = 10

func testGlobalVar() {
    fmt.Printf("num=%d\n", num) //函数访问全局变量
}

func main() {
    testGlobalVar()             //num=10
}
```

**局部变量无法使用在全局上,局部变量优先级别比全局变量高**

## 常用在`if，for循环，switch语句`上使用定义变量的方式
```
func testLocalVar2(x, y int) {
    fmt.Println(x, y)       // 函数参数也只在本函数生效
    if x > 0 {
        z := 100        // 变量z 只在if语句里生效
        fmt.Println(z)
    }
}
```

# 函数类型与变量
## 定义函数类型
`type` 关键字定义函数类型
```
type calculation func(int, int) int
```
上面定义了一个 `calculation` 类型，属于函数类型，这种函数接收两个`int`参数和返回一个 `int`;

```
func add(x, y int) int {
	return x + y
}

func sub(x, y int) int {
	return x - y
}

// 满足 calculation条件 两个参数，返回一个数
add 和 sub 都能复制给 calculation类型的变量
var c calculation

c = add

```
## 函数类型变量
我们可以声明函数类型的变量并且喂该变量赋值：
```
func main() {
	var c calculation               // 声明一个calculation类型的变量c
	c = add                         // 把add赋值给c 
	
	fmt.Printf("type of c:%T\n", c) // type of c:main.calculation
	fmt.Println(c(1, 2))            // 像调用add一样调用c

	f := add                        // 将函数add赋值给变量f
	fmt.Printf("type of f:%T\n", f) // type of f:func(int, int) int
	fmt.Println(f(10, 20))          // 像调用add一样调用f
}
```
### 函数也可以作为参数
```
func add(x, y int) int {
	return x + y
}
func calc(x, y int, op func(int, int) int) int { //第三个参数为函数
	return op(x, y)
}
func main() {
	ret2 := calc(10, 20, add)
	fmt.Println(ret2) //30
}
```
或者不适用calculation类型
```
package main

import "fmt"

func added(msg string, a func(a, b int) int) {
    fmt.Println(msg, ":", a(33,44))
}

func main() {
    //函数内部不能嵌套命名函数
    //所以main()中只能定义匿名函数
    f := func(a, b int) int {
        return a + b
    }
    added("a+b", f)
}
```

### 函数作为返回值
```
func do(s string) (func(int, int) int, error) {
	switch s {
	case "+":
		return add, nil
	case "-":
		return sub, nil
	default:
		err := errors.New("无法识别的操作符")
		return nil, err
	}
}
```
或者
```
func added() func(a, b int) int {
    //使用匿名函数返回
    f := func(a, b int) int {
        return a + b
    }
    return f
}
func main() {
    m := added()
    fmt.Println(m(33, 44))
}
```

## 匿名函数和闭包
### 匿名函数
```
func(参数)(返回值){         // 相比正常函数，少了函数名部分
    函数体
}
```
匿名函数没有函数名成，没办法像普通函数那样调用，所以**匿名函数需要保存到某个变量中或者立即执行**
```
func main() {
    //将匿名函数保存到变量中
    add := func(x, y int) {
        fmt.Println(x + y)
    }
    add(10, 20) //通过变量调用匿名函数
    
    //立即执行调用
    func(x, y int) {
        fmt.Println(x + y)
    }(10, 20)
}
```
### 回调函数`(sort.SliceStable)`
```
package main

import (
    "fmt"
    "sort"
    )
    
func main() {
    s1 := []int{112, 22, 52, 32, 12}
    // 定义排序函数
    less := func(i, j int) bool {
        //降序排序
        return s1[i] > s1[j]
    }
    // 使用回调函数
    sort.SliceStable(s1, less)
    fmt.Println(s1)
}
```
或者按照字符串长度比较
```
func main() {
    s1 :=[]string{"hello", "世界", "gaoxiao"}
    sort.SliceStable(s1, func(i, j int) bool {
        //按字节大小顺序降序排序
        return len([]rune(s1[i])) > len([]rune(s1[j]))
    })
    fmt.Println(s1)
}
```


### 闭包
指的是**一个函数和与其相关引用环境组成而成的实体**，就是 `闭包=函数+引用`
```
func adder() func(int) int {
    var x int           //开始默认值为0
    return func(y int) int{     //返回函数的引用
        x += y
        return x
    }
}

func main(){
    var f = adder()
    fmt.Println(f(10))  //10
    fmt.Println(f(20))  //30
    fmt.Println(f(30))  //60
    
    f1 := adder()
    fmt.Println(f1(40)) //40
	fmt.Println(f1(50)) //90
}
```

变量`f`是一个函数并且它引用了其外部作用域中的`x`的变量，此时`f`就是一个闭包，在`f`的生命周期里，变量`x`一直有效

闭包开始传参：
```
func adder2(x int) func(int) int {
	return func(y int) int {
		x += y
		return x
	}
}
func main() {
	var f = adder2(10)
	fmt.Println(f(10)) //20
	fmt.Println(f(20)) //40
	fmt.Println(f(30)) //70

	f1 := adder2(20)
	fmt.Println(f1(40)) //60
	fmt.Println(f1(50)) //110
}
```
闭包示例2

```
func makeSuffixFunc(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

func main() {
	jpgFunc := makeSuffixFunc(".jpg")
	txtFunc := makeSuffixFunc(".txt")
	fmt.Println(jpgFunc("test")) //test.jpg
	fmt.Println(txtFunc("test")) //test.txt
}
```
闭包示例3
```
func calc(base int) (func(int) int, func(int) int) {
	add := func(i int) int {
		base += i
		return base
	}

	sub := func(i int) int {
		base -= i
		return base
	}
	return add, sub
}

func main() {
	f1, f2 := calc(10)
	fmt.Println(f1(1), f2(2)) //11 9
	fmt.Println(f1(3), f2(4)) //12 8
	fmt.Println(f1(5), f2(6)) //13 7
}
```

## `defer` 语句
**延迟处理**,在`defer`归属的函数即将返回时，将延迟处理的语句按`defer`定义的逆序进行执行，也就是说，先被`defer`的语句最后被执行，最后被`defer`的语句，最先被执行。

```
func main() {
	fmt.Println("start")
	defer fmt.Println(1)
	defer fmt.Println(2)
	defer fmt.Println(3)
	fmt.Println("end")
}

//输出结果
start
end 
3
2
1
```
所以`defer`语句能非常方便的处理资源释放问题。比如：**资源清理、文件关闭、解锁及记录时间等**。

### `defer`执行时机
在`Go`语言的函数中`return`语句在底层并不是原子操作，它分为给返回值赋值和RET指令两步。而`defer`语句执行的时机就在返回值赋值操作后，`RET`指令执行前。具体如下图所示：
![image](https://www.liwenzhou.com/images/Go/func/defer.png)

### `defer` 面试题
```
func calc(index string, a, b int) int {
	ret := a + b
	fmt.Println(index, a, b, ret)
	return ret
}

func main() {
	x := 1
	y := 2
	defer calc("AA", x, calc("A", x, y))
	x = 10
	defer calc("BB", x, calc("B", x, y))
	y = 20
}

//因为延迟特性
x, y= 10, 20
```

# 内置函数

| 内置函数         | 介绍                          |
| ---------------- | ----------------------------- |
| `close`          | 主要用于关闭`channel`         |
| `len`            | 用于求长度                    |
| `new`            | 用来分配内存                  |
| `make`           | 用来分配内存                  |
| `append`         | 用于追加元素到数组，`slice`章 |
| `panic和recover` | 用于错误处理                  |

## `panic/recover`
`panic`可以让程序触发宕机，将堆栈和`goroutine`信息输出到控制台
```
语句： panic(value inerface{})  //panic()参数可以任意类型
```
编译正则表达式捕获错误
```
函数介绍：func Compile(expr string) (*Regexp, error) 编译正则表达式，发生错误时返回编译错误并且Regexp为nil

func MustComplie(str string) *Regexp {
    regexp, error := Compile(str)
    if error != nil {
        panic('regexp: Compile('+ quote(str)+ '):' +error.Error())  //触发宕机
    }
    return regexp       //只返回Regexp
}

```

## 在宕机时触发延迟执行语句
`panic()`触发宕机后，是不会执行后面代码的，但是它之前执行过的 `defer` 语句可以在宕机发生时发生作用

```
func funcA() {
	fmt.Println("func A")
}

func funcB() {
	panic("panic in B")
}

func funcC() {
	fmt.Println("func C")
}
func main() {
	funcA()
	funcB()
	funcC()
}

// 输出

func A
panic: panic in B       //只运行到funcB panic语句

goroutine 1 [running]:
main.funcB(...)
        .../code/func/main.go:12
main.main()
        .../code/func/main.go:20 +0x98
```

程序运行期间 `funcB` 中引发了 `panic` 导致程序崩溃，异常退出，这时候我们就可以通过`recover`将程序恢复回来，继续往后执行。
```
func funcA() {
    fmt.Println("func A")
}

func funcB() {
    defer func() {
        err := recover()        // 程序出现了panic错误，可以通过recover恢复过来
        if err != nil {         //产生报错，error不为空
            fmt.Println("recover in B")
        }
    }()
    panic("panic in B")
}

func funcC() {
    fmt.Println("func C")
}

func main() {
    funcA()
    funcB()
    funcC()
}
```
注意：
* `recover()`**必须搭配**`defer`**使用**
* `defer` **一定要在引发panic语句之前定义**


# 面试题
```
/*
你有50枚金币，需要分配给以下几个人：Matthew,Sarah,Augustus,Heidi,Emilie,Peter,Giana,Adriano,Aaron,Elizabeth。
分配规则如下：
a. 名字中每包含1个'e'或'E'分1枚金币
b. 名字中每包含1个'i'或'I'分2枚金币
c. 名字中每包含1个'o'或'O'分3枚金币
d: 名字中每包含1个'u'或'U'分4枚金币
写一个程序，计算每个用户分到多少金币，以及最后剩余多少金币？
程序结构如下，请实现 ‘dispatchCoin’ 函数
*/
package main

import "fmt"

var (
	coins = 50
	users = []string{
	"Matthew", "Sarah", "Augustus", "Heidi", "Emilie", "Peter", "Giana", "Adriano", "Aaron", "Elizabeth",
	}
	distribution = make(map[string]int, len(users))
)

func dispatchCoin () int {
	//定义剩下的金币
	var left int = coins
	//遍历每个users
	for i := 0; i < len(users); i++ {
		name := users[i]
		fmt.Println("给" + name + "分配金币中")
		var use int
		//遍历每个name中字符
		for j := 0; j < len(name); j++ {
			str := string(name[j])	//需要使用string强制转换
			fmt.Println("字母：",str)
			switch str {
			case "e", "E":
				use++
			case "i", "l":
				use += 2
			case "o", "O":
				use += 3
			case "u", "U":
				use += 4
			default:
				fmt.Println("其他字母，跳过，不添加金币")
			}
		}
		distribution[name] = use
		// 得到剩下的
		left = left - use
		fmt.Println("现在只剩下：",left)
	}
	return left
}


func main() {
	left := dispatchCoin()
	fmt.Println("剩下：", left)
}
```