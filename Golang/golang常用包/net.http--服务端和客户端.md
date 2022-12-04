- [`net/http`](#nethttp)
	- [HTTP服务端](#http服务端)
		- [默认`Server`示例](#默认server示例)
		- [自定义`Server`](#自定义server)
		- [自定义`Handler`](#自定义handler)
- [`Go`代码执行流程](#go代码执行流程)


# `net/http`
内置 `net/http`包提供`HTTP`客户端和服务端的实现

## HTTP服务端
`ListenAndServe` 使用指定监听地址和处理器启动一个`HTTP`服务端, 处理器参数通常是`nil`, --表示采用包变量`DefaultServeMux` 作为处理器

可以使用`Handle`和`HandleFunc`函数 向`DefaultServeMux`添加处理器
```golang
http.Handle("/foo", fooHandler)
http.HandleFunc("/bar", func(writer http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(writer, "hello, %q", html.EscapeString(r.URL.Path))
})
// 传入的是"/bar"路径, 和func

log.Fatal(http.ListenAndServe(":8080", nil))    //监听端口, nil默认DefaultServeMux
```

### 默认`Server`示例
```golang
// http server

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello 沙河！")
}

func main() {
	http.HandleFunc("/", sayHello)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Printf("http server failed, err:%v\n", err)
		return
	}
}
```
### 自定义`Server`
```golang
s := &http.Server{
	Addr:           ":8080",
	Handler:        myHandler,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
```

### 自定义`Handler`
需要实现自己的`ServerHttp`方法
```golang
package main

import (
	"fmt"
	"log"
	"net/http"
)

type dollars float32
// String()方法的类型，默认输出的时候会调用该方法，实现字符串的打印。
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }		

type MyHandler map[string]dollars
func (self *MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {		//实现了 ServeHTTP方法
	switch req.URL.Path {
	case "/list":
		for item, price := range self {
		fmt.Fprintf(w, "%s: %s\n", item, price)
			}
	case "/price":
		item := req.URL.Query().Get("item")
		price, ok := self[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound) // 404
			fmt.Fprintf(w, "no such item: %q\n", item)
			return
			}
		fmt.Fprintf(w, "%s\n", price)
	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}

func main() {
	handler := &MyHandler{"shoes": 50, "socks": 5}		//handler 需要实现自己的ServerHttp方法
	log.Fatal(http.ListenAndServe("localhost:8000", handler))
}

简单版本：
package main

import (
	"fmt"
	"net/http"
)

type MyHandler struct {}

func (mh *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		sayhelloworld(w, r)
		return
	}
	http.NotFound(w, r)
}

func sayhelloworld(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HelloWorld!")
	fmt.Println("Hello world")
}

func main() {
	myhandler := &MyHandler{}
	http.ListenAndServe(":9090", myhandler)
}

```

# `Go`代码执行流程
* 首先调用`Http.HandleFunc`

按顺序：
1. 调用`DefaultServeMux`的`HandleFunc` ，`HandleFunc`是类型，实现了`ServeHTTP`方法
2. 调用`DefaultServeMux`的`Handle`
3. 往`DefaultServeMux`的`map[string]muxEntry`中增加对应`handler`和`pattern`路由规则

* 其次调用`http.ListenAndServe(":9090", nil)

1. 实例化`Server`
2. 调用`Server`的`ListenAndServe()`
3. 启动`for`循环，`Accept`请求
4. 对每个请求实例化一个`Conn`并且开启一个`goroutine`进行服务, `go c.serve()`
5. 读取每个请求内容 `w ,err := c.readRequest()`
6. 判断`handler`是否为空，如果为空，没有设置，`handler`就设置为`DefaultServeMux`
7. 调用`handler`的`ServeHttp`， 即进入到`DefaultServeMux.ServeHttp`中
8. 根据`request`选择`handler`,并且进入到`handler`的`ServeHTTP`：`mux.handler(r).ServeHTTP(w, r)`
9. 选择`handler`:
10. 判断是否有路由满足这个`request`(循环遍历`ServeMux`的`muxEntry`)
11. 看情况调用这个路由`handler`的`ServeHTTP`或者`NotFoundHandler`的`ServeHTTP`