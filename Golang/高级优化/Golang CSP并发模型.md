[CSP并发模型https://www.jianshu.com/p/36e246c6153d](https://www.jianshu.com/p/36e246c6153d)

# 总结
* `Golang实现了 CSP` 并发模型做为并发基础，底层使用`goroutine`做为并发实体，`goroutine`非常轻量级可以创建几十万个实体。
* 实体间通过 `channel` 继续匿名消息传递使之解耦，在语言层面实现了自动调度，这样屏蔽了很多内部细节，对外提供简单的语法关键字，大大简化了并发编程的思维转换和管理线程的复杂性。
* `goroutine`并发是语言层面实现的，实现自己的调度系统，在用户态上运行，很少涉及内核态切换，减少开销

# `CSP`
今天介绍一下 go语言的并发机制以及它所使用的CSP并发模型

## `CSP`并发模型
`CSP`模型是上个世纪七十年代提出的，用于描述两个独立的并发实体通过共享的通讯 `channel`(管道)进行通信的并发模型。 `CSP`中`channel`是第一类对象，它不关注发送消息的实体，而关注与发送消息时使用的`channel`。

## `Golang CSP`
`Golang` 就是借用`CSP`模型的一些概念为之实现并发进行理论支持，其实从实际上出发，go语言并没有，完全实现了`CSP`模型的所有理论，仅仅是借用了 `process`和`channel`这两个概念。`process`是在`go`语言上的表现就是 `goroutine` 是实际并发执行的实体，每个实体之间是通过`channel`通讯来实现数据共享。

## `Channel`
Golang中使用 CSP中 channel 这个概念。channel 是被单独创建并且可以在进程之间传递，它的通信模式类似于 boss-worker 模式的，一个实体通过将消息发送到channel 中，然后又监听这个 channel 的实体处理，两个实体之间是匿名的，这个就实现实体中间的解耦，其中 channel 是同步的一个消息被发送到 channel 中，最终是一定要被另外的实体消费掉的，在实现原理上其实是一个阻塞的消息队列。

## `Goroutine`
Goroutine 是实际并发执行的实体，它底层是使用协程(coroutine)实现并发，coroutine是一种运行在用户态的用户线程，类似于 greenthread，go底层选择使用coroutine的出发点是因为，它具有以下特点：

用户空间 避免了内核态和用户态的切换导致的成本
可以由语言和框架层进行调度
更小的栈空间允许创建大量的实例
可以看到第二条 用户空间线程的调度不是由操作系统来完成的，像在java 1.3中使用的greenthread的是由JVM统一调度的(后java已经改为内核线程)，还有在ruby中的fiber(半协程) 是需要在重新中自己进行调度的，而goroutine是在golang层面提供了调度器，并且对网络IO库进行了封装，屏蔽了复杂的细节，对外提供统一的语法关键字支持，简化了并发程序编写的成本。

Goroutine 调度器
上节已经说了，golang使用goroutine做为最小的执行单位，但是这个执行单位还是在用户空间，实际上最后被处理器执行的还是内核中的线程，用户线程和内核线程的调度方法有：

M:N 用户线程和内核线程是多对多的对应关系
![image](https://upload-images.jianshu.io/upload_images/1767848-9c4b06362907280d.png?imageMogr2/auto-orient/strip|imageView2/2/w/350/format/webp)


golang 通过为goroutine提供语言层面的调度器，来实现了高效率的M:N线程对应关系

调度示意图中
![image](https://upload-images.jianshu.io/upload_images/1767848-fc23b15dc52e407f.jpg?imageMogr2/auto-orient/strip|imageView2/2/w/307/format/webp)


M：是内核线程
P : 是调度协调，用于协调M和G的执行，内核线程只有拿到了 P才能对goroutine继续调度执行，一般都是通过限定P的个数来控制golang的并发度
G : 是待执行的goroutine，包含这个goroutine的栈空间
Gn : 灰色背景的Gn 是已经挂起的goroutine，它们被添加到了执行队列中，然后需要等待网络IO的goroutine，当P通过 epoll查询到特定的fd的时候，会重新调度起对应的，正在挂起的goroutine。
Golang为了调度的公平性，在调度器加入了steal working 算法 ，在一个P自己的执行队列，处理完之后，它会先到全局的执行队列中偷G进行处理，如果没有的话，再会到其他P的执行队列中抢G来进行处理。

# `GPM`
首先GPM是golang runtime里面的东西，是语言层面的实现。也就是说go实现了自己的调度系统。 理解了这一点 再往下看
M（machine）是runtime对操作系统内核线程的虚拟， M与内核线程一般是一一映射的关系， 一个groutine最终是要放到M上执行的；
P管理着一组Goroutine队列，P里面一般会存当前goroutine运行的上下文环境（函数指针，堆栈地址及地址边界），P会对自己管理的goroutine队列做一些调度（比如把占用CPU时间较长的goroutine暂停 运行后续的goroutine等等。。）当自己的队列消耗完了 会去全局队列里取， 如果全局队列里也消费完了 会去其他P对立里取。
G 很好理解，就是个goroutine的，里面除了存放本goroutine信息外 还有与所在P的绑定等信息。

GPM协同工作 组成了runtime的调度器。

P与M一般也是一一对应的。他们关系是： P管理着一组G挂载在M上运行。当一个G长久阻塞在一个M上时，runtime会新建一个M，阻塞G所在的P会把其他的G 挂载在新建的M上。当旧的G阻塞完成或者认为其已经死掉时 回收旧的M。

P的个数是通过runtime.GOMAXPROCS设定的，现在一般不用自己手动设，默认物理线程数（比如我的6核12线程， 值会是12）。 在并发量大的时候会增加一些P和M，但不会太多，切换太频繁的话得不偿失。内核线程的数量一般大于12这个值， 不要错误的认为M与物理线程对应，M是与内核线程对应的。 如果服务器没有其他服务的话，M才近似的与物理线程一一对应。

说了这么多。初步了解了go的调度，我想大致也明白了， 单从线程调度讲，go比起其他语言的优势在哪里了？
go的线程模型是M：N的。 其一大特点是goroutine的调度是在用户态下完成的， 不涉及内核态与用户态之间的频繁切换，包括内存的分配与释放，都是在用户态维护着一块大的内存池， 不直接调用系统的malloc函数（除非内存池需要改变）。 另一方面充分利用了多核的硬件资源，近似的把若干goroutine均分在物理线程上， 再加上本身goroutine的超轻量，以上种种保证了go调度方面的性能。