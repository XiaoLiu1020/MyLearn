- [`Go`语言中并发编程](#go语言中并发编程)
- [`goroutine`](#goroutine)
	- [使用`goroutine`](#使用goroutine)
	- [启动单个 `goroutine`](#启动单个-goroutine)
	- [启动多个`goroutine`](#启动多个goroutine)
- [`goroutine`与线程](#goroutine与线程)
	- [可增长的栈](#可增长的栈)
	- [`goroutine`调度](#goroutine调度)
	- [`GOMAXPROCS`](#gomaxprocs)
- [`channel` 通信](#channel-通信)
	- [`channel` 类型](#channel-类型)
	- [`channel` 操作](#channel-操作)
	- [无缓冲的通道](#无缓冲的通道)
	- [有缓冲的通道](#有缓冲的通道)
		- [`for range` 循环取值](#for-range-循环取值)
	- [单向通道](#单向通道)
	- [通道总结：](#通道总结)
- [`worker pool (goroutine)`池](#worker-pool-goroutine池)
	- [`select` 多路复用](#select-多路复用)
	- [`select`设置超时](#select设置超时)
- [并发安全和锁](#并发安全和锁)
	- [互斥锁](#互斥锁)
	- [读写互斥锁](#读写互斥锁)


# `Go`语言中并发编程
`Go`语言并发通过 `goroutine`实现。`goroutine` 类似于线程，属于用户态的线程；

可以创建成千上万个`goroutine`并发工作，`goroutine`由`Go`语言运行时(runtime)调度完成，而线程由操作系统调度完成

`Go`语言还提供 `channel` 在多个 `goroutine`间通信。

# `goroutine`
引言：我们要实现并发编程的时候，我们通常需要自己维护一个线程池，并且需要自己去包装一个又一个的任务，同时需要自己去调度线程执行任务并维护上下文切换，这一切通常会耗费程序员大量的心智。那么能不能有一种机制，**程序员只需要定义很多个任务，让系统去帮助我们把这些任务分配到CPU上实现并发执行呢？**

`Go`程序会智能将 `goroutine`中任务合理分配给每个`CPU`, 因为它**在语言层面已经内置了调度和上下文切换的机制**。

在`Go`语言编程中你不需要去自己写进程、线程、协程，你的技能包里只有一个技能–`goroutine`，**当你需要让某个任务并发执行的时候，你只需要把这个任务包装成一个函数，开启一个**`goroutine`**去执行这个函数就可以了，就是这么简单粗暴**。

## 使用`goroutine`
只需要在调用函数时候在前面加上 `go`关键字，就可以为一个函数创建一个 `goroutine`。
## 启动单个 `goroutine`
```
func hello() {
    fmt.Println("Hello Goroutine!")
}
func main() {
    // 串行执行方式
    hello()
    fmt.Println("main goroutine done!" )  
    // 启动并发     
    go hello()      //启动另外一个goroutine 去执行hello函数
    fmt.Println("main goroutine done!")
}
```
并发执行时候：`main()`函数创建一个默认的 `goroutine`，**相当于主进程**，当 `main()`函数返回时候该`goroutine`就结束了，它里面的 `goroutine`也一起结束，所以里面开启的子 `goroutine`不一定运行完；

## 启动多个`goroutine`
这里也使用了 `sync.WaitGroup` 来实现 `goroutine`的同步

```
var wg sync.WaitGroup

func hello(i int) {
    defer wg.Done()     // goroutine结束就登记-1, defer 最后执行语法
    fmt.Println("Hello Goroutine!", i)
}

func main() {
    for i :=0; i < 10; i++ {
        wg.Add(1)       //启动一个goroutine就登记+1
        go hello(i)
    }
    wg.Wait()           //等待所有登记的goroutine都结束
}
```
10个 `goroutine` 是并发执行的，`goroutine`调度随机的，每次打印数字顺序都不一致

---

# `goroutine`与线程
## 可增长的栈
1. `os`线程 一般有固定的栈内存(通常为`2MB`),而一个`goroutine`的栈在声明周期开始时只有很小(典型情况`2KB`)
2. `goroutine` 栈会按需增大喝缩小，大小限制可以达到`1GB`

## `goroutine`调度
`GPM`是`Go`语言运行时`（runtime）`层面的实现，是`go`语言自己实现的一套调度系统。区别于操作系统调度OS线程。
* `G` - `goroutine` 存放本`goroutine`信息还有与所在 `P`的绑定信息；
* `P` `Pipeline` 管理一组`goroutine`队列，会存储当前 `goroutine`**运行的上下文环境(函数指针，堆栈地址及地址边界)**, `P`会对自己管理的 `goroutine`队列做一些调度， 当自己队列消费完了就去全局队列取，全局队列消费完会去其他 `P`队列抢任务
* `M(machine)`是 `Go`运行时`(runtime)`对操作系统内核线程的虚拟，一个`goroutine` 最终是要放到 `M`上执行的

`P`与`M`一般也是一一对应的，他们关系是： `P`管理着一组`G`挂载在`M`上运行，当一个`G`长久阻塞在一个`M`上时，`runtime`会新建一个`M`，阻塞`G`所在的`P`会把其他的`G`挂载在新建的`M`上，当旧的`G`阻塞完成或者认为已经死掉了，就会回收旧的`M`。


单从线程调度讲，Go语言相比起其他语言的优势在于OS线程是由OS内核来调度的，goroutine则是由Go运行时（runtime）自己的调度器调度的，这个调度器使用一个称为m:n调度的技术（复用/调度m个goroutine到n个OS线程）。 其一大特点是goroutine的调度是在用户态下完成的， 不涉及内核态与用户态之间的频繁切换，包括内存的分配与释放，都是在用户态维护着一块大的内存池， 不直接调用系统的malloc函数（除非内存池需要改变），成本比调度OS线程低很多。 另一方面充分利用了多核的硬件资源，近似的把若干goroutine均分在物理线程上， 再加上本身goroutine的超轻量，以上种种保证了go调度方面的性能。

## `GOMAXPROCS`
`Go`运行时的调度器使用`GOMAXPROCS`参数来确定需要使用多少个OS线程来同时执行Go代码。默认值是机器上的`CPU核心数`。例如在一个8核心的机器上，调度器会把Go代码同时调度到8个OS线程上（`GOMAXPROCS`是`m:n`调度中的`n`）。
Go语言中可以通过`runtime.GOMAXPROCS()`函数设置当前程序并发时占用的CPU逻辑核心数。

```
func a() {
	for i := 1; i < 10; i++ {
		fmt.Println("A:", i)
	}
}

func b() {
	for i := 1; i < 10; i++ {
		fmt.Println("B:", i)
	}
}

func main() {
    // 如果设定为1，此时会做完一个任务再做另一个任务
    // 如果设定为2，两个任务会并行执行
	runtime.GOMAXPROCS(1)
	go a()
	go b()
	time.Sleep(time.Second)
}
```
**Go语言中操作系统线程和goroutine的关系:**
1. 一个操作系统线程对应用户态多个`goroutine`
2. `go`程序可以同时使用多个`OS`线程
3. `goroutine` 和 `OS`线程时多对多关系，即`m:n`

# `channel` 通信
Go语言的并发模型是`CSP（Communicating Sequential Processes）`，提倡**通过通信共享内存而不是通过共享内存而实现通信。**

如果说`goroutine`是`Go`程序并发的执行体，`channel`就是它们之间的连接。`channel`是可以让一个`goroutine`发送特定值到另一个`goroutine`的通信机制。

`Go `语言中的通道`（channel`）是一种特殊的类型。通道像一个传送带或者队列，总是遵循先入先出（`First In First Out`）的规则，保证收发数据的顺序。每一个通道都是一个具体类型的导管，也就是**声明channel的时候需要为其指定元素类型**

## `channel` 类型
`channel`是一种引用类型

```
// 语法
var 变量 chan 元素类型

var ch1 chan int   // 声明一个传递整型的通道
var ch2 chan bool  // 声明一个传递布尔型的通道
var ch3 chan []int // 声明一个传递int切片的通道
```

## `channel` 操作
通道是引用类型，`channel` 空值是 `nil`
```
var ch chan int
fmt.Println(ch) // <nil>
```
声明之后需要初始化，使用make
```
make(chan 元素类型, [缓冲大小])     // 缓冲大小可选


ch4 := make(chan int)
ch5 := make(chan bool)
ch6 := make(chan []int)

```

通道有发送`（send）`、接收`(receive）`和关闭`（close）`三种操作。

发送和接收都使用 `<-` 符号。

```
channel := make(chan int)

//发送

channel <- 10    //发送10到channel

//接收

x := <- channel // 从channel中接收并赋值到 x中
<- channel           // 从channel接收值，忽略结果

//关闭
close(channel)
```
* 只有在通知接收方`goroutine`所有的数据都发送完毕的时候才需要关闭通道。
* `channel`可以被垃圾回收机制回收，所以不是必须关闭
    * 对已经关闭的`channel`再 发送会引发 `panic`
    * 对一个关闭的通道进行 `<- channel` 接收，会一直获取直到通道为空
    * 关闭且没有值的通道执行接收**会得到对应类型的 零值**

## 无缓冲的通道
也称为阻塞通道
```
func main() {
    ch := make(chan int)    //没输入缓冲大小，所以无缓冲
    ch <- 10
    fmt.Println("发送成功")
}
// 报错
fatal error: all goroutines are asleep - deadlock!

```
无缓冲的通道必须有接收才能发送
```
// 可以启用一个 goroutine`去接受值
func recv(c chan int) {
    ret := <- c
    fmt.Println("接收成功", ret)
}
func main() {
    ch := make(chan int)
    go recv(ch)     // 启用goroutine从通道接收值, 传入ch作为参数
    ch <- 10
    fmt.Println("发送成功")
}
```
**无缓冲通道上的发送操作会阻塞**，直到另一个`goroutine`在该通道上执行接收操作，这时两个 `goroutine`才继续执行。相反也是有个阻塞等待。

无缓冲通道会使`goroutine`同步化，也可以称为`同步通道`

## 有缓冲的通道
```
func main() {
    channel := make(chan int, 1)    //创建容量为1的有缓冲通道
    ch <- 10
    fmt.Println("发送成功")
}
```
通道满了装不下时，就会产生阻塞，等到别人接收，才能继续发送

使用内置 `len()`函数获取通道内元素数量， `cap()`函数获取通道容量

### `for range` 循环取值
当通道关闭时，再发送或者再关闭会引发 `panic` ，那如何判断一个通道是否被关闭了呢？
```
// channel 练习
func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	// 开启goroutine将0~100的数发送到ch1中
	go func() {
		for i := 0; i < 100; i++ {
			ch1 <- i
		}
		close(ch1)
	}()
	// 开启goroutine从ch1中接收值，并将该值的平方发送到ch2中
	go func() {
		for {
			i, ok := <-ch1 // 通道关闭后再取值ok=false
			if !ok {
				break
			}
			ch2 <- i * i
		}
		close(ch2)
	}()
	// 在主goroutine中从ch2中接收值打印
	for i := range ch2 { // 通道关闭后会退出for range循环
		fmt.Println(i)
	}
}
```
`for range` 遍历通道，当通道被关闭时候就会退出 `for range` 循环

## 单向通道
限制通道在函数中只能发送或只能接收，相当于生产者或消费者
```
// 生产者，传入只发送的通道
func counter(out chan<- int) {
    for i := 0; i < 100; i++ {
        out <- i        //发送
    }
    close(out)
}

// 计算平方，接收之后，发送
func squarer(out chan<- int, in <-chan int) {
    for i := range in { //接收in
        out <- i * i    //发送
    }
    close(out)
}

func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)
    go counter(ch1)     //传入通道ch1,在里面只发送
    go squarer(ch2, ch1)    //ch2在里面负责接收，ch1负责发送
    printer(ch2)
}
```
其中：
* `chan <- int` 是一个只能发送的通道
* `<- chan int` 是一个只能接收的通道
* 在函数传参及任何赋值操作中是可以**将双向通道转换为单向通道是可以的**，但**反过来不可以**

## 通道总结：
`channel`异常总结：
![异常总结](https://www.liwenzhou.com/images/Go/concurrence/channel01.png)

# `worker pool (goroutine)`池
可以指定启动的 `goroutine`数量- `worker pool`模式，控制`goroutine`数量，防止`goroutine`泄露和暴涨

```
package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <- chan int, res chan <- int) {
	for i := range jobs {
		fmt.Printf("worker: %d start job:%d \n", id, i)
		time.Sleep(time.Second)
		fmt.Printf("worker: %d end job:%d \n", id, i)
		res <- i * 2
	}
}

func main() {
	jobs := make(chan int, 100)
	res := make(chan int, 100)
	//开启3个goroutine, 等待发送
	for j := 1; j <= 3; j++ {
		go worker(j, jobs, res)
	}
	// 开启5个任务，发送到jobs中
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)
	// 遍历接收通道结果
	for i := 1; i <= 5; i++ {
		<- res
	}
	close(res)

}

//输出
/*
worker: 3 start job:3
worker: 1 start job:1
worker: 2 start job:2
worker: 1 end job:1
worker: 1 start job:4
worker: 3 end job:3
worker: 3 start job:5
worker: 2 end job:2
worker: 1 end job:4
worker: 3 end job:5
*/
```

## `select` 多路复用
某些场景下需要同时从多个通道接收数据，通道接收数据时，如果没有数据可以接收将会发生阻塞；

如果使用遍历每个通道方式实现，运行性能会差很多，`Go`内置了`select`关键字，同时响应多个通道操作

`select`使用类似 `switch`语句，有一系列 `case`分支和一个`default`分支，每个`case`会对应一个通道通信过程。 `select`会一致等待，直到某个`case`通信操作完成就会执行`case`分支对应语句。

```
select {
    case <- ch1:
        ...
    case data := <- ch2:
        ...
    case ch3 <- data:
        ...
    default:
        默认操作
}

func main() {
    channel := make(chan int, 1)
    for i :=0; i < 10; i++ {
        select {
        case x := <- channel:
            fmt.Println(x)
        case channale <- i:
        }
    }
}
```
使用 `select`语句能提高代码可读性：
* 可处理一个或多个 `channel`的发送/接收操作
* 如果多个 `case` 同时满足， `select`会随机选择一个
* 对于没有 `case`的 `select{}`会一致等待，可用于阻塞`main`函数

## `select`设置超时
```golang
func main() {
	c := make(chan int)
	o := make(chan bool)
	go func() {
		for {
			select {
				case v := <- c:
					println(v)
				case <- time.After(5 * time.Second):
					println("timeout")
					o <- true
					break
			}
		}
	}()
	<- o
}
```

---

# 并发安全和锁
可能会存在多个 `goroutine` 同时操作一个资源，产生 `竞态问题`(数据竞态)

例子：
```
var x int64
var wg sync.WaitGroup

func add() {
	for i := 0; i < 5000; i++ {
		x = x + 1
	}
	wg.Done()
}
func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)
}
```
上面的代码中我们开启了两个`goroutine`去累加变量`x`的值，这两个`goroutine`在访问和修改`x`变量的时候就会存在数据竞争,导致最后结果与期待不符合

## 互斥锁
保证同时只有一个 `goroutine`可以访问共享资源

使用 `sync`包的`Mutex`类型实现

```
var x int64
var wg sync.WaitGroup
var lock sync.Mutex

func add() {
	for i := 0; i < 5000; i++ {
		lock.Lock() // 加锁
		x = x + 1
		lock.Unlock() // 解锁
	}
	wg.Done()
}
func main() {
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	fmt.Println(x)
}
```
使用互斥锁能够保证同一时间有且只有一个`goroutine`进入临界区，其他的`goroutine`则在等待锁；当互斥锁释放后，等待的`goroutine`才可以获取锁进入临界区，多个`goroutine`同时等待一个锁时，唤醒的策略为随机。

## 读写互斥锁
读多写少情况下，**当我们并发去读取一个资源不涉及资源修改时候是没有必要加锁的**，因为读写锁是一个合理选择,使用 `sync`包中的`RWMutex`

读写锁`RWMutex`分读锁和写锁：
1. 当一个 `goroutine`获取读锁后，其他的`goroutine`如果是获取读锁会可以继续获得锁，**如果是获取写锁就会等待**
2. 当一个`goroutine`获取写锁之后，其他的`goroutine`无论是获取写或读锁都会等待。

```

var (
    x       int64
    wg      sync.WaitGroup
    lock    sync.Mutex
    rwlock  sync.RWMutex
)

func write() {
	// lock.Lock()   // 加互斥锁
	rwlock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond) // 假设读操作耗时10毫秒
	rwlock.Unlock()                   // 解写锁
	// lock.Unlock()                     // 解互斥锁
	wg.Done()
}

func read() {
	// lock.Lock()                  // 加互斥锁
	rwlock.RLock()               // 加读锁
	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
	rwlock.RUnlock()             // 解读锁
	// lock.Unlock()                // 解互斥锁
	wg.Done()
}

func main() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()			//写锁阻塞会产生阻塞，10*10 = 100ms
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go read()			//读锁，同样读不影响，是并发
	}

	wg.Wait()
	end := time.Now()
	fmt.Println(end.Sub(start))
}
```