- [`sync.Map`原理](#syncmap原理)
	- [数据结构](#数据结构)
	- [ 原理流程图](#-原理流程图)
		- [ Load](#-load)
		- [ Store](#-store)
		- [Delete](#delete)
	- [总结](#总结)
- [`skipmap`](#skipmap)
	- [适合场景](#适合场景)
	- [基本使用](#基本使用)
	- [关键步骤 `Store`](#关键步骤-store)
- [`concurrent-map` 并发Map](#concurrent-map-并发map)
	- [介绍](#介绍)
	- [简单使用](#简单使用)
	- [关键原理](#关键原理)


# `sync.Map`原理

go 1.9 官方提供了sync.Map 来优化线程安全的并发读写的map。该实现也是基于内置map关键字来实现的。

这个实现类似于一个线程安全的 `map[interface{}]interface{}` . 这个map的优化主要适用了以下场景：

1.  给定key的键值对只写了一次，但是**读了很多次**，比如在只增长的缓存中；&#x20;
2.  当多个`goroutine`读取、写入和覆盖的key值不相交时。

参考讲解: <https://blog.csdn.net/m0_67402013/article/details/124244876>

## 数据结构

```go
// 封装的线程安全的map
type Map struct {
	// lock
	mu Mutex

	// 实际是readOnly这个结构
	// 一个只读的数据结构，因为只读，所以不会有读写冲突。
	// readOnly包含了map的一部分数据，用于并发安全的访问。(冗余，内存换性能)
	// 访问这一部分不需要锁。
	read atomic.Value // readOnly -- 这里可以看出比较适合多读场景

	// dirty数据包含当前的map包含的entries,它包含最新的entries(包括read中未删除的数据,虽有冗余，但是提升dirty字段为read的时候非常快，不用一个一个的复制，
	//而是直接将这个数据结构作为read字段的一部分),有些数据还可能没有移动到read字段中。
	// 对于dirty的操作需要加锁，因为对它的操作可能会有读写竞争。
	// 当dirty为空的时候， 比如初始化或者刚提升完，下一次的写操作会复制read字段中未删除的数据到这个数据中。
	dirty map[interface{}]*entry

	// 当从Map中读取entry的时候，如果read中不包含这个entry,会尝试从dirty中读取，这个时候会将misses加一，
	// 当misses累积到 dirty的长度的时候， 就会将dirty提升为read,避免从dirty中miss太多次。因为操作dirty需要加锁。
	misses int
}

// readOnly is an immutable struct stored atomically in the Map.read field.
type readOnly struct {
	m       map[interface{}]*entry
	// 如果Map.dirty有些数据不在m中，这个值为true
	amended bool 
}
// entry 就是具体的内容
// An entry is a slot in the map corresponding to a particular key.
type entry struct {
	// *interface{}
	p unsafe.Pointer 
}

```

个人理解:&#x20;

`Map`里面还是会有锁, 因为涉及到并发写入, `sync.Map`会把存储内容读写分离, 单独把读分离,可以不需要加锁,做到高效, 先读取 `readOnly`部分,再去使用锁读取`dirty`部分

`sync.Map`里对存储进行的特殊处理, 写入时候会先写入到`dirty`中, 会在**适的时候把数据结构直接作为read字段的一部分(减少复制)**

删除操作, 找到元素, 先找`readOnly` 没有,就加上锁找`dirty`, 然后把数据标记起来,` p == expunged 惰性删除`

## &#x20;原理流程图

### &#x20;Load

![https://img-blog.csdnimg.cn/20200112110724714.png?x-oss-process=image/watermark,type\_ZmFuZ3poZW5naGVpdGk,shadow\_10,text\_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size\_16,color\_FFFFFF,t\_70](https://img-blog.csdnimg.cn/20200112110724714.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size_16,color_FFFFFF,t_70)

### &#x20;Store

![https://img-blog.csdnimg.cn/20200112110342685.png?x-oss-process=image/watermark,type\_ZmFuZ3poZW5naGVpdGk,shadow\_10,text\_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size\_16,color\_FFFFFF,t\_70](https://img-blog.csdnimg.cn/20200112110342685.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size_16,color_FFFFFF,t_70)

### Delete

![https://img-blog.csdnimg.cn/20200112111125713.png?x-oss-process=image/watermark,type\_ZmFuZ3poZW5naGVpdGk,shadow\_10,text\_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size\_16,color\_FFFFFF,t\_70](https://img-blog.csdnimg.cn/20200112111125713.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size_16,color_FFFFFF,t_70)

## 总结

无锁读与读写分离；

![https://img-blog.csdnimg.cn/20200112110232202.png?x-oss-process=image/watermark,type\_ZmFuZ3poZW5naGVpdGk,shadow\_10,text\_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size\_16,color\_FFFFFF,t\_70](https://img-blog.csdnimg.cn/20200112110232202.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size_16,color_FFFFFF,t_70)

写加锁与延迟提升；

![https://img-blog.csdnimg.cn/20200112110214618.png?x-oss-process=image/watermark,type\_ZmFuZ3poZW5naGVpdGk,shadow\_10,text\_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size\_16,color\_FFFFFF,t\_70](https://img-blog.csdnimg.cn/20200112110214618.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9sb3V5dXRpbmcuYmxvZy5jc2RuLm5ldA==,size_16,color_FFFFFF,t_70)

***指针与惰性删除，map保存的value都是指针。惰性删除，实际删除是在 Store时候去check 然后删除***

# `skipmap`

跳表原理参考: \[Skip List--跳表（全网最详细的跳表文章没有之一）]\(<https://www.jianshu.com/p/9d8296562806?u_atoken=5ce682b7-b395-43cb-960f-dee2f7f8e3ab&u_asession=01iC6xtKR2RdODDKV6qe-pYbADqLaXkQwy2DSm-NCSf9iZclEWLyqvEsKa7XYgB9yRX0KNBwm7Lovlpxjd_P_q4JsKWYrT3W_NKPr8w6oU7K_5HD0hOnTi8YcueOzRXTuO3KmjkU3JT7ddtoHBlecZWGBkFo3NEHBv0PZUm6pbxQU&u_asig=05s4r5di7KSYa8I9JfvIpj9aa5Qt56lzfgDNAP8nd_DmUtEjQk85YYhcOXcKS776QfrmbET43Eatnuq5QOLDtJlvPTD04Bq5h1PC4i3otoaVOkaQq2A0Zx-fAm9JH4H4WY8Fzx4Xevfk3DuEOzKIz5yIAMuCOVpcwbUeG5zlCZTXr9JS7q8ZD7Xtz2Ly-b0kmuyAKRFSVJkkdwVUnyHAIJzcMgP2TqTJQPpRJQbBaHTz40BQBm_zGMT6FTf_aiLqTZzKnPGeiYgOeAvNODIGQOu-3h9VXwMyh6PgyDIVSG1W-xoyRAjwp8eky6Ph6HdID7c4qWnPeNnkyWEERpY2mCvY2fjjkWLSUV2cBSy_65vP65B247cIonsGh3LdhR6wUNmWspDxyAEEo4kbsryBKb9Q&u_aref=4gtnaOdWU1QP5MO3y9O61rtbSpw%3D>)

第三方包: <https://github.com/XiaoLiu1020/skipmap>

注意: 要求使用Go1.18版本以上, 因为使用了泛型

Go Doc: <https://pkg.go.dev/github.com/zhangyunhao116/skipmap>

## 适合场景

*   需要key排序
*   并发进行操作, 比如调用 `Range` and `Store`同时, 这些情况下会大大提高性能

## 基本使用

```go
package main

import (
	"fmt"

	"github.com/zhangyunhao116/skipmap"
)

//  需要使用泛型

type Node struct {
	name string
	id   int
}

func main() {
	// Typed key and generic value.
	//  使用泛型,存储Node类型
	m0 := skipmap.NewString[Node]()

	node := Node{name: "liukaitao", id: 1}

	// 存储
	m0.Store("yanshiyun", node)

	v, ok := m0.Load("10")
	if ok {
		fmt.Println("skipmap load key 10 with value ", v)
	}

	v, ok = m0.Load("yanshiyun")
	if ok {
		fmt.Println("skipmap load key yanshiyun with value, ", v)
	}

	m0.Range(func(key string, value Node) bool {
		fmt.Println("m0 Found: ", key, value)
		return true
	})
	m0.Delete("yanshiyun")
	fmt.Printf("skipmap contains %d items\r\n", m0.Len())
}

```

## 关键步骤 `Store`

```go
// Store sets the value for a key.
func (s *StringMap[valueT]) Store(key string, value valueT) {
	level := s.randomlevel()
	var preds, succs [maxLevel]*stringnode[valueT]
	for {
		nodeFound := s.findNode(key, &preds, &succs)
		if nodeFound != nil { // indicating the key is already in the skip-list
			if !nodeFound.flags.Get(marked) {
				// We don't need to care about whether or not the node is fully linked,
				// just replace the value.
				nodeFound.storeVal(value)
				return
			}
			// If the node is marked, represents some other goroutines is in the process of deleting this node,
			// we need to add this node in next loop.
			continue
		}

		// Add this node into skip list.
		var (
			highestLocked        = -1 // the highest level being locked by this process
			valid                = true
			pred, succ, prevPred *stringnode[valueT]
		)
		for layer := 0; valid && layer < level; layer++ {
			pred = preds[layer]   // target node's previous node
			succ = succs[layer]   // target node's next node
			if pred != prevPred { // the node in this layer could be locked by previous loop
				pred.mu.Lock()
				highestLocked = layer
				prevPred = pred
			}
			// valid check if there is another node has inserted into the skip list in this layer during this process.
			// It is valid if:
			// 1. The previous node and next node both are not marked.
			// 2. The previous node's next node is succ in this layer.
			valid = !pred.flags.Get(marked) && (succ == nil || !succ.flags.Get(marked)) && pred.loadNext(layer) == succ
		}
		if !valid {
			unlockstring(preds, highestLocked)
			continue
		}

		nn := newStringNode(key, value, level)
		for layer := 0; layer < level; layer++ {
			nn.storeNext(layer, succs[layer])
			preds[layer].atomicStoreNext(layer, nn)
		}
		nn.flags.SetTrue(fullyLinked)
		unlockstring(preds, highestLocked)
		atomic.AddInt64(&s.length, 1)
		return
	}
}

```

# `concurrent-map` 并发Map
第三方包: https://github.com/orcaman/concurrent-map

中文介绍: https://github.com/orcaman/concurrent-map/blob/master/README-zh.md

## 介绍
标准库中的`sync.Map`是专为`append-only`场景设计的。因此，如果`您想将Map用于一个类似内存数据库`，那么使用我们的版本可能会受益

## 简单使用
`go get "github.com/orcaman/concurrent-map/v2"`

```golang
import (
	"github.com/orcaman/concurrent-map/v2"
)

// 创建一个新的 map.
	m := cmap.New[string]()

	// 设置变量m一个键为“foo”值为“bar”键值对
	m.Set("foo", "bar")

	// 从m中获取指定键值.
	bar, ok := m.Get("foo")

	// 删除键为“foo”的项
	m.Remove("foo")
```

## 关键原理
对`map`进行分片处理, 减少锁粒度
```golang
// A "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (SHARD_COUNT) map shards.
type ConcurrentMap[K comparable, V any] struct {
	shards   []*ConcurrentMapShared[K, V]
	sharding func(key K) uint32
}

```