- [`Go`性能优化之路](#go性能优化之路)
- [1. `go tool pprof`](#1-go-tool-pprof)
	- [1.2 `runtime/pprof`](#12-runtimepprof)
	- [1.3 `net/http/pprof`---`httpserver 类型`](#13-nethttppprof---httpserver-类型)
		- [如果是`httpserver`使用`go-gin`包](#如果是httpserver使用go-gin包)
		- [1.4 `pprof CPU`分析例子](#14-pprof-cpu分析例子)
- [2 `pprof`使用方式](#2-pprof使用方式)
	- [2.1 使用`Web`界面](#21-使用web界面)
	- [2.2 通过交互式终端使用](#22-通过交互式终端使用)
	- [2.3 通过`可视化界面`](#23-通过可视化界面)
		- [启动 PProf 可视化界面](#启动-pprof-可视化界面)
		- [查看 PProf 可视化界面](#查看-pprof-可视化界面)
		- [查看火焰图](#查看火焰图)
- [3. 内存分析](#3-内存分析)
	- [3.1 使用系统`top`](#31-使用系统top)
	- [3.2 `GODEBUG与gctrace`　跟踪`GC`内存释放情况](#32-godebug与gctrace跟踪gc内存释放情况)
		- [涉及术语](#涉及术语)
		- [**格式**](#格式)
		- [**含义**](#含义)
		- [`STW`-`STOP the World`](#stw-stop-the-world)
	- [3.3 `runtime.ReadMemStats`－－运行方法查看](#33-runtimereadmemstats运行方法查看)
	- [3.4 `pprof`工具查看](#34-pprof工具查看)
- [4. `go tool trace`](#4-go-tool-trace)
- [5. `go test -bench` 校验](#5-go-test--bench-校验)
- [性能优化之路实践](#性能优化之路实践)
	- [1. 使用`sync.Pool`复用对象](#1-使用syncpool复用对象)
	- [2. 使用成员变量复用对象](#2-使用成员变量复用对象)
	- [3. 写时复制代替互斥锁](#3-写时复制代替互斥锁)
	- [4. 避免包含指针结构体作为`map`的`key`](#4-避免包含指针结构体作为map的key)
	- [5. 使用`strings.Builder`拼接字符串](#5-使用stringsbuilder拼接字符串)
	- [6. 分区：减少共享数据结构争夺](#6-分区减少共享数据结构争夺)
- [ 实践--参考go-zero部分](#-实践--参考go-zero部分)


# `Go`性能优化之路

![](https://img-blog.csdnimg.cn/img_convert/be9664275b677260a9cc5bd4a0eb60f8.webp?x-oss-process=image/format,png)

# 1. `go tool pprof`

采集自：　<https://zhuanlan.zhihu.com/p/33528194>

`Golang` 提供的两个官方包 [runtime/pprof](https://link.zhihu.com/?target=https%3A//golang.org/pkg/runtime/pprof/)，[net/http/pprof](https://link.zhihu.com/?target=https%3A//golang.org/pkg/net/http/pprof/) .能方便的采集程序运行的堆栈、`goroutine`、内存分配和占用、`io `等信息的 `.prof` 文件

使用` go tool pprof` 分析 `.prof` 文件

## 1.2 `runtime/pprof`

如果程序为`非 httpserver 类型`，使用此方式；在 main 函数中嵌入如下代码:

```go
import "runtime/pprof"

var cpuprofile = flag.String("cpuprofile"， ""， "write cpu profile `file`")
var memprofile = flag.String("memprofile"， ""， "write memory profile to `file`")

func main() {
    flag.Parse()
    if *cpuprofile != "" {
        f， err := os.Create(*cpuprofile)
        if err != nil {
            log.Fatal("could not create CPU profile: "， err)
        }
        if err := pprof.StartCPUProfile(f); err != nil {
            log.Fatal("could not start CPU profile: "， err)
        }
        defer pprof.StopCPUProfile()
    }

    // ... rest of the program ...

    if *memprofile != "" {
        f， err := os.Create(*memprofile)
        if err != nil {
            log.Fatal("could not create memory profile: "， err)
        }
        runtime.GC() // get up-to-date statistics
        if err := pprof.WriteHeapProfile(f); err != nil {
            log.Fatal("could not write memory profile: "， err)
        }
        f.Close()
    }
}
```

运行程序

```bash
./logger -cpuprofile cpu.prof -memprofile mem.prof
```

可以得到 cpu.prof 和 mem.prof 文件，使用 go tool pprof 分析。

```bash
go tool pprof logger cpu.prof
go tool pprof logger mem.prof
```

## 1.3 `net/http/pprof`---`httpserver 类型`

如果程序为 `httpserver 类型`， 则只需要导入该包:

```go
import _ "net/http/pprof"
```

### 如果是`httpserver`使用`go-gin`包

而不是使用默认的 `http` 包启动，则需要手动添加 `/debug/pprof` 对应的 handler，`github `有[封装好的模版](https://github.com/DeanThompson/ginpprof):

```go
import "github.com/DeanThompson/ginpprof"
...
router := gin.Default()
ginpprof.Wrap(router)
...
```

导入包重新编译程序后运行,在浏览器中访问 `http://host:port/debug/`，可以看到性能信息

通过浏览器查看的数据不能直观反映程序性能问题，`go tool pprof` 命令行工具提供了丰富的工具集:

查看 heap 信息

```bash
go tool pprof http://127.0.0.1:4500/debug/pprof/heap
```

查看 30s 的 CPU 采样信息

```bash
go tool pprof http://127.0.0.1:4500/debug/pprof/profile
```

其他功能使用参见 [官方 net/http/pprof 库](https://link.zhihu.com/?target=https%3A//golang.org/pkg/net/http/pprof/)

### 1.4 `pprof CPU`分析例子

采集 profile 数据之后，可以分析 CPU 热点代码。 执行下面命令：

```bash
go tool pprof http://127.0.0.1:4500/debug/pprof/profile
```

会采集 30s 的 profile 数据，之后进入终端交互模式，输入 `top` 指令。

```bash
~ # go tool pprof http://127.0.0.1:4500/debug/pprof/profile
Fetching profile over HTTP from http://127.0.0.1:4500/debug/pprof/profile
Saved profile in /home/vagrant/pprof/pprof.logger.samples.cpu.012.pb.gz
File: logger
Type: cpu
Time: Jan 19， 2018 at 2:01pm (CST)
Duration: 30s， Total samples = 390ms ( 1.30%)
Entering interactive mode (type "help" for commands， "o" for options)
(pprof) top
Showing nodes accounting for 360ms， 92.31% of 390ms total
Showing top 10 nodes out of 74
      flat  flat%   sum%        cum   cum%
     120ms 30.77% 30.77%      180ms 46.15%  compress/flate.(*compressor).findMatch /usr/local/go/src/compress/flate/deflate.go
     100ms 25.64% 56.41%      310ms 79.49%  compress/flate.(*compressor).deflate /usr/local/go/src/compress/flate/deflate.go
      60ms 15.38% 71.79%       60ms 15.38%  compress/flate.matchLen /usr/local/go/src/compress/flate/deflate.go
      20ms  5.13% 76.92%       20ms  5.13%  compress/flate.(*huffmanBitWriter).indexTokens /usr/local/go/src/compress/flate/huffman_bit_writer.go
      10ms  2.56% 79.49%       10ms  2.56%  compress/flate.(*huffmanBitWriter).writeTokens /usr/local/go/src/compress/flate/huffman_bit_writer.go
      10ms  2.56% 82.05%       10ms  2.56%  hash/adler32.update /usr/local/go/src/hash/adler32/adler32.go
      10ms  2.56% 84.62%       10ms  2.56%  runtime.futex /usr/local/go/src/runtime/sys_linux_amd64.s
      10ms  2.56% 87.18%       10ms  2.56%  runtime.memclrNoHeapPointers /usr/local/go/src/runtime/memclr_amd64.s
      10ms  2.56% 89.74%       10ms  2.56%  runtime.pcvalue /usr/local/go/src/runtime/symtab.go
      10ms  2.56% 92.31%       10ms  2.56%  runtime.runqput /usr/local/go/src/runtime/runtime2.go
(pprof)
```

# 2 `pprof`使用方式

暂时使用如下demo

```go
//demo.go
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof"
    "github.com/EDDYCJY/go-pprof-example/data"
)

func main() {
    go func() {
        for {
            log.Println(data.Add("https://github.com/EDDYCJY"))
        }
    }()

    http.ListenAndServe("0.0.0.0:6060", nil)
}

// data/d.go
package data

var datas []string

func Add(str string) string {
    data := []byte(str)
    sData := string(data)
    datas = append(datas, sData)

    return sData
}

```

运行这个文件，你的 HTTP 服务会多出 /debug/pprof 的 endpoint 可用于观察应用程序的情况

## 2.1 使用`Web`界面

查看当前总览：访问 `http://127.0.0.1:6060/debug/pprof/`

```cpp
/debug/pprof/

profiles:
0   block
5   goroutine
3   heap
0   mutex
9   threadcreate

full goroutine stack dump
```

*   cpu（CPU Profiling）: `$HOST/debug/pprof/profile`，默认进行 30s 的 CPU Profiling，得到一个分析用的 profile 文件
*   block（Block Profiling）：`$HOST/debug/pprof/block`，查看导致阻塞同步的堆栈跟踪
*   goroutine：`$HOST/debug/pprof/goroutine`，查看当前所有运行的 goroutines 堆栈跟踪
*   heap（Memory Profiling）: `$HOST/debug/pprof/heap`，查看活动对象的内存分配情况
*   mutex（Mutex Profiling）：`$HOST/debug/pprof/mutex`，查看导致互斥锁的竞争持有者的堆栈跟踪
*   threadcreate：`$HOST/debug/pprof/threadcreate`，查看创建新OS线程的堆栈跟踪

## 2.2 通过交互式终端使用

（1）go tool pprof <http://localhost:6060/debug/pprof/profile?seconds=60>

```bash
$ go tool pprof http://localhost:6060/debug/pprof/profile\?seconds\=60

Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=60
Saved profile in /Users/eddycjy/pprof/pprof.samples.cpu.007.pb.gz
Type: cpu
Duration: 1mins, Total samples = 26.55s (44.15%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) 
```

执行该命令后，需等待 60 秒（可调整 seconds 的值），pprof 会进行 CPU Profiling。结束后将默认进入 pprof 的交互式命令模式，可以对分析的结果进行查看或导出。具体可执行 `pprof help` 查看命令说明

```bash
(pprof) top10
Showing nodes accounting for 25.92s, 97.63% of 26.55s total
Dropped 85 nodes (cum <= 0.13s)
Showing top 10 nodes out of 21
      flat  flat%   sum%        cum   cum%
    23.28s 87.68% 87.68%     23.29s 87.72%  syscall.Syscall
     0.77s  2.90% 90.58%      0.77s  2.90%  runtime.memmove
     0.58s  2.18% 92.77%      0.58s  2.18%  runtime.freedefer
     0.53s  2.00% 94.76%      1.42s  5.35%  runtime.scanobject
     0.36s  1.36% 96.12%      0.39s  1.47%  runtime.heapBitsForObject
     0.35s  1.32% 97.44%      0.45s  1.69%  runtime.greyobject
     0.02s 0.075% 97.51%     24.96s 94.01%  main.main.func1
     0.01s 0.038% 97.55%     23.91s 90.06%  os.(*File).Write
     0.01s 0.038% 97.59%      0.19s  0.72%  runtime.mallocgc
     0.01s 0.038% 97.63%     23.30s 87.76%  syscall.Write
```

*   `flat`：给定函数上运行耗时
*   `flat%`：同上的 CPU 运行耗时总比例
*   `sum%`：给定函数累积使用 CPU 总比例
*   `cum`：当前函数加上它之上的调用运行总耗时
*   `cum%`：同上的 CPU 运行耗时总比例

最后一列为函数名称，在大多数的情况下，我们可以通过这五列得出一个应用程序的运行情况，加以优化 🤔

（2）go tool pprof <http://localhost:6060/debug/pprof/heap>

```bash
$ go tool pprof http://localhost:6060/debug/pprof/heap
Fetching profile over HTTP from http://localhost:6060/debug/pprof/heap
Saved profile in /Users/eddycjy/pprof/pprof.alloc_objects.alloc_space.inuse_objects.inuse_space.008.pb.gz
Type: inuse_space
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 837.48MB, 100% of 837.48MB total
      flat  flat%   sum%        cum   cum%
  837.48MB   100%   100%   837.48MB   100%  main.main.func1
```

*   \-inuse\_space：分析应用程序的常驻内存占用情况
*   \-alloc\_objects：分析应用程序的内存临时分配情况

（3） go tool pprof <http://localhost:6060/debug/pprof/block>

（4） go tool pprof <http://localhost:6060/debug/pprof/mutex>

## 2.3 通过`可视化界面`

需要测试用例,因为监控的是运行中采集的数据

```go
package data

import "testing"

const url = "https://github.com/EDDYCJY"

func TestAdd(t *testing.T) {
    s := Add(url)
    if s == "" {
        t.Errorf("Test.Add error!")
    }
}

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(url)
    }
}
```

（2）执行测试用例

```bash
$ go test -bench=. -cpuprofile=cpu.prof
pkg: github.com/EDDYCJY/go-pprof-example/data
BenchmarkAdd-4      10000000           187 ns/op
PASS
ok      github.com/EDDYCJY/go-pprof-example/data    2.300s
```

\-memprofile 也可以了解一下

### 启动 PProf 可视化界面

方法一：

```bash
$ go tool pprof -http=:8080 cpu.prof
```

方法二：

```bash
$ go tool pprof cpu.prof 
$ (pprof) web
```

需要安装`graphviz`

***

### 查看 PProf 可视化界面

（1）Top

![img](https://upload-images.jianshu.io/upload_images/12294974-6394cba1b0a00696.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200)

（2）Graph

![img](https://upload-images.jianshu.io/upload_images/12294974-9154d8ef2970217b.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200)

框越大，线越粗代表它占用的时间越大哦

（3）Peek

![img](https://upload-images.jianshu.io/upload_images/12294974-9b7bcea7f44e2029.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200)

（4）Source

![img](https://upload-images.jianshu.io/upload_images/12294974-876571115869e640.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200)

### 查看火焰图

（1） 安装 PProf

```bash
$ go get -u github.com/google/pprof
```

（2） 启动 PProf 可视化界面:

```bash
$ pprof -http=:8080 cpu.prof
```

​	(3) 查看

![img](https://upload-images.jianshu.io/upload_images/12294974-0a076fdc295db7aa.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/1200)

它就是本次的目标之一，它的最大优点是动态的。调用顺序由上到下\*\*（A -> B -> C -> D）\*\*，每一块代表一个函数，越大代表占用 CPU 的时间更长。同时它也支持点击块深入进行分析！

# 3. 内存分析

## 3.1 使用系统`top`

```go
package main

import (
    "log"
    "runtime"
    "time"
)

func test() {
    //slice 会动态扩容，用slice来做堆内存申请
    container := make([]int, 8)

    log.Println(" ===> loop begin.")
    for i := 0; i < 32*1000*1000; i++ {
        container = append(container, i)
    }
    log.Println(" ===> loop end.")
}

func main() {
    log.Println("Start.")

    test()

    log.Println("force gc.")
    runtime.GC() //强制调用gc回收

    log.Println("Done.")

    time.Sleep(3600 * time.Second) //睡眠，保持程序不退出
}
```

编译运行

```bash
$go build -o snippet_mem && ./snippet_mem
```

使用`top`命令

```bash
$top -p $(pidof snippet_mem)
```

## 3.2 `GODEBUG与gctrace`　跟踪`GC`内存释放情况

直接对可执行文件添加变量

执行`snippet_mem`程序之前添加环境变量`GODEBUG='gctrace=1'`来跟踪打印垃圾回收器信息

```bash
$ GODEBUG='gctrace=1' ./snippet_mem
```

设置`gctrace=1`会使得垃圾回收器在每次回收时汇总所回收内存的大小以及耗时，
并将这些内容汇总成单行内容打印到标准错误输出中。

### 涉及术语

*   `mark`：标记阶段。
*   `markTermination`：标记结束阶段。
*   `mutator assist`：辅助 GC，是指在 GC 过程中 mutator 线程会并发运行，而 mutator assist 机制会协助 GC 做一部分的工作。
*   `heap_live`：在 Go 的内存管理中，span 是内存页的基本单元，每页大小为 8kb，同时 Go 会根据对象的大小不同而分配不同页数的 span，而 heap\_live 就代表着所有 span 的总大小。
*   `dedicated / fractional / idle`：在标记阶段会分为三种不同的 mark worker 模式，分别是 dedicated、fractional 和 idle，它们代表着不同的专注程度，其中 dedicated 模式最专注，是完整的 GC 回收行为，fractional 只会干部分的 GC 行为，idle 最轻松

### **格式**

```bash
gc # @#s #%: #+#+# ms clock, #+#/#/#+# ms cpu, #->#-># MB, # MB goal, # P
```

### **含义**

> `gc # `       GC次数的编号，每次GC时递增
> `@#s `        距离程序开始执行时的时间
> `#%  `        GC占用的执行时间百分比
> `#+...+#`     GC使用的时间
> `#->#-># MB ` GC开始，结束，以及当前活跃堆内存的大小，单位M
>
> `MB goal`   全局堆内存大小
>
> `P`         使用processor的数量

如果每条信息最后，以`(forced)`结尾，那么该信息是由`runtime.GC()`调用触发

例子

    gc 17 @0.149s 1%: 0.004+36+0.003 ms clock, 0.009+0/0.051/36+0.006 ms cpu, 181->181->101 MB, 182 MB goal, 2 P

该条信息含义如下：

*   `gc 17`: Gc 调试编号为17

*   `@0.149s`:此时程序已经执行了0.149s

*   `1%`: 0.149s中其中gc模块占用了1%的时间

*   `0.004+36+0.003 ms clock`: 垃圾回收的时间，分别为STW（stop-the-world）清扫的时间+并发标记和扫描的时间+STW标记的时间

*   `0.009+0/0.051/36+0.006 ms cpu`: 垃圾回收占用cpu时间

*   `181->181->101 MB`：GC开始前堆内存181M， GC结束后堆内存181M，当前活跃的堆内存101M

*   `182 MB goal`: 全局堆内存大小

*   `2 P`: 本次GC使用了2个P(调度器中的Processer)

### `STW`-`STOP the World`

Stop The World（STW），STW 代指在执行某个垃圾回收算法的某个阶段时，需要将整个应用程序暂停去处理 GC 相关的工作事项。

<https://eddycjy.gitbook.io/golang/di-9-ke-gong-ju/godebug-gc>

## 3.3 `runtime.ReadMemStats`－－运行方法查看

```go
package main

import (
    "log"
    "runtime"
    "time"
)

func readMemStats() {

    var ms runtime.MemStats

    runtime.ReadMemStats(&ms)

    log.Printf(" ===> Alloc:%d(bytes) HeapIdle:%d(bytes) HeapReleased:%d(bytes)", ms.Alloc, ms.HeapIdle, ms.HeapReleased)
}

func test() {
    //slice 会动态扩容，用slice来做堆内存申请
    container := make([]int, 8)

    log.Println(" ===> loop begin.")
    for i := 0; i < 32*1000*1000; i++ {
        container = append(container, i)
        if ( i == 16*1000*1000) {
            readMemStats()
        }
    }

    log.Println(" ===> loop end.")
}

func main() {
    log.Println(" ===> [Start].")

    readMemStats()
    test()
    readMemStats()

    log.Println(" ===> [force gc].")
    runtime.GC() //强制调用gc回收

    log.Println(" ===> [Done].")
    readMemStats()

    go func() {
        for {
            readMemStats()
            time.Sleep(10 * time.Second)
        }
    }()

    time.Sleep(3600 * time.Second) //睡眠，保持程序不退出
}
```

封装了一个函数`readMemStats()`，这里面主要是调用`runtime`中的`ReadMemStats()`方法获得内存信息，然后通过`log`打印出来。

运行发现`runtime.GC()`强制运行GC回收内存

```bash
$ go run demo2.go
2020/03/02 18:21:17  ===> [Start].
2020/03/02 18:21:17  ===> Alloc:71280(bytes) HeapIdle:66633728(bytes) HeapReleased:66600960(bytes)
2020/03/02 18:21:17  ===> loop begin.
2020/03/02 18:21:18  ===> Alloc:132535744(bytes) HeapIdle:336756736(bytes) HeapReleased:155721728(bytes)
2020/03/02 18:21:38  ===> loop end.
2020/03/02 18:21:38  ===> Alloc:598300600(bytes) HeapIdle:609181696(bytes) HeapReleased:434323456(bytes)
2020/03/02 18:21:38  ===> [force gc].
2020/03/02 18:21:38  ===> [Done].
2020/03/02 18:21:38  ===> Alloc:55840(bytes) HeapIdle:1207427072(bytes) HeapReleased:434266112(bytes)
2020/03/02 18:21:38  ===> Alloc:56656(bytes) HeapIdle:1207394304(bytes) HeapReleased:434266112(bytes)
2020/03/02 18:21:48  ===> Alloc:56912(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
2020/03/02 18:21:58  ===> Alloc:57488(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
2020/03/02 18:22:08  ===> Alloc:57616(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(bytes)
c2020/03/02 18:22:18  ===> Alloc:57744(bytes) HeapIdle:1207394304(bytes) HeapReleased:1206493184(by
```

可以看到，打印`[Done].`之后那条trace信息，Alloc已经下降，即内存已被垃圾回收器回收。在`2020/03/02 18:21:38`和`2020/03/02 18:21:48`的两条trace信息中，HeapReleased开始上升，即垃圾回收器把内存归还给系统。

## 3.4 `pprof`工具查看

`pprof`工具支持网页上查看内存的使用情况，需要在代码中添加一个协程即可。

跟 # 1.3 类似

添加以下代码

```go
 import(
    "net/http"
    _ "net/http/pprof"
)
 
 //启动pprof
    go func() {
        log.Println(http.ListenAndServe("0.0.0.0:10000", nil))
    }()
```

# 4. `go tool trace`

可以继续沿用`go tool pprof`的demo

运行以下命令开启跟踪`trace`，参数seconds为跟踪时间，保存在`trace.out`文件中

```bash
curl -o trace.out http://127.0.0.1:6060/debug/pprof/trace?seconds=10
```

对`trace.out`运行

```bash
go tool trace trace.out
```

# 5. `go test -bench` 校验

```makefile
ver1:
	go test -bench=. -count=10 | tee ver1.txt

ver2:
	go test -bench=. -count=10 | tee ver2.txt

benchstat:
	benchstat ver1.txt ver2.txt

```

对比两个版本

安装`benchstat`: `go get golang.org/x/perf/cmd/benchstat`

# 性能优化之路实践

go-performance-code <https://github.com/first-giver/go-performance-code>

`Go pprof 与线上事故`：一次成功的定位与失败的复现 <https://mp.weixin.qq.com/s/c6fU9t951Mv167Ivsy8iXA>

## 1. 使用`sync.Pool`复用对象

本质：　定期进行GC处理的用户定义的对象列表

原理：　复用已经分配的对象，减少分配数量，降低GC压力

*   `必须重置被复用对象`
*   保证`使用后放回池中`，与任何手动内存管理方案一样

```go
package no3_syncpool

import (
	"sync"
	"testing"
)

type Book struct {
	Title    string
	Author   string
	Pages    int
	Chapters []string
}

var pool = sync.Pool{
	New: func() interface{} {
		return &Book{}
	},
}

func BenchmarkNoPool(b *testing.B) {
	var book *Book

	for n := 0; n < b.N; n++ {
		book = &Book{
			Title:  "The Art of Computer Programming, Vol. 1",
			Author: "Donald E. Knuth",
			Pages:  672,
		}
	}

	_ = book
}

func BenchmarkPool(b *testing.B) {
	for n := 0; n < b.N; n++ {
        // 	重置被复用对象
		book := pool.Get().(*Book)
		book.Title = "The Art of Computer Programming, Vol. 1"
		book.Author = "Donald E. Knuth"
		book.Pages = 672
		// 保证使用后放回池子中，无论什么情况
		pool.Put(book)
	}
}

```

运行基准测试

```bash
go test -bench=.* 
```

典型实例

*   利用`sync.Pool`实现接受`UDP`请求的数据`buf`缓冲区－－－－避免`[]byte频繁分配和释放`
*   `updPool`自身可以作为全局变量，更好方式实现为`Server`中的成员变量

```go
var udpPool = sync.Pool{
	New : func () interface {
		return make([]byte, defaultUDPBufferSize)
	}
}

func EchoUDP(address string) error {
	for {
		// 丛缓冲池取出
		buf := udpPool.Get().([]byte)
		// u 为udp的套接字
		num, addr, err := u.ReadFrom(buf)
		if err != nil {
			// 记得释放
			udpPool.Put(buf[:defaultUDPBufferSize])
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			}
			return err
		}
		go handleUDP(u, buf[:num], addr)
	}
}

func handleUDP(u *net.UDPConn, buf []byte, addr net.Addr){
	_, err := u.WriteTo(buf, addr)
	if err != nil {}
	udpPool.Put(buf[:defaultUDPBufferSize])
	return 
}
```

## 2. 使用成员变量复用对象

典型实例：

TCP服务端，将每个`buf`缓冲区和对应`TCP Conn`绑定，每次此`Conn`读取数据都会复用此`buf`

根本目的：避免频繁创建对象

```go
type framer struct{}

// ReadFrame 从io reader拆分出完整数据帧
func (f *framer) ReadFrame(reader io.Reader) (msgbuf []byte, err error) {
    head := make([]byte, framerHeadLen) // 创建
    _, err = io.ReadFull(reader, head[:frameHeadLen])
    totalLen := binary.BigEndian.Uint32(head[4:8])
    msg := make([]byte, totalLen) 				// 创建
    copy(msg, head[:])
    _, err = io.ReadFull(reader, msg[frameHeadLen : totalLen])
}
```

改进后

```go
type framer struct{
    reader io.Reader
    head	[16]byte
    msg		[]byte
}

// ReadFrame 从io reader拆分出完整数据帧
func (f *framer) ReadFrame() (msgbuf []byte, err error) {
    var num int
    _, err = io.ReadFull(reader, f.head[:frameHeadLen])
    totalLen := binary.BigEndian.Uint32(f.head[4:8])
    
    if int(totalLen) > len(f.msg) {
        f.msg = make([]byte, totalLen)
    }

    copy(f.msg, f.head[:])
    num,  err = io.ReadFull(f.reader,f. msg[frameHeadLen : totalLen])
    return f.msg[:totalLen], nil
}
```

## 3. 写时复制代替互斥锁

应用场景：受保护的数据不会经常被修改，并且可以对其进行复制

实现：使用`atomic.Value`保证数据加载和存储操作`原子性`

`atomic_map.go`

```go
package cow

import (
	"sync"
	"sync/atomic"
)

// AtomicMap 原子Map实现
//
// 利用atomic.Value原子(无锁)的加载数据
// 利用Copy-on-Write实现数据更新
type AtomicMap struct {
	mu    sync.Mutex
	clean atomic.Value
}

func (m *AtomicMap) Load(key interface{}) (interface{}, bool) {
	data, _ := m.clean.Load().(map[interface{}]interface{})
	v, ok := data[key]
	return v, ok
}

func (m *AtomicMap) Store(key, value interface{}) {
	m.mu.Lock()
	dirty := m.dirty()
	dirty[key] = value
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *AtomicMap) dirty() map[interface{}]interface{} {
	data, _ := m.clean.Load().(map[interface{}]interface{})
	dirty := make(map[interface{}]interface{}, len(data)+1)

	for k, v := range data {
		dirty[k] = v
	}
	return dirty
}

func (m *AtomicMap) LoadOrStore(key, value interface{}) (interface{}, bool) {
	data, _ := m.clean.Load().(map[interface{}]interface{})
	v, ok := data[key]
	if ok {
		return v, ok
	}

	m.mu.Lock()
	// Lock阻塞获取锁期间,可能数据已经存在，再次Load检查数据
	data, _ = m.clean.Load().(map[interface{}]interface{})
	v, ok = data[key]
	if !ok {
		dirty := m.dirty()
		dirty[key] = value
		v = value
		m.clean.Store(dirty)
	}
	m.mu.Unlock()
	return v, ok
}

func (m *AtomicMap) Delete(key interface{}) {
	m.mu.Lock()
	dirty := m.dirty()
	delete(dirty, key)
	m.clean.Store(dirty)
	m.mu.Unlock()
}

func (m *AtomicMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	data, _ := m.clean.Load().(map[interface{}]interface{})
	for k, v := range data {
		if !f(k, v) {
			break
		}
	}
}

```

`rwmutex_map.go`

```go
package cow

import "sync"

type RWMutexMap struct {
	mu    sync.RWMutex
	dirty map[interface{}]interface{}
}

func (m *RWMutexMap) Load(key interface{}) (interface{}, bool) {
	m.mu.RLock()
	value, ok := m.dirty[key]
	m.mu.RUnlock()
	return value, ok
}

func (m *RWMutexMap) Store(key, value interface{}) {
	m.mu.Lock()
	if m.dirty == nil {
		m.dirty = make(map[interface{}]interface{})
	}
	m.dirty[key] = value
	m.mu.Unlock()
}

func (m *RWMutexMap) LoadOrStore(key, value interface{}) (actual interface{}, loaded bool) {
	m.mu.Lock()
	actual, loaded = m.dirty[key]
	if !loaded {
		actual = value
		if m.dirty == nil {
			m.dirty = make(map[interface{}]interface{})
		}
		m.dirty[key] = value
	}
	m.mu.Unlock()
	return actual, loaded
}

func (m *RWMutexMap) Delete(key interface{}) {
	m.mu.Lock()
	delete(m.dirty, key)
	m.mu.Unlock()
}

func (m *RWMutexMap) Range(f func(key, value interface{}) (shouldContinue bool)) {
	m.mu.RLock()
	keys := make([]interface{}, 0, len(m.dirty))
	for k := range m.dirty {
		keys = append(keys, k)
	}
	m.mu.RUnlock()

	for _, k := range keys {
		v, ok := m.Load(k)
		if !ok {
			continue
		}
		if !f(k, v) {
			break
		}
	}
}

```

![image-20200507153915901](/home/lkt/桌面/notes/pprof.assets/image-20200507153915901.png)

## 4. 避免包含指针结构体作为`map`的`key`

原理：**在垃圾回收期间，　运行时`runtime`扫描包含指针的对象，并进行追踪**

优化方案：　需要在插入`map`之前将字符串散列为整数

![image-20200507155508169](/home/lkt/桌面/notes/pprof.assets/image-20200507155508169.png)

```go
package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"
)

const (
	numElements = 10000000
)

func timeGC() {
	t := time.Now()
	runtime.GC()
	fmt.Printf("gc took: %s\n", time.Since(t))
}

var pointers = map[string]int{}

func main() {
	for i := 0; i < 10000000; i++ {
		pointers[strconv.Itoa(i)] = i
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}
}

```

以下`GC`时间减少

```go
package main

import (
	"fmt"
	"runtime"
	"time"
)

func timeGC() {
	t := time.Now()
	runtime.GC()
	fmt.Printf("gc took: %s\n", time.Since(t))
}

type Entity struct {
	A int
	B float64
}

// 相比没有频繁创建对象，只是更改map.key
var entities = map[Entity]int{}

func main() {
	for i := 0; i < 10000000; i++ {
		entities[Entity{
			A: i,
			B: float64(i),
		}] = i
	}

	for {
		timeGC()
		time.Sleep(1 * time.Second)
	}
}

```

## 5. 使用`strings.Builder`拼接字符串

![image-20200507155642438](/home/lkt/桌面/notes/pprof.assets/image-20200507155642438.png)

```go
package no6_strings_builder

import (
	"strings"
	"testing"
)

var str string

var strs = []string{
	"here's",
	"a",
	"some",
	"long",
	"list",
	"of",
	"strings",
	"for",
	"you",
}

func BuildStrRaw(strs []string) string {
	var s string

	for _, v := range strs {
		s += v
	}

	return s
}

func BuildStrBuilder(strs []string) string {
	b := strings.Builder{}

	for _, v := range strs {
		b.WriteString(v)
	}

	return b.String()
}

func BuildStrPreAllocBuilder(strs []string) string {
	b := strings.Builder{}
	b.Grow(128)

	for _, v := range strs {
		b.WriteString(v)
	}

	return b.String()
}

func BenchmarkStringBuildRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str = BuildStrRaw(strs)
	}
}

func BenchmarkStringBuildBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str = BuildStrBuilder(strs)
	}
}

func BenchmarkStringPreAllocBuildBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		str = BuildStrPreAllocBuilder(strs)
	}
}

```

## 6. 分区：减少共享数据结构争夺

原理:`减少加锁力度`

![image-20200507155923806](/home/lkt/桌面/notes/pprof.assets/image-20200507155923806.png)

# &#x20;实践--参考go-zero部分

大部分来自文档\:Go 服务自动收集线上问题现场 <https://mp.weixin.qq.com/s/yYFM3YyBbOia3qah3eRVQA>

ps aux | grep service\_name 找到采集的服务进程id

源码位置: <https://github.com/zeromicro/go-zero/blob/master/core/proc/signals.go>

```go
func init() {
  go func() {
    ...
    signals := make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGTERM)

    for {
      v := <-signals
      switch v {
      ...
      case syscall.SIGUSR2:	// 这里收到USR2信号进行采集, 再次收到结束采集  kill -usr2 <process_id>
        if profiler == nil {
          profiler = StartProfile()
        } else {
          profiler.Stop()
          profiler = nil
        }
      ...
    }
  }()
}
```

调用StartProfile()

源码: <https://github.com/zeromicro/go-zero/blob/90828a0d4ae28fffd18f4d3e7c53246c802e7c1c/core/proc/profile.go#L168>

```go
func StartProfile() Stopper {
	if !atomic.CompareAndSwapUint32(&started, 0, 1) {
		logx.Error("profile: Start() already called")
		return noopStopper
	}

	var prof Profile
	prof.startCpuProfile()	//   每个都会创建采集的指标的文件
	prof.startMemProfile()
	prof.startMutexProfile()
	prof.startBlockProfile()
	prof.startTraceProfile()
	prof.startThreadCreateProfile()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		<-c

		logx.Info("profile: caught interrupt, stopping profiles")
		prof.Stop()

		signal.Reset()
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	return &prof
}

func (p *Profile) startCpuProfile() {
	fn := createDumpFile("cpu")
	f, err := os.Create(fn)
	if err != nil {
		logx.Errorf("profile: could not create cpu profile %q: %v", fn, err)
		return
	}

	logx.Infof("profile: cpu profiling enabled, %s", fn)
	pprof.StartCPUProfile(f)
	p.closers = append(p.closers, func() {
		pprof.StopCPUProfile()
		f.Close()
		logx.Infof("profile: cpu profiling disabled, %s", fn)
	})
}

...

func createDumpFile(kind string) string {
	command := path.Base(os.Args[0])
	pid := syscall.Getpid()
	return path.Join(os.TempDir(), fmt.Sprintf("%s-%d-%s-%s.pprof",
		command, pid, kind, time.Now().Format(timeFormat)))
}

// Stop 会执行关闭所有的采集
func (p *Profile) Stop() {
	if !atomic.CompareAndSwapUint32(&p.stopped, 0, 1) {
		// someone has already called close
		return
	}
	p.close()
	atomic.StoreUint32(&started, 0)
}

```

值得注意的是收集的信息都在 `/tmp` 文件夹下，以这个服务名命名的如下：

    - xxxx-mq-cpu-xxx.pprof
    - xxxx-mq-memory-xxx.pprof
    - xxxx-mq-mutex-xxx.pprof
    - xxxx-mq-block-xxx.pprof
    - .......

