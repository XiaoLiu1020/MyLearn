- [问题：](#问题)
- [解决方案](#解决方案)
- [控制`goroutine`并发数量](#控制goroutine并发数量)
  - [尝试`chan` + `sync`](#尝试chan--sync)
  - [使用信号量`Semaphore`](#使用信号量semaphore)
      - [优势](#优势)
- [灵活控制`goroutine`并发数量](#灵活控制goroutine并发数量)
      - [优势](#优势-1)
- [方案三,第三方库](#方案三第三方库)

# 问题：
```golang
func main() {
    userCount := math.MaxInt64
    for i := 0; i < userCount; i++ {
        go func(i int) {
            // 做一些各种各样的业务逻辑处理
            fmt.Printf("go func: %d\n", i)
            time.Sleep(time.Second)
        }(i)
    }
}
```
**上面会开启多个**`goroutine`，使系统资源占用率不断上涨， 输出一定数量后，不刷新输出最新值，发出信号`signal:killed`

不控制并发的 `goroutine` 数量 会发生什么问题？大致如下：

* `CPU` 使用率浮动上涨
* `Memory` 占用不断上涨。也可以看看 `CMPRS`，它表示进程的压缩数据的字节数。已经到达 `114G+` 了
* 主进程崩溃（被杀掉了）

# 解决方案
* 控制/限制`goroutine`同时并发运行的数量
* 改变应用程序的逻辑写法(避免大规模的使用系统资源和等待)

# 控制`goroutine`并发数量
## 尝试`chan` + `sync`
```golang
...
var wg = sync.WaitGroup{}

func main() {
    userCount := 10
    ch := make(chan bool, 2)
    for i := 0; i < userCount; i++ {
        wg.Add(1)
        go Read(ch, i)
    }

    wg.Wait()
}

func Read(ch chan bool, i int) {
    defer wg.Done()

    ch <- true
    fmt.Printf("go func: %d, time: %d\n", i, time.Now().Unix())
    time.Sleep(time.Second)
    <-ch
}
```
## 使用信号量`Semaphore`
```golang
/*
信号量，控制协程同时并发处理数量，方法有Acquire,获取许可，AcquireWithTime，指定时间内获取许可
releases释放
*/
package main

import "time"

type Semaphore struct {
	permits int			// amount of the acquire
	channel chan int 	// acquire or release semaphore
}
// 创建信号量
func New(p int) *Semaphore {
	return &Semaphore{
		permits: p,
		channel: make(chan int, p),
	}
}

//获取许可, 当channel满了，会进行阻塞
func (s *Semaphore) Acquire() {
	s.channel <- 0
}

// 释放许可，channel长度减一
func (s *Semaphore) Release() {
	<- s.channel
}

// 尝试获取许可
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.channel <- 0:
		return true
	default:
		return false
	}
}

// 尝试指定时间内获取许可
func (s *Semaphore) TryAcquireOnTime(timeout time.Duration) bool{
	for {
		select {
		case s.channel <- 0:
			return true
		case <-time.After(timeout):
			return false
		}
	}
}

// 当前可用许可书
func (s *Semaphore) AvailablePermits() int {
	return s.permits - len(s.channel)
}
```

使用信号量，并发数量达到一定程度就可以达到阻塞`goroutine`目的，也可以使用`TryAcquire`方法，直接返回`bool`

#### 优势
* 适合量不大、复杂度低的使用场景
    * 几百几千个、几十万个也是可以接受的（看具体业务场景）
    * 实际业务逻辑在运行前就已经被阻塞等待了（因为并发数受限），基本实际业务逻辑损耗的性能比 goroutine 本身大
    * goroutine 本身很轻便，仅损耗极少许的内存空间和调度。这种等待响应的情况都是躺好了，等待任务唤醒
* Semaphore 操作复杂度低且流转简单，容易控制

# 灵活控制`goroutine`并发数量
要控制输入的数量，以此达到改变允许并发运行 `goroutine` 的数量

```golang
package main

import (
    "fmt"
    "sync"
    "time"
)

var wg sync.WaitGroup

func main() {
    userCount := 10     //开启10个gc
    ch := make(chan int, 5)     //缓冲通道只有5个
    for i := 0; i < userCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for d := range ch {
                fmt.Printf("go func: %d, time: %d\n", d, time.Now().Unix())
                time.Sleep(time.Second * time.Duration(d))
            }
        }()
    }

    for i := 0; i < 10; i++ {
        ch <- 1
        ch <- 2
        //time.Sleep(time.Second)
    }

    close(ch)       //关闭通道
    wg.Wait()
}
```

#### 优势
* 变更 `channel` 的输入数量
* 能够根据特殊情况，变更 `channel` 的循环值
* 变更最大允许并发的 `goroutine` 数量

# 方案三,第三方库

* `go-playground/pool`
* `nozzle/throttler`
* `Jeffail/tunny`
* `panjf2000/ants`
比较成熟的第三方库也不少，基本都是以生成和管理 `goroutine` 为目标的池工具。