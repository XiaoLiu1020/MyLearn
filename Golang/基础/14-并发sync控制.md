**继续上一节**
## `sync.WaitGroup`
可以使用 `sync.WaitGroup` 实现并发任务的同步，有如下几种方法：

方法名 | 功能
---|---
`(wg *WaitGroup) Add(delta int)` | 计数器 + `delta`
`(wg *WaitGroup) Done()` | 计数器-1
`(wg *WaitGroup) Wait()` | 阻塞直到计数器为0

`sync.WaitGroup` 内部维护着计数器，启动任务将计数值添加N，完成一个就调用一个`Done()`，通过调用`Wait()`等待并发任务执行完

```
var wg sync.WaitGroup

func hello() {
    defer wg.Done()
    fmt.Println("Hello Goroutine!")
}

func main() {
    wg.Add(1)
    go hello()      //启动另外一个goroutine去执行hello函数
    fmt.Println("main goroutine done!")
    wg.Wait()
}
```

## `sync.Once`
很多场景下，我们需要确保某些操作在高并发场景下只执行一次，例如只加载一次配置文件，只关闭一次通道等

`sync.Once`只有一个`Do`方法

```
func (o *Once) Do(f func()) {
    
}
```
备注：**如果要执行的函数**`f`**需要传递参数就需要搭配闭包使用**

#### 加载配置文件示例
```
var icons map[string]image.Image

func loadIcons() {
	icons = map[string]image.Image{
		"left":  loadIcon("left.png"),
		"up":    loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down":  loadIcon("down.png"),
	}
}

// Icon 被多个goroutine调用时不是并发安全的
func Icon(name string) image.Image {
	if icons == nil {
		loadIcons()
	}
	return icons[name]
}
```
1. 多个`goroutine` 并发调用`Icon`函数并不是并发安全的
2. 现代的编译器和`cpu`可能会保证每个`goroutine`都满足串行一致的基础上自由重排访问内存顺序
3. `loadIcons`函数可能被重排以下结果：
```
func loadIcons() {
    icons = make(map[string]image.Image)    //产生并发问题
	icons["left"] = loadIcon("left.png")
	icons["up"] = loadIcon("up.png")
	icons["right"] = loadIcon("right.png")
	icons["down"] = loadIcon("down.png")
}
```
在这种情况下就会出现 即使判断了 `icons`也不是`nil`也不意味着变量初始化完成了。

`sync.Once` 其实内部包含一个互斥锁和一个布尔值，互斥锁保证布尔值和数据安全，而布尔值用来记录初始化是否完成了。

```
var icons map[string]image.Image

var loadIconsOnce sync.Once

func loadIcons() {
	icons = map[string]image.Image{
		"left":  loadIcon("left.png"),
		"up":    loadIcon("up.png"),
		"right": loadIcon("right.png"),
		"down":  loadIcon("down.png"),
	}
}

// Icon 是并发安全的
func Icon(name string) image.Image {
	loadIconsOnce.Do(loadIcons) //只会执行初始化一次
	return icons[name]
}
```

## `sync.Map`
`Go`语言中内置的 `map`并不是并发安全的

```
var m = make(map[string]int)

func get(key string) int {
	return m[key]
}

func set(key string, value int) {
	m[key] = value
}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			set(key, n)
			fmt.Printf("k=:%v,v:=%v\n", key, get(key))
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```
当并发多了之后，会报`fatal error: concurrent map writes`错误。

`sync`包中提供了并发安全版 `map`-`sync.Map`，开箱即用，一样使用`make`函数初始化，同时 `sync.Map`内置了诸如 `Store`,`Load`,`LOadOrStore`,`Delete`,`Range`等方法

```
var m = sync.Map{}

func main() {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			key := strconv.Itoa(n)
			m.Store(key, n)         //Store 储存
			value, _ := m.Load(key)
			fmt.Printf("k=:%v,v:=%v\n", key, value)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
```

# 原子操作
* 代码中的加锁操作 涉及**内核态的上下文切换**会比较耗时和代价比较高。
* **针对基本数据类型**，还可以使用**原子操作**来保证操作安全
* 原子操作是`Go`语言提供的方法，**在用户态就可以完成**，性能比加锁操作更好

## `atomic`包
传入参数为 `&x` 内存地址

方法 | 解释
---|---
func LoadInt32(addr *int32) (val int32) |读取操作
func LoadInt64(addr *int64) (val int64) |
func LoadUint32(addr *uint32) (val uint32)|
func LoadUint64(addr *uint64) (val uint64)|
func LoadUintptr(addr *uintptr) (val uintptr)|
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)	|
func StoreInt32(addr *int32, val int32)|写入操作
func StoreInt64(addr *int64, val int64)|
func StoreUint32(addr *uint32, val uint32)|
func StoreUint64(addr *uint64, val uint64)|
func StoreUintptr(addr *uintptr, val uintptr)|
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)	|
func AddInt32(addr *int32, delta int32) (new int32)|修改操作
func AddInt64(addr *int64, delta int64) (new int64)|
func AddUint32(addr *uint32, delta uint32) (new uint32)|
func AddUint64(addr *uint64, delta uint64) (new uint64)|
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)	|
func SwapInt32(addr *int32, new int32) (old int32)|交换操作
func SwapInt64(addr *int64, new int64) (old int64)|
func SwapUint32(addr *uint32, new uint32) (old uint32)|
func SwapUint64(addr *uint64, new uint64) (old uint64)|
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)|
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)	|
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)| 比较并交换操作
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)|
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)|
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)|
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)|
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)|

#### 示例
比较一下 `sync.Mutex`互斥锁 和 原子操作性能 
```
var x int64
var l sync.Mutex
var wg sync.WaitGroup

// 普通版加函数
func add() {
	// x = x + 1
	x++ // 等价于上面的操作
	wg.Done()
}

// 互斥锁版加函数
func mutexAdd() {
	l.Lock()
	x++
	l.Unlock()
	wg.Done()
}

// 原子操作版加函数
func atomicAdd() {
	atomic.AddInt64(&x, 1)
	wg.Done()
}

func main() {
	start := time.Now()
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		// go add()       // 普通版add函数 不是并发安全的
		// go mutexAdd()  // 加锁版add函数 是并发安全的，但是加锁性能开销大
		go atomicAdd() // 原子操作版add函数 是并发安全，性能优于加锁版
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(x)
	fmt.Println(end.Sub(start))
}
```
