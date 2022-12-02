# `fmt`标准库
实现格式化`I/O`，分为向外输出内容和获取输入内容两部分;

# 向外输出
## `Print`
`Print`系列函数会将内容输出到系统的标准输出

`Print`函数直接输出内容，`Printf`函数支持格式化输出字符串, `Println`函数会在输出内容结尾条件`\n`

```
func Print(a ...interface{}) (n int, err error)
func Printf(format string, a ...interface{}) (n int, err error)
func Println(a ...interface{}) (n int, err error)
```
## `Fprint`
`Fprint`系列函数会将内容输出到一个`io.Writer`接口类型的变量`w`中，通常用这个函数往文件中写入内容。
```
//interface{}为交互内容
func Fprint(w io.Writer, a ...interface{}) (n int, err error)
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
func Fprintln(w io.Writer, a ...interface{}) (n int, err error)
```
举个例子
```
// 向标准输出写入内容
fmt.Fprintln(os.Stdout, "向标准输出写入内容")
fileObj, err := os.OpenFile("./xx.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
if err != nil {
	fmt.Println("打开文件出错，err:", err)
	return
}
name := "沙河小王子"
// 向打开的文件句柄中写入内容
fmt.Fprintf(fileObj, "往文件中写如信息：%s", name)
```
**注意，只要满足**`io.Writer`**接口的类型都支持写入。**

## `Sprint`
`Sprint`系列函数会把传入数据生成并返回一个字符串
```
func Sprint(a ...interface{}) string
func Sprintf(format string, a ...interface{}) string
func Sprintln(a ...interface{}) string

//示例
s1 := fmt.Sprint("沙河小王子")
name := "沙河小王子"
age := 18
s2 := fmt.Sprintf("name:%s,age:%d", name, age)
s3 := fmt.Sprintln("沙河小王子")
fmt.Println(s1, s2, s3)
```
## `Errorf`
`Errorf` 函数根据`format`参数生成格式化字符串并返回一个包含该类型的错误
```
func Errorf(format string, a ...interface{}) error
//通常使用这种方式来自定义错误类型，例如：
err := fmt.Errorf("这是一个错误")
```

# 格式化占位符
`*printf`系列函数都支持`format`格式化参数，划分方便记忆：
## 通用占位符
**占位符**|**说明**
---|---
`%v`|值的默认格式表示
`%+v`|类似`%v`，但输出结构体会添加字段名
`%#v`|值的`Go`语法表示
`%T`|打印值的类型
`%%`|百分号

```
fmt.Printf("%v\n", 100)
fmt.Printf("%v\n", false)
o := struct{ name string }{"小王子"}
fmt.Printf("%v\n", o)
fmt.Printf("%#v\n", o)
fmt.Printf("%T\n", o)
fmt.Printf("100%%\n")

//输出结果如下：
100
false
{小王子}
struct { name string }{name:"小王子"}   //%#v
struct { name string }                  //%T
100%
```
## 布尔型
**占位符**|**说明**
---|---
`%t`|true 或 false

## 整型
**占位符**|**说明**
---|---
`%b`|	表示为二进制
`%c`|	该值对应的unicode码值
`%d`|	表示为十进制
`%o`|	表示为八进制
`%x`|	表示为十六进制，使用a-f
`%X`|	表示为十六进制，使用A-F
`%U`|	表示为Unicode格式：U+1234，等价于”U+%04X”
`%q`|	该值对应的单引号括起来的go语法字符字面值，必要时会采用安全的转义表示

## 浮点数与复数
**占位符**|**说明**
---|---
`%b`|无小数部分，二进制指数的科学计数法
`%e`|科学计数法
`%E`|科学计数法
`%f`|有小数部分但无指数部分
`%F`|等价于`%f`
`%g`|根据实际情况采用`%e`或`%f`格式(以获得更简洁，准确输出)
`%G`|根据实际情况采用`%E`或`%F`格式(以获得更简洁，准确输出)

## 字符串和`[]byte`
**占位符**|**说明**
---|---
`%s`|直接输出字符串或者`[]byte`
`%q`|该值对应的双引号括起来的`go`语法字符串字面值必要时会采用安全转义表示
`%x`|每个字节用两字符十六进制表示(使用a-f)
`%X`|每个字节用两字符十六进制表示(使用A-F)

示例代码如下：
```golang
s := "小王子"
fmt.Printf("%s\n", s)
fmt.Printf("%q\n", s)
fmt.Printf("%x\n", s)
fmt.Printf("%X\n", s)

//输出结果如下：
小王子
"小王子"
e5b08fe78e8be5ad90
E5B08FE78E8BE5AD90
```

## 指针
**占位符**|**说明**
---|---
`%p`|表示为十六进制，并加上前导的`ox`
`%#p`|表示为十六进制

示例代码如下：
```golang
a := 10
fmt.Printf("%p\n", &a)
fmt.Printf("%#p\n", &a)
输出结果如下：

0xc000094000
c000094000
```

## 宽度标识符
如果未指定精度，会使用默认精度；如果点号后没有跟数字，表示精度为0
**占位符**|**说明**
---|---
`%f`|默认宽度，默认精度
`%9f`|宽度9，默认禁锢
`%.2f`|默认宽度，精度2
`%9.2f`|宽度9，精度2
`%9.f`|宽度9，精度0

示例代码如下：
```golang
n := 12.34
fmt.Printf("%f\n", n)
fmt.Printf("%9f\n", n)
fmt.Printf("%.2f\n", n)
fmt.Printf("%9.2f\n", n)
fmt.Printf("%9.f\n", n)
//输出结果如下：

12.340000
12.340000
12.34
    12.34
       12
```

# 获取输入
`fmt`包下有`fmt.Scan, fmt.Scanf, fmt.Scanln`三个函数，可以在程序运行过程中从标准输入获取用户输入

读取输入有三大家族:
* `Scan` : 从标准输入`os.Stdin`读取数据, 包括`Scan(),Scanf(),Scanln()`
* `SScan`: 从字符串读取数据 ...
* `Fscan`: 从`io.Reader`读取数据 ...

## `fmt.Scan`
```golang
func Scan(a ...interface{}) (n int, err error)
```
* `Scan`从标准输入扫描文本，读取由空白符分割的值保存到传递给本函数的参数中，换行符视为空白符
* 本函数返回成功扫描数据个数和错误

```golang
func main() {
	var (
		name    string
		age     int
		married bool
	)
	fmt.Scan(&name, &age, &married)
	fmt.Printf("扫描结果 name:%s age:%d married:%t \n", name, age, married)
}

//终端执行，依次输入 小王子，28，false使用空格分割
$ ./scan_demo 
小王子 28 false
扫描结果 name:小王子 age:28 married:false 
```

`fmt.Scan`从标准输入中扫描用户输入数据，将以空白符分隔的数据分别存入指定参数。

## `fmt.Scanf`
```golang
func Scanf(format string, a ...interface{}) (n int, err error)
```
* `Scanf`从标准输入扫描文本，根据`format`参数指定格式去读取由空白符分隔的值保存到传递给本函数的参数中。
* 本函数返回成功扫描数据个数和错误

```golang
func main() {
	var (
		name    string
		age     int
		married bool
	)
	fmt.Scanf("1:%s 2:%d 3:%t", &name, &age, &married)
	fmt.Printf("扫描结果 name:%s age:%d married:%t \n", name, age, married)
}
```
将上面的代码编译后在终端执行，在终端按照指定的格式依次输入小王子、28和false。

```golang
$ ./scan_demo 
1:小王子 2:28 3:false       //要按照Scanf里面格式输入，空格分隔
扫描结果 name:小王子 age:28 married:false 
fmt.Scanf不同于fmt.Scan简单的以空格作为输入数据的分隔符，fmt.Scanf为输入数据指定了具体的输入内容格式，只有按照格式输入数据才会被扫描并存入对应变量。
```

## `fmt.Scanln`
```golang
func Scanln(a ...interface{}) (n int, err error)
```
* 遇到换行时才停止扫描，最后一个数据后面必须由换行或者到达结束位置
* 返回成功个数和遇到错误

示例代码
```golang
func main() {
	var (
		name    string
		age     int
		married bool
	)
	fmt.Scanln(&name, &age, &married)
	fmt.Printf("扫描结果 name:%s age:%d married:%t \n", name, age, married)
}
//将上面的代码编译后在终端执行，在终端依次输入小王子、28和false使用空格分隔。

$ ./scan_demo 
小王子 28 false
扫描结果 name:小王子 age:28 married:false 
//fmt.Scanln遇到回车就结束扫描了，这个比较常用。
```

## `bufio.NewReader`
获取完整输入，输入内容可能包含空格

示例代码：
```
func bufioDemo() {
	reader := bufio.NewReader(os.Stdin) // 从标准输入生成读对象
	fmt.Print("请输入内容：")
	text, _ := reader.ReadString('\n') // 读到换行
	text = strings.TrimSpace(text)
	fmt.Printf("%#v\n", text)
}
```

## `Fscan`系列
类似于`fmt.Scan、fmt.Scanf、fmt.Scanln`三个函数，只不过不是从标准输入中读取数据而是hi从`io.Reader`中读取数据。
```golang
func Fscan(r io.Reader, a ...interface{}) (n int, err error)
func Fscanf(r io.Reader, format string, a ...interface{}) (n int, err error)
func Fscanln(r io.Reader, a ...interface{}) (n int, err error)
```

## `Sscan`系列
类似于`fmt.Scan、fmt.Scanf、fmt.Scanln`三个函数,只不过是从**指定字符串中读取数据**。
```golang
func Sscan(str string, a ...interface{}) (n int, err error)
func Sscanf(str string, format string, a ...interface{}) (n int, err error)
func Sscanln(str string, a ...interface{}) (n int, err error)
```