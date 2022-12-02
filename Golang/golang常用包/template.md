# `template`

`Go`语言中输出`HTML`场景使用`html/template`包,一般文件使用`text/template`包,都提供了相同的接口

# `Go`语言模板引擎

作用机制简单归纳:

*   模板文件通常定义为`.tmpl`和`.tpl`为后缀（也可以使用其他的后缀），必须使用`UTF8`编码。
*   模板文件中使用`{{和}}`包裹和标识需要传入的数据。
*   传给模板这样的数据就可以通过点号`（.）`来访问，如果数据是复杂类型的数据，可以通过`{{ .FieldName }}`来访问它的字段。
*   除`{{和}}`包裹的内容外，其他内容均不做修改

# 模板引擎使用

## 定义模板文件

修改数据使用\`\`{{ }}\`包裹

## 解析

定义好模板文件后,可以使用下面方法去解析,得到模板对象:

```golang
func (t *Template) Parse(src string) (*Template, error)
func ParseFiles(filenames ...string) (*Template, error)
func ParseGlob(pattern string) (*Template, error)
```

也可以`func New(name string) *Template` 函数创建一个名为`name`的模板

## 模板渲染

```golang
func (t *Template) Execute(wr io.Writer, data interface{}) error
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error
```

# 基本示例

## 定义模板文件

定义`hello.tmpl`模板文件：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Hello</title>
</head>
<body>
    <p>Hello {{.}}</p>
</body>
</html>
```

## 解析和渲染模板

创建`main.go`,写下`HTTP server`端:

```golang
// main.go

func sayHello(w http.ResponseWriter, r *http.Request) {
	// 解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	// 利用给定数据渲染模板，并将结果写入w
	tmpl.Execute(w, "沙河小王子")
}
func main() {
	http.HandleFunc("/", sayHello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Println("HTTP server failed,err:", err)
		return
	}
}
```

# 模板语法

## `{{.}}`

其中`{{.}}`中的点表示当前对象,可以根据`.`访问结构体对应字段

```golang
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", sayHello)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type UserInfo struct {
	Name	string
	Gender	string
	Age		int
}

func sayHello(writer http.ResponseWriter, request *http.Request) {
	//解析指定文件生成模板对象
	tmpl, err := template.ParseFiles("./hello.tmpl")
	if err != nil {
		fmt.Println("create template failed, err: ", err)
		return
	}
	user := UserInfo{
		Name:   "liukaitao",
		Gender: "man",
		Age:    18,
	}
	tmpl.Execute(writer, user)
}

```

模板`hello.tmpl`内容:

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Hello</title>
</head>
<body>
    <p>Hello {{.Name}}</p>
    <p>性别：{{.Gender}}</p>
    <p>年龄：{{.Name}}</p>
</body>
</html>
```

## 注释

    {{/* a comment */}}
    注释,执行时会忽略, 可以多行

## `pipeline`

`pipeline`是指产生数据的操作。比如`{{.}},{{.Name}}`等。

## 变量

    $obj := {{.}}

其中 `$obj`是变量名字,后续代码中就可以使用该变量

## 移除空格

有时候我们在使用模板语法的时候会不可避免的引入一下空格或者换行符，这样模板最终渲染出来的内容可能就和我们想的不一样。

例如：

    {{- .Name -}}

注意：`-`要紧挨`{{和}}`，同时与模板值之间需要使用空格分隔

## 条件判断

    {{if pipeline}} T1 {{end}}

    {{if pipeline}} T1 {{else}} T0 {{end}}

    {{if pipeline}} T1 {{else if pipeline}} T0 {{end}}

## `range`遍历

    {{range pipeline}} T1 {{end}}
    如果pipeline的值其长度为0，不会有任何输出

    {{range pipeline}} T1 {{else}} T0 {{end}}
    如果pipeline的值其长度为0，则会执行T0。

## `with`

    {{with pipeline}} T1 {{end}}
    如果pipeline为empty不产生输出，否则将dot设为pipeline的值并执行T1。不修改外面的dot。

    {{with pipeline}} T1 {{else}} T0 {{end}}
    如果pipeline为empty，不改变dot并执行T0，否则dot设为pipeline的值并执行T1。

## 预定义函数

执行模板时，函数从两个函数字典中查找：首先是模板函数字典，然后是全局函数字典。一般不在模板内定义函数，而是使用Funcs方法添加函数到模板里。

预定义的全局函数如下：

    and
        函数返回它的第一个empty参数或者最后一个参数；
        就是说"and x y"等价于"if x then y else x"；所有参数都会执行；
    or
        返回第一个非empty参数或者最后一个参数；
        亦即"or x y"等价于"if x then x else y"；所有参数都会执行；
    not
        返回它的单个参数的布尔值的否定
    len
        返回它的参数的整数类型长度
    index
        执行结果为第一个参数以剩下的参数为索引/键指向的值；
        如"index x 1 2 3"返回x[1][2][3]的值；每个被索引的主体必须是数组、切片或者字典。
    print
        即fmt.Sprint
    printf
        即fmt.Sprintf
    println
        即fmt.Sprintln
    html
        返回与其参数的文本表示形式等效的转义HTML。
        这个函数在html/template中不可用。
    urlquery
        以适合嵌入到网址查询中的形式返回其参数的文本表示的转义值。
        这个函数在html/template中不可用。
    js
        返回与其参数的文本表示形式等效的转义JavaScript。
    call
        执行结果是调用第一个参数的返回值，该参数必须是函数类型，其余参数作为调用该函数的参数；
        如"call .X.Y 1 2"等价于go语言里的dot.X.Y(1, 2)；
        其中Y是函数类型的字段或者字典的值，或者其他类似情况；
        call的第一个参数的执行结果必须是函数类型的值（和预定义函数如print明显不同）；
        该函数类型值必须有1到2个返回值，如果有2个则后一个必须是error接口类型；
        如果有2个返回值的方法返回的error非nil，模板执行会中断并返回给调用模板执行者该错误；

## 比较函数

布尔函数会将任何类型的零值视为假，其余视为真。

下面是定义为函数的二元比较运算的集合：

    eq      如果arg1 == arg2则返回真
    ne      如果arg1 != arg2则返回真
    lt      如果arg1 < arg2则返回真
    le      如果arg1 <= arg2则返回真
    gt      如果arg1 > arg2则返回真
    ge      如果arg1 >= arg2则返回真

为了简化多参数相等检测，`eq`（只有eq）可以接受2个或更多个参数，它会**将第一个参数和其余参数依次比较**，返回下式的结果：

    {{eq arg1 arg2 arg3}}

比较函数只适用于基本类型

## 自定义函数

```golang
func sayHello(w http.ResponseWriter, r *http.Request) {
	htmlByte, err := ioutil.ReadFile("./hello.tmpl")
	if err != nil {
		fmt.Println("read html failed, err:", err)
		return
	}
	// 自定义一个夸人的模板函数
	kua := func(arg string) (string, error) {
		return arg + "真帅", nil
	}
	// 采用链式操作在Parse之前调用Funcs添加自定义的kua函数
	tmpl, err := template.New("hello").Funcs(template.FuncMap{"kua": kua}).Parse(string(htmlByte))
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}

	user := UserInfo{
		Name:   "小王子",
		Gender: "男",
		Age:    18,
	}
	// 使用user渲染模板，并将结果写入w
	tmpl.Execute(w, user)
}
```

我们可以在模板文件`hello.tmpl`中按照如下方法使用我们自定义的`kua`函数

    {{kua .Name}}

## 嵌套`template`

子`template`可以是单独文件,也可以是通过`define`定义的`template`

`t.tmpl`文件如下:

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>tmpl test</title>
</head>
<body>
    <h1>测试嵌套template语法</h1>
    <hr>
    {{template "ul.tmpl"}}
    <hr>
    {{template "ol.tmpl"}}
</body>
</html>

{{ define "ol.tmpl"}}
<ol>
    <li>吃饭</li>
    <li>睡觉</li>
    <li>打豆豆</li>
</ol>
{{end}}
```

`ul.tmpl`文件如下:

```html
<ul>
    <li>注释</li>
    <li>日志</li>
    <li>测试</li>
</ul>
```

注册路由处理函数:

```golang
http.HandleFunc("/tmpl", tmplDemo)

func tmplDemo(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./t.tmpl", "./ul.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	user := UserInfo{
		Name:   "小王子",
		Gender: "男",
		Age:    18,
	}
	tmpl.Execute(w, user)
}
```

## `block`

    {{block "name" pipeline}} T1 {{end}}

`block`是定义模板`{{define "name"}} T1 {{end}}`和执行`{{template "name" pipeline}}`缩写，典型的用法是定义一组根模板，然后通过在其中重新定义块模板进行自定义。

定义一个基础模板`templates/base.tmpl`:

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <title>Go Templates</title>
</head>
<body>
<div class="container-fluid">
    {{block "content" . }}{{end}}
</div>
</body>
</html>
```

然后定义一个`templates/index.tmpl`,继承`base.tmpl`

```html
{{template "base.tmpl"}}

{{define "content"}}
    <div>Hello world!</div>
{{end}}
```

然后使用`template.ParseGlob`按照正则匹配规则解析模板文件，然后通过`ExecuteTemplate`渲染指定的模板：

```golang
func index(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		fmt.Println("create template failed, err:", err)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.tmpl", nil)
	if err != nil {
		fmt.Println("render template failed, err:", err)
		return
	}
}
```

## 修改默认标识符

```golang
template.New("test").Delims("{[", "]}").ParseFiles("./t.tmpl")
```

## 取消跨站脚本攻击防护

```golang
{{ . | safe }}
```

