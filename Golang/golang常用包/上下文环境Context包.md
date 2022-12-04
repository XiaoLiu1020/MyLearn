
- [`Go`标准库`Context`](#go标准库context)
- [为什么需要`Context`](#为什么需要context)
	- [全局变量方式](#全局变量方式)
	- [通道方式](#通道方式)
	- [官方版方案 `context`](#官方版方案-context)
- [`Context`](#context)
	- [`Context`接口](#context接口)
	- [`Background()和TODO()`](#background和todo)
	- [`With`系列函数](#with系列函数)
		- [`WithCancel`](#withcancel)
		- [`WithDeadline`](#withdeadline)
		- [`WithTimeout`](#withtimeout)
		- [`WithValue`](#withvalue)
		- [使用`Context`注意事项](#使用context注意事项)
- [客户端超时取消示例](#客户端超时取消示例)
	- [`server`端](#server端)
	- [`client`](#client)

# `Go`标准库`Context`
在`Go http`包的`Server`中,每一个请求都有一个自己独立的对应的`goroutine`去处理; 当一个请求被取消或者超时时,所有用来处理该请求的`goroutine`都应该退出,然后系统才会释放这些`goroutine`占用资源

# 为什么需要`Context`
需求: 优雅可控制地结束子`goroutine`
## 全局变量方式
```golang
package main

import (
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup
var exit bool

// 全局变量方式存在的问题：
// 1. 使用全局变量在跨包调用时不容易统一
// 2. 如果worker中再启动goroutine，就不太好控制了。

func worker() {
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		if exit {
			break
		}
	}
	wg.Done()
}

func main() {
	wg.Add(1)
	go worker()
	time.Sleep(time.Second * 3) // sleep3秒以免程序过快退出
	exit = true                 // 修改全局变量实现子goroutine的退出
	wg.Wait()
	fmt.Println("over")
}
```
## 通道方式
```golang
package main

import (
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup

// 管道方式存在的问题：
// 1. 使用全局变量在跨包调用时不容易实现规范和统一，需要维护一个共用的channel

func worker(exitChan chan struct{}) {
LOOP:
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		select {
		case <-exitChan: // 等待接收上级通知
			break LOOP
		default:
		}
	}
	wg.Done()
}

func main() {
	var exitChan = make(chan struct{})
	wg.Add(1)
	go worker(exitChan)
	time.Sleep(time.Second * 3) // sleep3秒以免程序过快退出
	exitChan <- struct{}{}      // 给子goroutine发送退出信号
	close(exitChan)
	wg.Wait()
	fmt.Println("over")
}
```
## 官方版方案 `context`
```golang
package main

import (
	"context"
	"fmt"
	"sync"

	"time"
)

var wg sync.WaitGroup

func worker(ctx context.Context) {
LOOP:
	for {
		fmt.Println("worker")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done(): // 等待上级通知
			break LOOP
		default:
		}
	}
	wg.Done()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 3)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
```

# `Context`
`context`为新标准库,定义了`context`类型,专门简化对于处理单个请求的多个`goroutine`之间与请求域的数据,取消信息,截止时间等相关操作.这些操作可能涉及多个`API`调用。

对服务器传入的请求应该创建上下文，而对服务器的传出调用应该接受上下文。它们之间的函数调用链必须传递上下文，或者可以使用`WithCancel、WithDeadline、WithTimeout或WithValue`创建的派生上下文。当一个上下文被取消时，它派生的所有上下文也被取消。

## `Context`接口
`context.Context`是一个结构,定义了四个需要实现的方法.
```golang
type Context interface{
    Deadline() (deadline time.Time, ok bool)    //返回当前context被取消时间
    Done() <-chan struct{}  //返回一个channel, 这个Channel会在当前工作完成或者上下文被取消之后关闭
    Err() error //返回当前context结束原因, Done返回的channel被关闭时才会返回非空值
    Value(key interface{}) interface{}  //从context中返回键对应值
}
```
## `Background()和TODO()`
`Go`内置两个函数：`Background()和TODO()`，这两个函数分别返回一个实现了`Context`接口的`background`和`todo`。我们代码中最开始都是以这两个内置的上下文对象作为最顶层的`parent context`，衍生出更多的子上下文对象。

`Background()`主要用于`main`函数、初始化以及测试代码中，作为Context这个树结构的最顶层的Context，也就是根Context。

`TODO()`，它目前还不知道具体的使用场景，如果我们不知道该使用什么Context的时候，可以使用这个。

background和todo本质上都是emptyCtx结构体类型，是一个不可取消，没有设置截止时间，没有携带任何值的Context。

## `With`系列函数

### `WithCancel`
```golang
func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
```
`WithCancel`返回带新`Done`通道的父节点副本, 当调用返回的cancel函数或当关闭父上下文的Done通道时，将关闭返回上下文的Done通道，无论先发生什么情况。

取消此上下文将释放与其关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel。
```golang
func gen(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return // return结束该goroutine，防止泄露
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 当我们取完需要的整数后调用cancel, 为ctx.Done()发送信号

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
```
`gen`函数在单独的`goroutine`中生成整数并将它们发送到返回的通道。 `gen`的调用者在使用生成的整数之后需要取消上下文，以免`gen`启动的内部`goroutine`发生泄漏。

### `WithDeadline`
```golang
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
```
当截止日过期时，当调用返回的`cancel`函数时，或者当父上下文的`Done`通道关闭时，返回上下文的`Done`通道将被关闭，以最先发生的情况为准。
```golang
func main() {
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)

	// 尽管ctx会过期，但在任何情况下调用它的cancel函数都是很好的实践。
	// 如果不这样做，可能会使上下文及其父类存活的时间超过必要的时间。
	defer cancel()

	select {
	case <-time.After(1 * time.Second): //一秒后触发
		fmt.Println("overslept")
	case <-ctx.Done():      //过时触发channel
		fmt.Println(ctx.Err())
	}
}
```
`ctx.Done()`会先接受到值,打印`ctx`退出原因

### `WithTimeout`
```golang
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
```
`WithTimeout`返回`WithDeadline(parent, time.Now().Add(timeout))`。
```golang
package main

import (
	"context"
	"fmt"
	"sync"

	"time"
)

// context.WithTimeout

var wg sync.WaitGroup

func worker(ctx context.Context) {
LOOP:
	for {
		fmt.Println("db connecting ...")
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func main() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
```

### `WithValue`
```golang
func WithValue(parent Context, key, val interface{}) Context
```
`WithValue`返回父节点的副本，其中与key关联的值为`val`。

仅对`API`和进程间传递请求域的数据使用上下文值，而不是使用它来传递可选参数给函数。

**所提供的键必须是可比较的，并且不应该是**`string`**类型或任何其他内置类型，以避免使用上下文在包之间发生冲突。**`WithValue`**的用户应该为键定义自己的类型。**

```golang
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// context.WithValue

type TraceCode string       //自定义TraceCode类型

var wg sync.WaitGroup

func worker(ctx context.Context) {
	key := TraceCode("TRACE_CODE")
	traceCode, ok := ctx.Value(key).(string) // 在子goroutine中获取trace code
	// 使用类型推断
	if !ok {
		fmt.Println("invalid trace code")
	}
    LOOP:
	for {
		fmt.Printf("worker, trace code:%s\n", traceCode)
		time.Sleep(time.Millisecond * 10) // 假设正常连接数据库耗时10毫秒
		select {
		case <-ctx.Done(): // 50毫秒后自动调用
			break LOOP
		default:
		    continute
		}
	}
	fmt.Println("worker done!")
	wg.Done()
}

func main() {
	// 设置一个50毫秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
	// 在系统的入口中设置trace code传递给后续启动的goroutine实现日志数据聚合
	// 为ctx传入key, value
	ctx = context.WithValue(ctx, TraceCode("TRACE_CODE"), "12512312234")
	wg.Add(1)
	go worker(ctx)
	time.Sleep(time.Second * 5)
	cancel() // 通知子goroutine结束
	wg.Wait()
	fmt.Println("over")
}
```

### 使用`Context`注意事项
* 推荐以参数的方式显示传递`Context`
* 以`Context`作为参数的函数方法，应该把`Context`作为第一个参数。
* 给一个函数方法传递`Context`的时候，不要传递`nil`，如果不知道传递什么，就使用`context.TODO()`
* `Context`的`Value`相关方法应该传递请求域的必要数据，**不应该用于传递可选参数**
* `Context`是线程安全的，**可以放心的在多个goroutine中传递**


# 客户端超时取消示例
## `server`端
```golang
package main

import (
	"fmt"
	"math/rand"
	"time"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		panic(err)
	}
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	number := rand.Intn(2)  //随机数
	if number == 0 {
		time.Sleep(time.Second * 5)
		fmt.Fprintf(writer, "slow response")
		fmt.Printf("finish slow response")
	}
	fmt.Fprintf(writer, "quick response")
}
```

## `client`
```golang
package main

import (
	"context"
	"golang.org/x/tools/tools/go/ssa/interp/testdata/src/fmt"
	"golang.org/x/tools/tools/godoc/analysis"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	//定义超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond * 2)
	defer cancel()
	doCall(ctx)
}

type respData struct {
	resp *http.Response
	err error
}

func doCall(ctx context.Context) {
	transport := http.Transport{
		//请求频繁可定义全局的client对象并启用长链接
		//请求不频繁使用短链接
		DisableKeepAlives: true,
	}
	//定义client对象
	client := http.Client{
		Transport: &transport,
	}
	//定义通道，存放*respData 类型, struct
	respChan := make(chan *respData, 1)
	req, err := http.NewRequest("GET", "http://127.0.0.1:8000/", nil)
	if err != nil {
		fmt.Printf("new request failed, err: %v\n", err)
		return
	}
	// 使用带超时的ctx创建一个新的client request
	req = req.WithContext(ctx)
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()
	go func() {
		resp, err := client.Do(req)	//请求
		fmt.Printf("client.do resp:%v, err:%v\n", resp, err)
		rd := &respData{
			resp: resp,
			err: err,
		}
		respChan <- rd		//传递请求数据
		wg.Done()
	}()
	//不需要for循环，看谁先发来信号满足case，就运行哪一个
	select {
	case <-ctx.Done():			//超时ctx会发送信息
		//transport.CancelRequest(req)
		fmt.Println("call api timeout")
	case result := <-respChan:
		fmt.Println("call server api success")
		if result.err != nil {
			fmt.Printf("call server api failed, err: %v\n", err)
			return
		}
		data, _ := ioutil.ReadAll(result.resp.Body)
		fmt.Printf("resp: %v\n", string(data))
		result.resp.Body.Close()
	}
}

```
