- [接口](#接口)
	- [接口类型](#接口类型)
	- [为什么使用接口](#为什么使用接口)
	- [接口定义](#接口定义)
	- [实现接口的条件](#实现接口的条件)
	- [接口类型变量](#接口类型变量)
	- [值接收者和指针接收者实现接口的区别](#值接收者和指针接收者实现接口的区别)
		- [使用 值接收者 实现接口](#使用-值接收者-实现接口)
		- [使用 指针接收者 实现接口](#使用-指针接收者-实现接口)
	- [完整例子](#完整例子)
- [类型与接口关系](#类型与接口关系)
	- [一个类型实现多个接口](#一个类型实现多个接口)
	- [多个类型实现同一个接口](#多个类型实现同一个接口)
	- [接口嵌套](#接口嵌套)
- [空接口](#空接口)
	- [空接口定义](#空接口定义)
	- [空接口应用](#空接口应用)
		- [空接口作为函数参数](#空接口作为函数参数)
		- [空接口作为`map`的值](#空接口作为map的值)
- [类型断言](#类型断言)
	- [接口值](#接口值)
- [更多详细类型断言看这里](#更多详细类型断言看这里)



接口`(interface)`定义了一个对象的行为规范，只定义规范不是先，由具体对象实现规范细节

# 接口
## 接口类型
`interface` 是一组 `method`集合，是鸭子类型 `duck-type programming`的一种体现。接口做的事情就是定义一个协议（规则）

只要一台机器由洗衣服和甩干功能，就叫他洗衣机，不关心属性（数据），只关心行为（方法）

**为了保护你的Go语言职业生涯，请牢记接口（interface）是一种类型**

## 为什么使用接口
接口区别于我们之前所有的具体类型，**接口是一种抽象的类型。**

当你看到一个接口类型的值时，你不知道它是什么，**唯一知道的是通过它的方法能做什么。**

## 接口定义
`Go`语言提倡面向接口编程

```
type 接口类型名 interface{
    方法名1（参数列表1） 返回值列表1
    方法名2（参数列表2） 返回值列表2
    ...
}
```
其中：
* 接口类型名：使用 `type`定义，接口命名时，一般会再单词后面添加 `er`
* 方法名： 当方法名且这个接口类型名首字母大写，表示可以被包之外代码访问
* 参数列表，返回值列表： 参数变量名可以省略

## 实现接口的条件
接口就是**一个需要实现的方法列表**
```
package main
//Sayer 接口
import "fmt"

type Sayer interface {
	say()
}

//定义 dog 和 cat 两个结构体
type dog struct {}

type cat struct {}

// dog实现了Sayer接口
func (d dog) say() {
	fmt.Println("汪汪汪")
}

// cat实现了Sayer接口
func (c cat) say() {
	fmt.Println("喵喵喵")
}

func main () {
	//实例化
	c := &cat{}
	c.say()
}
```
## 接口类型变量
那实现了接口怎么使用呢？ `Sayer` 类型变量能够存储 `dog`和 `cat`类型变量
```
func main() {
	var x Sayer // 声明一个Sayer类型的变量x
	a := cat{}  // 实例化一个cat
	b := dog{}  // 实例化一个dog
	x = a       // 可以把cat实例直接赋值给x
	x.say()     // 喵喵喵
	x = b       // 可以把dog实例直接赋值给x
	x.say()     // 汪汪汪
}
```

## 值接收者和指针接收者实现接口的区别
大家都有`Mover`接口和一个`dog`结构体
```
type Mover interface {
    move()
}

type dog struct {}
```
### 使用 值接收者 实现接口
```
func (d dog) move() {       //接收者类型为 dog值
    fmt.Println("狗会动")
}
```
此时实现接口的是 `dog` 类型：
```
func main(){
    var x Mover
    var wangcai = dog{}     //旺财是dog类型
    x = wangcai             // x可以接收dog类型
    x.move()
    
    var fugui = &dog{}      // 富贵是*dog类型
    x = fugui               // x可以接收*dog类型
    //语法糖：等价于 x = *fugui
    x.move()
}

//输出
狗会动
狗会动
``` 
不管是dog结构体还是结构体指针 `*dog` 类型的变量都可以赋值给该接口变量。因为 `Go` 语言中有对指针类型变量求值的语法糖，`dog`指针`fugui`内部会自动求值`*fugui`。
### 使用 指针接收者 实现接口
```
func (d *dog) move() {
	fmt.Println("狗会动")
}
func main() {
	var x Mover
	var wangcai = dog{} // 旺财是dog类型
	x = wangcai         // x不可以接收dog类型
	
	var fugui = &dog{}  // 富贵是*dog类型
	x = fugui           // x可以接收*dog类型
}
```
此时实现Mover接口的是`*dog`类型，所以不能给`x`传入`dog`类型的`wangcai`，此时`x`只能存储`*dog`类型的值。

## 完整例子
```
type People interface {
	Speak(string) string
}

type Student struct{}

func (stu *Student) Speak(think string) (talk string) {
	if think == "sb" {
		talk = "你是个大帅比"
	} else {
		talk = "您好"
	}
	return
}

func main() {
    // var peo People
    // var stu = &Student{}
    // peo = stu
	var peo People = &Student{}
	think := "bitch"
	fmt.Println(peo.Speak(think))
}
```

# 类型与接口关系
## 一个类型实现多个接口
一个类型可以同时实现多个接口，而接口间彼此独立
## 多个类型实现同一个接口
不同的类型还可以实现同一接口

并且一个接口的方法，不一定需要由一个类型完全实现，接口的方法可以通过在类型中嵌入其他类型或者结构体来实现。

```
// WashingMachine 洗衣机
type WashingMachine interface {
	wash()
	dry()
}

// 甩干器
type dryer struct{}

// 实现WashingMachine接口的dry()方法
func (d dryer) dry() {
	fmt.Println("甩一甩")
}

// 海尔洗衣机
type haier struct {
	dryer //嵌入甩干器
}

// 实现WashingMachine接口的wash()方法
func (h haier) wash() {
	fmt.Println("洗刷刷")
}
```

## 接口嵌套
接口与接口间可以通过嵌套创造新的接口,使用方法跟普通接口一样
```
// Sayer 接口
type Sayer interface {
	say()
}

// Mover 接口
type Mover interface {
	move()
}

// 接口嵌套
type animal interface {
	Sayer
	Mover
}

type cat struct {
	name string
}

func (c cat) say() {
	fmt.Println("喵喵喵")
}

func (c cat) move() {
	fmt.Println("猫会动")
}

func main() {
	var x animal
	x = cat{name: "花花"}
	// var x animal = &cat{name: "花花"}
	x.move()
	x.say()
}
```

# 空接口
## 空接口定义
空接口是指没有定义任何方法的接口,空接口类型的变量可以存储任意类型的变量。
```
func main() {
	// 定义一个空接口x
	var x interface{}
	s := "Hello 沙河"
	x = s
	fmt.Printf("type:%T value:%v\n", x, x)
	i := 100
	x = i
	fmt.Printf("type:%T value:%v\n", x, x)
	b := true
	x = b
	fmt.Printf("type:%T value:%v\n", x, x)
}
```
## 空接口应用
### 空接口作为函数参数
使用空接口实现可以接收任意类型的函数参数
```
// 空接口作为函数参数
func show(a interface{}) {
	fmt.Printf("type:%T value:%v\n", a, a)
}
```
### 空接口作为`map`的值
空接口实现可以保存任意值的`map`映射
```
// 空接口作为map值
	var studentInfo = make(map[string]interface{})
	studentInfo["name"] = "沙河娜扎"
	studentInfo["age"] = 18
	studentInfo["married"] = false
	fmt.Println(studentInfo)
```

# 类型断言
空接口可以存储任意类型的值，那如何获取其存储的具体数据呢？

## 接口值
一个接口的值（简称接口值）是由 `一个具体类型` 和 `具体类型的值`两部分组成的。

这两部分分别称为接口的 `动态类型`和 `动态值`;

```
var w io.Writer
w = nil
w = os.Stdout
w = new(bytes.Buffer)
```
![分解图](https://www.liwenzhou.com/images/Go/interface/interface.png)

**想要判断空接口中的值，这个时候就可以使用类型断言**, 格式如下：
```
x.(T)
```
其中：
* `x`：表示类型为 `interface{}`的变量
* `T`: 表示断言， `x`可能是类型
*  返回两个参数，第一个是 `x` 转化为 `T`类型后变量；第二个是**布尔值**，若为`true`表示断言成功

```
func main() {
	var x interface{}
	x = "Hello 沙河"
	v, ok := x.(string)
	if ok {
		fmt.Println(v)
	} else {
		fmt.Println("类型断言失败")
	}
}

//断言多次
func justifyType(x interface{}) {
	switch v := x.(type) {      //断言，x.(T)
	case string:
		fmt.Printf("x is a string，value is %v\n", v)
	case int:
		fmt.Printf("x is a int is %v\n", v)
	case bool:
		fmt.Printf("x is a bool is %v\n", v)
	default:
		fmt.Println("unsupport type！")
	}
}
```
**只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要定义接口。**

# 更多详细类型断言看这里
[Go基础系列：接口类型断言和type-switch](https://www.cnblogs.com/f-ck-need-u/p/9893347.html)

