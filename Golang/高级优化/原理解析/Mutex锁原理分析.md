- [参考文档](#参考文档)
- [图解`Go`互斥锁`mutex`核心实现原理](#图解go互斥锁mutex核心实现原理)
- [1. 锁基本概念](#1-锁基本概念)
	- [1.1 `CAS`与轮询](#11-cas与轮询)
		- [普通`cas`--`Compare And Swap`](#普通cas--compare-and-swap)
		- [轮询锁](#轮询锁)
	- [1.2 锁公平性](#12-锁公平性)
	- [1.3 饥饿与排队](#13-饥饿与排队)
		- [锁饥饿](#锁饥饿)
		- [排队机制](#排队机制)
	- [1.4 位计数](#14-位计数)
- [2. `mutex`实现](#2-mutex实现)
	- [2.1 成员变量与模式](#21-成员变量与模式)
		- [2.1.1 成员变量](#211-成员变量)
		- [2.1.2 锁模式](#212-锁模式)
		- [2.1.3 `normal`模式](#213-normal模式)
		- [2.1.4 `starvation`模式](#214-starvation模式)
	- [2.2 锁计数](#22-锁计数)
		- [2.2.1 锁状态](#221-锁状态)
		- [2.2.2 等待计数](#222-等待计数)
	- [自旋：](#自旋)
		- [什么是自旋:](#什么是自旋)
		- [自旋条件](#自旋条件)
		- [自旋问题:](#自旋问题)
	- [2.3 唤醒机制](#23-唤醒机制)
		- [2.3.1 唤醒标志`Woken`](#231-唤醒标志woken)
		- [2.3.2 唤醒流程](#232-唤醒流程)
	- [2.3 加锁流程](#23-加锁流程)
		- [2.3.1 快速模式](#231-快速模式)
		- [2.3.2 自旋与唤醒](#232-自旋与唤醒)
		- [2.3.3 更改锁状态](#233-更改锁状态)
		- [2.4.5 加锁排队与状态转换](#245-加锁排队与状态转换)
	- [2.5 释放锁逻辑](#25-释放锁逻辑)
		- [2.5.1 释放锁代码](#251-释放锁代码)


# 参考文档
https://rainbowmango.gitbook.io/go/chapter02/2.4-mutex#42-starvation-mo-shi

https://www.cnblogs.com/buyicoding/p/12082162.html


# 图解`Go`互斥锁`mutex`核心实现原理

# 1. 锁基本概念

## 1.1 `CAS`与轮询

### 普通`cas`--`Compare And Swap`

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619037-1751266629.png)

原理就是调用系统cpu的原子性命令：拿查询到的`old`值跟`真正old`值比较，是就替换为`new值`，返回`swapped = true`, 不是就返回`false`，替换不成功

利用`处理器的CAS指令`来实现对给定变量的值交换来进行锁的获取

### 轮询锁

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619198-1740057135.png)

在多线程并发的情况下`很有可能会有线程CAS失败`，通常就会`配合for循环采用轮询的方式去尝试重新获取锁`

## 1.2 锁公平性

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619335-1215377483.png)

- 先进行锁获取的线程是否比后续的线程更先获得锁，如果是则就是公平锁
- 多个线程按照获取锁的顺序依次获得锁，否则就是非公平性

## 1.3 饥饿与排队

### 锁饥饿

锁饥饿是指因为大量线程都同时进行获取锁，`某些线程可能在锁的CAS过程中一直失败，从而长时间获取不到锁`

### 排队机制

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619478-1432181166.png)

上面提到了CAS和轮询锁进行锁获取的方式，可以发现如果已经有线程获取了锁，但是在当前线程在多次轮询获取锁失败的时候，就`没有必要再继续进行反复尝试浪费系统资源`，通常就会`采用一种排队机制`，来进行排队等待

## 1.4 位计数

在大多数编程语言中针对实现基于CAS的锁的时候，通常都会采用一个32位的整数来进行锁状态的存储

------

# 2. `mutex`实现

## 2.1 成员变量与模式

### 2.1.1 成员变量

在go的mutex中核心成员变量只有两个`state`和`sema`,其通过`state来进行锁的计数`，而`通过sema来实现排队`

```go
type Mutex struct {
	state int32			// 锁的计数
	sema  uint32	// 实现排队
}
```

### 2.1.2 锁模式

|          | 描述                                                                                                                                                              | 公平性 |
| -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| 正常模式 | 默认情况下，Mutex的模式为normal。<br /><br />该模式下，协程如果加锁不成功不会立即转入阻塞排队，而`是判断是否满足自旋的条件，如果满足则会启动自旋过程，尝试抢锁。` | 否     |
| 饥饿模式 | 处于饥饿模式下，`不会启动自旋过程`，也即一旦有协程释放了锁，那么一定会唤醒协程，被唤醒的协程将会成功获取锁，同时也会把等待计数减1。                               | 是     |

- 在正常模式下，其实锁的性能是最高的, 如果`多个goroutine进行锁获取后立马进行释放`则可以避免多个线程的排队消耗
- 切换到饥饿模式后，在进行锁获取的时候，如果满足一定的条件也会切换回正常模式，从而保证锁的高性能



### 2.1.3 `normal`模式

默认情况下，Mutex的模式为normal。

该模式下，协程如果加锁不成功不会立即转入阻塞排队，而是判断是否满足自旋的条件，如果满足则会启动自旋过程，尝试抢锁。

### 2.1.4 `starvation`模式

切换原理：`自旋过程中能抢到锁，一定意味着同一时刻有协程释放了锁，`我们知道释放锁时如果发现有阻塞等待的协程，还会释放一个信号量来唤醒一个等待协程，`被唤醒的协程得到CPU后开始运行，此时发现锁已被抢占了，自己只好再次阻塞，不过阻塞前会判断自上次阻塞到本次阻塞经过了多长时间`，`如果超过1ms的话，会将Mutex标记为"饥饿"模式，然后再阻塞。`

处于饥饿模式下，不会启动自旋过程，也即一旦有协程释放了锁，那么一定会唤醒协程，被唤醒的协程将会成功获取锁，同时也会把等待计数减1。



## 2.2 锁计数

```go
const (
    // 三个标记位
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken						// 锁唤醒
	mutexStarving					// 饥饿
    
	mutexWaiterShift = iota		//

	starvationThresholdNs = 1e6
)
```



### 2.2.1 锁状态

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619613-1159743212.png)

在mutex中锁有三个标志位，其中其二进制位分别位001(mutexLocked)、010(mutexWoken)、100(mutexStarving), 注意这三者并不是互斥的关系，比如一个锁的状态可能是锁定的饥饿模式并且已经被唤醒

```go
const (
    // 锁状态常量
    // 三个标记位
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken						// 锁唤醒
	mutexStarving					// 饥饿
)
```

> `iota`相当于`其在const命名块的行索引`，　第一行`iota`为０
>
> `<<`为左移，　`1<<3` 相当于`二进制１左移动三位`，`=8`

### 2.2.2 等待计数

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619740-1472434277.png)

mutex中通过低3位存储了当前mutex的三种状态，剩下的29位全部用来存储尝试正在等待获取锁的goroutine的数量

```go
	mutexWaiterShift = iota // 3
```

------

## 自旋：

加锁时，如果当前Locked位为1，说明该锁当前由其他协程持有，`尝试加锁的协程并不是马上转入阻塞，而是会持续的探测Locked位是否变为0，这个过程即为自旋过程。`

自旋时间很短，但如果在自旋过程中发现锁已被释放，那么协程可以立即获取锁。此时即便有协程被唤醒也无法获取锁，只能再次阻塞。

自旋的好处是，`当加锁失败时不必立即转入阻塞，有一定机会获取到锁，这样可以避免协程的切换。`

### 什么是自旋:

自旋对应于CPU的"PAUSE"指令，CPU对该指令什么都不做，相当于CPU空转，对程序而言相当于sleep了一小段时间，时间非常短，当前实现是30个时钟周期。

自旋过程中会持续探测Locked是否变为0，连续两次探测间隔就是`执行这些PAUSE指令，它不同于sleep，不需要将协程转为睡眠状态。`



### 自旋条件

加锁时程序会自动判断是否可以自旋，无限制的自旋将会给CPU带来巨大压力，所以判断是否可以自旋就很重要了。

自旋必须满足以下所有条件：

- `自旋次数要足够小，通常为4`，即自旋最多4次
- `CPU核数要大于1，否则自旋没有意义`，因为此时不可能有其他协程释放锁
- 协程调度机制中的Process数量要大于1，比如使用GOMAXPROCS()将处理器设置为1就不能启用自旋
- 协程调度机制中的`可运行队列必须为空`，否则会延迟协程调度

可见，自旋的条件是很苛刻的，总而言之就是`不忙的时候才会启用自旋。`

### 自旋问题:

如果自旋过程中获得锁，那么之前被阻塞的协程将无法获得锁，如果加锁的协程特别多，每次都通过自旋获得锁，那么之前被阻塞的进程将很难获得锁，从而进入饥饿状态。



`Starving状态`。这个状态下不会自旋，一旦有协程释放锁，那么一定会唤醒一个协程并成功加锁。

------

## 2.3 唤醒机制

暂时不看`Woken`的解锁并唤醒协程

![img](https://gblobscdn.gitbook.com/assets%2F-LQm0KQP9eyG1B9ntPkR%2F-LRQkwhW7HWUzIL11Bvy%2F-LRQkxH3uwy8jwK3gIQj%2Fmutex-05-unlock_with_waiter.png?alt=media)

协程A解锁过程分为两个步骤，一是把Locked位置0，二是查看到Waiter>0，所以释放一个信号量，唤醒一个阻塞的协程，被唤醒的协程B把Locked位置1，于是协程B获得锁。



### 2.3.1 唤醒标志`Woken`

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093619874-692221445.png)

唤醒标志其实就是上面说的第二位，　唤醒标志主要用于`标识当前尝试获取goroutine是否有正在处于唤醒状态的`

记得上面公平模式下，当前正在cpu上运行的goroutine可能会先获取到锁



`Woken状态`用于加锁和解锁过程的通信，举个例子，`同一时刻，两个协程一个在加锁，一个在解锁，在加锁的协程可能在自旋过程中，此时把Woken标记为1，用于通知解锁协程不必释放信号量了`，好比在说：你只管解锁好了，不必释放信号量，我马上就拿到锁了。



### 2.3.2 唤醒流程

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093620014-48496094.png)

1. 当`释放锁`的时候，如果当前有goroutine正在唤醒状态`Woken=1`，则只需要修改锁状态为释放锁，则处于woken状态的goroutine(可能在自旋中)就可以直接获取锁
2. `Woken=0`否则则需要唤醒一个goroutine(释放信号量), `并且等待这个goroutine修改state状态为mutexWoken`，才退出







## 2.3 加锁流程

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093620929-1160638259.png)





### 2.3.1 快速模式

如果当前没有goroutine加锁，则直接进行CAS成功，则直接获取锁成功

```go
		// Fast path: grab unlocked mutex.
	// cas 设置加锁，　成功证明直接cas成功，不成功证明有锁存在
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// 不成功走这
	m.lockSlow()　
```

### 2.3.2 自旋与唤醒

接上`m.lockSlow()`

```go
func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false					
	iter := 0									// 自旋次数
	old := m.state						// 锁状态 old
	for {
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
        // 判断　old锁　是否锁了并且是否可以自旋转
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			// Active spinning makes sense.
			// Try to set mutexWoken flag to inform Unlock
			// to not wake other blocked goroutines.
			
			// !awoke 如果当前线程不处于唤醒状态
			// old&mutexWoken == 0 如果当前不处于　正在唤醒的节点的状态
			// old>> mutexWaiterShift != 0 : 右移３位，　如果不为０，证明有等待中的`goroutine`
			// 设置当前状态为唤醒状态成功
			//　当满足以上四种条件
			if !awoke && old&mutexWoken == 0 && old>>mutexWaiterShift != 0 &&
				
				atomic.CompareAndSwapInt32(&m.state, old, old|mutexWoken) {
				// 更改当前线程为唤醒状态
				awoke = true
			}
			// 尝试自旋
			runtime_doSpin()
			// 自旋计数
			iter++
			// 从新获取状态
			old = m.state
			continue
		}
		new := old
        .....
```

> 101＆100　= 100　选出重合的位， 都不重合返回0
>
> 001 | 100 = 101 　合并位

### 2.3.3 更改锁状态

流程走到这里会有两种可能：

1. 锁状态当前已经不是锁定状态
2. 自旋超过指定的次数，不再允许自旋了	

```go
		new := old
		if old&mutexStarving == 0 {
			// 如果当前不是饥饿模式，则这里其实就可以尝试进行锁的获取了|=其实就是将锁的那个bit位设为1表示锁定状态
			new |= mutexLocked
		}
		if old&(mutexLocked|mutexStarving) != 0 {
			// 如果当前被锁定或者处于饥饿模式，则增等待一个等待计数
			new += 1 << mutexWaiterShift
		}
		if starving && old&mutexLocked != 0 {
			// 如果当前已经处于饥饿状态，并且当前锁还是被占用，则尝试进行饥饿模式的切换
			new |= mutexStarving
		}
		if awoke {
			if new&mutexWoken == 0 {
				throw("sync: inconsistent mutex state")
			}
			// awoke为true则表明当前线程在上面自旋的时候，修改mutexWoken状态成功
			// 清除唤醒标志位
            // 为什么要清除标志位呢？
            // 实际上是因为后续流程很有可能当前线程会被挂起,就需要等待其他释放锁的goroutine来唤醒
            // 但如果unlock的时候发现mutexWoken的位置不是0，则就不会去唤醒，则该线程就无法再醒来加锁
			new &^= mutexWoken
		}
```

### 2.4.5 加锁排队与状态转换

再加锁的时候实际上只会有一个goroutine加锁CAS成功，而其他线程则需要重新获取状态，进行上面的自旋与唤醒状态的重新计算，从而再次CAS

```go
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				// 如果原来的状态等于0则表明当前已经释放了锁并且也不处于饥饿模式下
                // 实际的二进制位可能是这样的 1111000, 后面三位全是0，只有记录等待goroutine的计数器可能会不为0
                // 那就表明其实
				break // locked the mutex with CAS
			}
			// 排队逻辑，如果发现waitStatrTime不为0，则表明当前线程之前已经再排队来，后面可能因为
            // unlock被唤醒，但是本次依旧没获取到锁，所以就将它移动到等待队列的头部
			queueLifo := waitStartTime != 0
			if waitStartTime == 0 {
				waitStartTime = runtime_nanotime()
			}
            // 这里就会进行排队等待其他节点进行唤醒
			runtime_SemacquireMutex(&m.sema, queueLifo)
			// 如果等待超过指定时间，则切换为饥饿模式 starving=true
            // 如果一个线程之前不是饥饿状态，并且也没超过starvationThresholdNs，则starving为false
            // 就会触发下面的状态切换
			starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
			// 重新获取状态
            old = m.state
			if old&mutexStarving != 0 { 
                // 如果发现当前已经是饥饿模式，注意饥饿模式唤醒的是第一个goroutine
                // 当前所有的goroutine都在排队等待
			// 一致性检查，
				if old&(mutexLocked|mutexWoken) != 0 || old>>mutexWaiterShift == 0 {
					throw("sync: inconsistent mutex state")
				}
				// 获取当前的模式
				delta := int32(mutexLocked - 1<<mutexWaiterShift)
				if !starving || old>>mutexWaiterShift == 1 {
					// 如果当前goroutine不是饥饿状态，就从饥饿模式切换会正常模式
                    // 就从mutexStarving状态切换出去
					delta -= mutexStarving
				}
                // 最后进行cas操作
				atomic.AddInt32(&m.state, delta)
				break
			}
            // 重置计数
			awoke = true
			iter = 0
		} else {
			old = m.state
		}
```



## 2.5 释放锁逻辑

![image.png](https://img2018.cnblogs.com/blog/1506724/201912/1506724-20191223093621483-1011216849.png)

### 2.5.1 释放锁代码

```go
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}
	// 直接进行cas操作
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		// 如果释放锁并且不是饥饿模式
		old := new
		for {
	
			if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
				// 如果已经有等待者并且已经被唤醒，就直接返回
				return
			}
			// 减去一个等待计数，然后将当前模式切换成mutexWoken
			new = (old - 1<<mutexWaiterShift) | mutexWoken
			if atomic.CompareAndSwapInt32(&m.state, old, new) {
				// 唤醒一个goroutine
				runtime_Semrelease(&m.sema, false)
				return
			}
			old = m.state
		}
	} else {
		// 唤醒等待的线程
		runtime_Semrelease(&m.sema, true)
	}
}

```

