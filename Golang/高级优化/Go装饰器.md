- [装饰器模式](#装饰器模式)
- [普通类型](#普通类型)
- [使用反射传入`接口类型参数`](#使用反射传入接口类型参数)
- [通过反射生成新的函数覆盖装饰](#通过反射生成新的函数覆盖装饰)
- [多个装饰器共同使用](#多个装饰器共同使用)
- [使用带有参数的装饰器](#使用带有参数的装饰器)

# 装饰器模式

# 普通类型

```go
package main

import (
	"fmt"
	"reflect"
)

// help 用于测试
func help(s string) string {
	fmt.Println("s: ", s)
	return s
}

// // decorator 装饰器
func decorator(fn func(name string) string) func(string) string {
	fmt.Println("define in decorator")
	return func(name string) string {
		fmt.Println("func in decorator")
		res := fn(name)
		return res
	}
}

// Go 实现装饰器
func main() {
	test := decorator(help)
	fmt.Println(test("Hello"))
}
```

普通`decorator`传入参数为`func`，`注意传入参数和返回参数`

返回参数也是`func`，相当于返回函数本身

# 使用反射传入`接口类型参数`

```go
// help 用于测试
func help(s string) string {
	fmt.Println("s: ", s)
	return s
}

// decorator　装饰器　使用interface
func decorator(fn interface{}, params ...interface{}) func() interface{} {
	fmt.Println("come in decorator")
	funcValue := reflect.ValueOf(fn)
	// 获取参数个数
	paramsLen := funcValue.Type().NumIn()
	if funcValue.Kind() != reflect.Func {					// 判断是否函数
		fmt.Errorf("fn must be a function")					// 这里返回错误
	}

	paramsList := make([]reflect.Value, paramsLen)

	if len(params) != paramsLen {
		fmt.Errorf("fn params length is not enough")	// 这里返回错误
	}

	for i, item := range params {
		paramsList[i] = reflect.ValueOf(item)		
	}

	return func() interface{} {
		return funcValue.Call(paramsList)
	}
}

func main() {
	res := decorator(help, "1")
	fmt.Println(res())
}

```

反射:`funcValue := reflect.ValueOf(fn)`

`decorator`传入`fn interface{}`，利用反射\`\`funcValue.Kind() == reflect.Func\`判断是否函数；

传入不定长参数`params ...interface{}`，利用反射`funcValue.Type().NumIn()`判断参数长度;

最后对`funcValue.Call(paramsList)`进行反射调用

# 通过反射生成新的函数覆盖装饰

```go
package main

import (
	"fmt"
	"reflect"
)

// help 用于测试
func help(s string) string {
	fmt.Println("s: ", s)
	return s
}

// decorator 装饰器
func decorator(deco, fn interface{}) {
	var decoratedFunc, targetFunc reflect.Value

	decoratedFunc = reflect.ValueOf(deco).Elem()
	targetFunc = reflect.ValueOf(fn)

	v := reflect.MakeFunc(targetFunc.Type(), // 通过反射生成一个新的函数
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("decorator self code")
            // in 就是传入参数，　out就是输出参数
			out = targetFunc.Call(in)
			return
		})

	decoratedFunc.Set(v)
	return
}

func main() {
	testFunc := help
	decorator(&testFunc, help)
	// 这里testFunc已经被改变
	testFunc("hah")
}

```

`reflect.ValueOf(deco).Elem()`　反射获取原本函数对应本身，作为输出的装饰完毕之后的函数，取代本身

`reflect.ValueOf(fn)`动态获取函数方法

使用`reflect.MakeFunc`返回`type`的函数，经过`后面func`包装的

最后把`decoratedFunc`的值给更换掉

# 多个装饰器共同使用

```go
package main

import (
	"fmt"
)

// help 用于测试
func help(s string) {
	fmt.Println("s: ", s)
}

type Func func(string)

type Decorator func(Func) Func

func decorator1(f Func) Func {
	return func(name string) {
		fmt.Println("decorator 1")
		f(name)
		fmt.Println("decorator 1")
	}
}

func decorator2(f Func) Func {
	return func(name string) {
		fmt.Println("decorator 2")
		f(name)
		fmt.Println("decorator 2")
	}
}

// decorator 装饰器
func handler(h Func, decors ...Decorator) Func { // 使用Decorator数组实现
	for i := range decors {
		d := decors[len(decors)-1-i] // 让后面的先包装，前面的函数会先调用
		h = d(h)
	}
	return h

}

func main() {
	// 经过装饰器装饰
	help := handler(help, decorator1, decorator2)
	help("haha")
}

```

定义函数类型，方便阅读

`handler`传入被装饰函数`h`和装饰器`...Decorator`，不定长参数

使用遍历包装`d := decors[len(decors) -1 -i ]    ,   h= d(h)`

最后覆盖自身`help := handler(help, decorator1, decorator2)`

# 使用带有参数的装饰器

```go
type Func func(string) // 定义函数类型方便使用

// 装饰器生成器, 传入参数为first
func decoratorGen(first string) func(Func) Func {
	return func(f Func) Func {
		return func(name string) {
			fmt.Println("first: ", first)
			f(name)
		}
	}
    
   
func main() {
	// 经过装饰器装饰
	decorator := decoratorGen("first")
    // 获取到装饰器，执行装饰调用
	help := decorator(help)
	help("nihao")
}

```

