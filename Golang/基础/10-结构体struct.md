- [类型别名和自定义类型](#类型别名和自定义类型)
	- [自定义类型](#自定义类型)
	- [类型别名](#类型别名)
	- [类型定义和类型别名区别](#类型定义和类型别名区别)
- [结构体](#结构体)
	- [结构体定义](#结构体定义)
	- [结构体实例化](#结构体实例化)
	- [创建指针类型结构体](#创建指针类型结构体)
	- [取结构体的地址实例化](#取结构体的地址实例化)
	- [`struct`的值和指针](#struct的值和指针)
	- [匿名结构体](#匿名结构体)
- [结构体初始化](#结构体初始化)
		- [使用键值初始化](#使用键值初始化)
- [结构体内存布局](#结构体内存布局)
	- [应用题：查看下面代码执行结果](#应用题查看下面代码执行结果)
- [构造函数](#构造函数)
	- [方法和接收者](#方法和接收者)
	- [指针类型的接收者](#指针类型的接收者)
	- [值类型的接收者](#值类型的接收者)
	- [任意类型添加方法](#任意类型添加方法)
		- [结构体匿名字段](#结构体匿名字段)
- [嵌套结构体](#嵌套结构体)
		- [递归`struct`：嵌套自身](#递归struct嵌套自身)
		- [嵌套匿名结构体](#嵌套匿名结构体)
	- [结构体的继承](#结构体的继承)
		- [结构体字段的可见性](#结构体字段的可见性)
- [结构体序列化数据](#结构体序列化数据)
- [结构体标签`(Tag)`](#结构体标签tag)
- [练习题](#练习题)


`Go`语言中没有“类”的概念，也不支持“类”的继承等面向对象的概念。`Go`语言通过结构体的内嵌再配合接口

# 类型别名和自定义类型
## 自定义类型
在基本数据类型中，如`string, int, float, bool`等数据类型，`Go` 语言可以使用`type`关键字自定义类型；

自定义类型是定义了一个全新的类型。我们可以基于内置的基本类型定义，也可以通过`struct`定义。
```
//将MyInt定义为int类型
type MyInt int
```
通过`Type`关键字定义， `MyInt`就是一种新类型，基于 `int`特性

## 类型别名
```
type TypeAlias = Type

//我们之前见过的 rune 和 byte 就是类型别名
type byte = uint8
type rune = int32
```
## 类型定义和类型别名区别
```
//类型定义
type NewInt int

//类型别名
type MyInt = int

func main() {
	var a NewInt
	var b MyInt
	
	fmt.Printf("type of a:%T\n", a) //type of a:main.NewInt 
	fmt.Printf("type of b:%T\n", b) //type of b:int
}
```
结果显示：
* `a`类型为`main.NewInt`，表示`main`包下定义的`NewInt`类型
* `b`类型为`int`,`MyInt` 类型只会在代码中存在，编译完成时并不会有`MyInt`类型

# 结构体
`Go`语言提供一种自定义数据类型，可以封装多个基本数据类型，这种数据类型叫**结构体**，`struct`定义

使用`struct`实现面向对象

## 结构体定义
使用 `type` 和 `struct`关键字定义
```
type name struct {
    字段名 字段类型
    字段名 字段类型
    ...
}
```
其中：
* `name` 类型名：标识结构体名称，同一个包内不能重复；
* `字段名`: 表示结构体字段名，结构体中字段名必须唯一；
* `字段类型`： 表示结构体字段具体类型

比如：定义一个 `Person`(人)结构体：
```
type Person struct {
    name string     
    city string
    age int8
}
//同种类型也可以写到一行种
type Person struct{
    name, city string
    age int8
}

//这样就可以使用 person结构体方便表示和存储
```
结构体`struct`是用来描述一组值的，**本质上是一种聚合型的数据类型**

## 结构体实例化
**只有当结构体实例化，才会真正地分配内存，也就是必须实例化后才能使用结构体字段**
```
var 结构体示例 结构体类型
```
基本的实例化形式
```
type Person struct {
    name, city string
    age int8
}

func main() {
    var p1 Person
    p1.name = "古力娜扎"
    p1.city = "北京"
    p1.age = 18
    fmt.Printf("p1=%v\n", p1)  //p1={沙河娜扎 北京 18}
	fmt.Printf("p1=%#v\n", p1) //p1=main.person{name:"沙河娜扎", city:"北京", age:18}
}
```
我们通过 `.`访问结构体字段（成员变量）,`p1.name` 和`p1.age`等

## 创建指针类型结构体
还可以通过使用 `new`关键字对结构体实例化，得到是结构体的地址
```
ins := new(Type) 
//Type 可以是结构体，整型，字符串
// ins 类型被实例化后保存在ins变量，ins类型为*Type,属于指针
```
```
var p2 = new(person)
fmt.Printf("%T\n", p2)     //*main.person
fmt.Printf("p2=%#v\n", p2) //p2=&main.person{name:"", city:"", age:0}

//p2 是一个结构体指针
//支持对结构体指针直接使用 . 访问结构体成员
var p2 = new(person)
p2.name = "小王子"
p2.age = 28
p2.city = "上海"
fmt.Printf("p2=%#v\n", p2) //p2=&main.person{name:"小王子", city:"上海",age:28}
```
也可以使用 `&` 对结构体进行取地址操作，**也相当于对该结构体类型进行了一次 `new`实例化操作：
```
p3 := &person{}
fmt.Printf("%T\n", p3)     //*main.person
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"", city:"", age:0}
//&取完地址，使用*取地址的值 ，p3.name 语法糖是 (*p3).name
p3.name = "七米"            
p3.age = 30
p3.city = "成都"
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"七米", city:"成都", age:30}
```
## 取结构体的地址实例化
```
ins := &Type{}
```
其中：
* `Type`表示结构体类型
* `ins`为结构体示例，类型为 `*Type`， 属于指针类型
```
type Command struct {
    Name    string    // 指令名称
    Var     *int      // 指令绑定的变量, 使用整型指针绑定一个指针
    Comment string    // 指令的注释
}
var version int = 1

cmd := &Command{}   // 实例化
cmd.Name = "version"
cmd.Var = &version  // &version 返回version的 *int类型
cmd.Comment = "show version"
```
## `struct`的值和指针
下面三种方式都可以构造`person struct`的实例 `p`：
```
p1 := person{}
p2 := &person{}
p3 := new(person)
```
输出一下：
```
type person struct {
    name string
    age int
}

func main() {
    p1 := person{}
    p2 := &person{}
    p3 := new(person)
    fmt.Println(p1)
    fmt.Println(p2)
    fmt.Println(p3)
}

//结果
{0}
&{0}
&{0}
```
`p1,p2,p3`都是`person struct`实例，但是`p2``p3`是完全等价的，都指向实例的指针，指针再指向实例，而`p1`直接指向实例
```
变量名      指针        数据对象(实例)
--------------------------------------
p1(addr)----------------> {0}
p2-------->ptr(addr)----->{0}
p3-------->ptr(addr)----->{0}
```
所以`p1`和`ptr(addr)`保存的都是数据对象地址，而`p2, p3`则保存`ptr(addr)`地址，我们**称指向指针的变量**`p2,p3`**为指针**，将直接指向数据对象的变量称为对象本身。

## 匿名结构体
**使用在一些临时数据结构等场景下**
```
package main
import(
    "fmt"
)

func main() {
    var user struct{Name string; Age int}
    user.Name = "小王子"
    user.Age = 18
    fmt.Printf("%#v\n", user)
}
```

# 结构体初始化
```
type person struct {
	name string
	city string
	age  int8
}

func main() {
	var p4 person       // 默认初始值为其类型零值
	fmt.Printf("p4=%#v\n", p4) //p4=main.person{name:"", city:"", age:0}
}
```

### 使用键值初始化
```
p5 := &person{
	name: "小王子",
	city: "北京",
	age:  18,
}
fmt.Printf("p5=%#v\n", p5) //p5=main.person{name:"小王子", city:"北京", age:18}
```
**也可以使用值的列表初始化**
```
p8 := &person{
	"沙河娜扎",
	"北京",
	28,
}
```
注意：
1. 必须初始化结构体所有字段
2. 填充顺序必须与字段在结构体声明顺序一致
3. 不能和键值初始化方式混用

# 结构体内存布局
**结构体占用一块连续的内存**
```
type test struct {
	a int8
	b int8
	c int8
	d int8
}
n := test{
	1, 2, 3, 4,
}
fmt.Printf("n.a %p\n", &n.a)
fmt.Printf("n.b %p\n", &n.b)
fmt.Printf("n.c %p\n", &n.c)
fmt.Printf("n.d %p\n", &n.d)

//输出
n.a 0xc0000a0060
n.b 0xc0000a0061
n.c 0xc0000a0062
n.d 0xc0000a0063
```

## 应用题：查看下面代码执行结果
```
type student struct {
	name string
	age  int
}

func main() {
	m := make(map[string]*student)     //值类型为 student指向的值 //定义map，存的类型为student结构体
	stus := []student{                  //定义stus切片，嵌套map
		{name: "小王子", age: 18},
		{name: "娜扎", age: 23},
		{name: "大王八", age: 9000},
	}

	for _, stu := range stus {          //遍历切片，获取每个值，保存在m(映射)中，键位stu.name,值为stu的内存地址
		m[stu.name] = &stu      
	}
	for k, v := range m {
		fmt.Println(k, "=>", v.name)    //遍历输出 键和值.name, 使用了语法糖 *v.name
	}
}
//输出
小王子 => 小王子
...
```

# 构造函数
`Go`语言结构体没有构造函数，我们可以自己实现
```
func newPerson(name, city string, age int8) *person {
    return &person{     //返回的是结构体指针类型，如果结构体复杂，值拷贝性能开销比较大
        name: name,
        city: city,
        age: age,
    }
}

//调用构造函数
p9 := newPerson("张三", "沙河"， 90)    //返回其实例化
fmt.Printf("%#v\n", p9) //&main.person{name:"张三", city:"沙河", age:90}
```
## 方法和接收者
`方法（Method）`是一种作用于特定类型变量的函数，这种特定类型变量叫做 `接收者(Receiver)`。 

接收者概念类似于其他语言中 `this`或者`self`
```
func (接收者变量 接收者类型) 方法名(参数列表) (返回参数) {
    函数体
}
```
其中：
* **接收者变量**：接收者中的参数变量名在命名时，官方建议使用接收者类型名的第一个小写字母，而不是 `self, this` 之类的命名。例如： `Person`类型接收者变量命名为 `p`，`Connector`类型接收者变量命名为 `c` 等
* 接收者类型：接收者类型和参数类似，可以是指针类型和非指针类型
* 方法名，参数列表，返回参数：具体格式与函数定义相同
```
//Person 结构体
type Person struct {
	name string
	age  int8
}

//NewPerson 构造函数
func NewPerson(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

//定义Dream Person做梦的方法 ，接收者类型为Person
func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好Go语言！\n", p.name)
}

func main() {
	p1 := NewPerson("小王子", 25)   //使用构造函数，返回的是结构体指针类型
	p1.Dream()          //调用方法
}
```

## 指针类型的接收者
接收者由一个结构体指针组成，指针特性，调用方法时修改接收者指针的任意成员变量，在方法结束猴，修改猴都是有效的，接近于其他语言中的`this`或者`self`
```
// SetAge 设置p的年龄
// 使用指针接收者
func (p *Person) SetAge(newAge int8) {
	p.age = newAge
}

func main() {
	p1 := NewPerson("小王子", 25)
	fmt.Println(p1.age) // 25
	p1.SetAge(30)
	fmt.Println(p1.age) // 30
}
```
## 值类型的接收者
当方法作用于值类型接收者时，**Go语言会在代码运行时将接收者的值复制一份**。在值类型接收者的方法中可以获取接收者的成员值，**但修改操作只是针对副本，无法修改接收者变量本身**。
```
// SetAge2 设置p的年龄
// 使用值接收者
func (p Person) SetAge2(newAge int8) {
	p.age = newAge
}

func main() {
	p1 := NewPerson("小王子", 25)
	p1.Dream()
	fmt.Println(p1.age) // 25
	p1.SetAge2(30) // (*p1).SetAge2(30)
	fmt.Println(p1.age) // 25
}
```
## 任意类型添加方法
不仅仅是结构体，任何类型都可以拥有方法。
```
//MyInt 将int定义为自定义MyInt类型
type MyInt int

//SayHello 方法
func (m MyInt) SayHello() {
    fmt.Println("Hello, i am an int")
}

func main() {
    var m1 MyInt
    m1.SayHello()   //Hello, i am an int
    m1 = 100
    fmt.Printf("%#v %T\n", m1, m1)  //100   main.MyInt
}
```
**注意：非本地类型不能定义方法，也就是不能给别的包类型定义方法**

### 结构体匿名字段
匿名字段默认采用类型名作为字段名，**结构体要求字段名必须唯一，所以同一个匿名类型只有一个**
```
//Person 结构体Person类型
type Person struct {
	string
	int
}

func main() {
	p1 := Person{
		"小王子",
		18,
	}
	fmt.Printf("%#v\n", p1)        //main.Person{string:"北京", int:18}
	fmt.Println(p1.string, p1.int) //北京 18
}
```

# 嵌套结构体
**一个结构体中可以嵌套包含另一个结构体或者结构体指针**
```
//Address 地址结构体
type Address struct {
	Province string
	City     string
}

//User 用户结构体
type User struct {
	Name    string
	Gender  string
	Address Address
}

func main() {
	user1 := User{
		Name:   "小王子",
		Gender: "男",
		Address: Address{
			Province: "山东",
			City:     "威海",
		},
	}
	// user1=main.User{Name:"小王子", Gender:"男", Address:main.Address{Province:"山东", City:"威海"}}
	fmt.Printf("user1=%#v\n", user1)
	
}
```
### 递归`struct`：嵌套自身
如果`struct`中嵌套的`struct`类型是自己的指针类型，可以用来生成特殊的数据结构：链表或者二叉树(双端链表)
```
type Node struct {
    data string
    ri *Node
}
//链表结构示意图：
-----|----          -----|-----
| data | ri  | --> | data | ri |
| ---- | --- |-----|-----
```
如果**给嵌套两个自己的指针，每个结果都有一个左指针和一个右指针，分别指向它的左边节点和右边节点，就形成了二叉树或是双端链表数据结构**

定义二叉树
```
type Tree struct {
    le *Tree
    data string
    ri *Tree
}
//最初生成二叉树时，`root`节点没有任何指向
//root节点：初始左右两端为空
root := new(Tree)
root.data = "root node"

//随着节点增加，root节点开始指向其他左节点，右节点或者其他，二叉树添加节点时候，只需将新生成节点赋值给它前一个节点的le或者ri字段即可
//生成两个新节点：
newLeft := new(Tree)
newLeft.data = "left node"
newRight := &Tree{nil, "Right node", nil}

//添加到树种
root.le = newLeft
root.ri = newRight

//再添加一个新节点到newLeft节点的右节点
anotherNode := &Tree{nil, "another Node", nil}
newLeft.ri = anotherNode

```
### 嵌套匿名结构体
**当访问结构体成员时会先在结构体中查找该字段，找不到再去匿名结构体中查找。**
```
//Address 地址结构体
type Address struct {
	Province string
	City     string
}

//User 用户结构体
type User struct {
	Name    string
	Gender  string
	Address //匿名结构体, 会在结构体寻找其字段，找不到才再去匿名结构体中
}

func main() {
	var user2 User
	user2.Name = "小王子"
	user2.Gender = "男"
	user2.Address.Province = "山东"    //通过匿名结构体.字段名访问
	user2.City = "威海"                //直接访问匿名结构体的字段名
	fmt.Printf("user2=%#v\n", user2) //user2=main.User{Name:"小王子", Gender:"男", Address:main.Address{Province:"山东", City:"威海"}}
}
```
**还有种情况：嵌套结构体内部可能存在相同的字段名。**

这个时候为了避免歧义 **需要指定具体的内嵌结构体的字段**

## 结构体的继承
```
//Animal 动物
type Animal struct {
	name string
}

func (a *Animal) move() {
	fmt.Printf("%s会动！\n", a.name)
}

//Dog 狗
type Dog struct {
	Feet    int8
	*Animal //通过嵌套匿名结构体实现继承
}

func (d *Dog) wang() {
	fmt.Printf("%s会汪汪汪~\n", d.name)
}

func main() {
	d1 := &Dog{
		Feet: 4,
		Animal: &Animal{ //注意嵌套的是结构体指针
			name: "乐乐",
		},
	}
	d1.wang() //乐乐会汪汪汪~
	d1.move() //乐乐会动！
}
```
### 结构体字段的可见性
结构体中字段**大写开头表示可公开访问，小写私有（仅在定义当前结构体包中可访问）**

# 结构体序列化数据
```
package main

import (
	"encoding/json"
	"fmt"
)

//student
type Student struct {
	ID int
	Gender, Name string
}

//class
type Class struct {
	Title string
	//存切片
	Students []*Student	//类型为*Student指针类型
}

func main() {
	//实例化
	class := new(Class)
	class.Title = "101"
	class.Students = make([]*Student, 0, 200)
	//或者使用
	/*class := &Class{
		Title:    "101",
		Students: make([]*Student, 0, 200),
	}*/

	//放入student数据
	for i := 0; i < 10; i++ {
		//实例化每个学生
		stu := &Student{
			Name:   fmt.Sprintf("stu%02d", i),
			Gender: "男",
			ID:     i,
		}
		//添加
		class.Students = append(class.Students, stu)
	}
	//Json序列化：结构体 --> Json格式字符串
	data, err := json.Marshal(class)
	if err != nil{
		fmt.Println("json marshal failed")
		return
	}
	fmt.Printf("json:%s\n", data)

	//Json反序列化：Json格式字符串 --> 结构体
	str := `{"Title":"101","Students":[{"ID":0,"Gender":"男","Name":"stu00"},{"ID":1,"Gender":"男","Name":"stu01"},{"ID":2,"Gender":"男","Name":"stu02"},{"ID":3,"Gender":"男","Name":"stu03"},{"ID":4,"Gender":"男","Name":"stu04"},{"ID":5,"Gender":"男","Name":"stu05"},{"ID":6,"Gender":"男","Name":"stu06"},{"ID":7,"Gender":"男","Name":"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"stu09"}]}`
	// 初始化空实例
	c1 := &Class{}
	err = json.Unmarshal([]byte(str), c1)
	if err != nil {
		fmt.Println("json unmarshal failed!")
		return
	}
	fmt.Printf("%#v\n", c1)
}

//序列化结果
json:{"Title":"101","Students":[{"ID":0,"Gender":"男","Name":"stu00"},{"ID":1,"Gender":"男","Name":"stu01"},{"ID":2,"Gender":"男","Name":"stu02"},{"ID":3,"Gender":
"男","Name":"stu03"},{"ID":4,"Gender":"男","Name":"stu04"},{"ID":5,"Gender":"男","Name":"stu05"},{"ID":6,"Gender":"男","Name":"stu06"},{"ID":7,"Gender":"男","Name"
:"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"stu09"}]}

//反序列化结果
&main.Class{Title:"101", Students:[]*main.Student{(*main.Student)(0xc000066810), (*main.Student)(0xc000066840), (*main.Student)(0xc000066870), (*main.Student)(0xc0
:"stu07"},{"ID":8,"Gender":"男","Name":"stu08"},{"ID":9,"Gender":"男","Name":"s
tu09"}]}
&main.Class{Title:"101", Students:[]*main.Student{(*main.Student)(0xc000066810)
, (*main.Student)(0xc000066840), (*main.Student)(0xc000066870), (*main.Student)(0xc0000668a0), (*main.Student)(0xc000066900), (*main.Student)(0xc000066930), (*main
.Student)(0xc000066960), (*main.Student)(0xc000066990), (*main.Student)(0xc0000669c0), (*main.Student)(0xc0000669f0)}}

```

# 结构体标签`(Tag)`
`Tag` 是结构体元信息，在运行时候通过反射机制读取出来。

`Tag` 在结构体字段后方定义, 由一对反引号包裹起来：
```
`key1:"value1" key2:"value2"`
```
**结构体标签由一个或多个键值对组成。键与值使用冒号分隔，值用双引号括起来。键值对之间使用一个空格分隔。**

**格式严谨，出错不会提示，不要在key和value之间添加空格**

```
//Student 学生
type Student struct {
	ID     int    `json:"id"`   //通过指定tag实现json序列化该字段时的key
	Gender string               //json序列化是默认使用字段名作为key
	name   string               //私有不能被json包访问
}

func main() {
	s1 := Student{
		ID:     1,
		Gender: "男",
		name:   "沙河娜扎",
	}
	data, err := json.Marshal(s1)
	if err != nil {
		fmt.Println("json marshal failed!")
		return
	}
	fmt.Printf("json str:%s\n", data) //json str:{"id":1,"Gender":"男"}
}
```

# 练习题
```
package main

import "fmt"

//需要先定义结构体，相当于Python的__init__初始化
type Person struct {
	name, gender string
	age int
}

//定义构造函数
func NewPerson(name string, age int, gender string) *Person {
	//返回其实例化
	return &Person{
		name: name,	//需要有逗号
		age: age,
		gender: gender,
	}
}

//定义方法, 实例化时候产生的是 *Person指针类型， 接收*Person也是指针类型，方便修改值
func (p *Person) run() {
	fmt.Println("he is running")
	p.age -= 5
	fmt.Printf("and %v become young, age is %v", p.name, p.age)
}

//使用继承, 需要先定义子类，并且实例化
type Man struct {
	height int 	//自己的属性
	//继承属性
	*Person		//接收的是*Person指针类型
}

//定义Man自己的方法
func (m Man) xuxu() {
	fmt.Printf("\n %v is xuxuing\n", m.name)
}

//构造函数
func NewMan (height int) *Man {
	return &Man{
		height,
		&Person{
			name:   "liukaitao",
			gender: "man",
			age:    20,
		},
	}
}

//主体
func main() {
	p1 := NewPerson("liukaitao", 24, "man")
	fmt.Printf("%#v\n", p1)	//&main.Person{name:"liukaitao", gender:"man", age:24}
	// 使用自定义方法，run()
	p1.run()
	//p1.xuxu()	// (type *Person has no field or method xuxu)
	fmt.Printf("\n age is %v\n", p1.age)

	//实例化Man
/*	liukaitao := &Man{
		height: 175,		//初始化自己属性
		Person: &Person{
			name:   "liukaitao",
			gender: "man",
			age:    20,
		},			//初始化父类属性
	}*/
	// 或者使用构造函数实例化
	liukaitao := NewMan(175)
	liukaitao.run()
	liukaitao.xuxu()
}
```