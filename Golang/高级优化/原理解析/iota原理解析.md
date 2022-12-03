- [`iota`](#iota)
- [前言](#前言)
- [热身](#热身)
- [规则](#规则)
- [编译原理](#编译原理)


# `iota`

# 前言

iota常用于const表达式

其值是`从零开始`，`const声明块中每增加一行iota值自增1`。

使用iota可以简化常量定义



# 热身

下面常量定义源于GO源码，下面每个常量的值是多少？

```go
type Priority int
const (
    LOG_EMERG Priority = iota
    LOG_ALERT
    LOG_CRIT
    LOG_ERR
    LOG_WARNING
    LOG_NOTICE
    LOG_INFO
    LOG_DEBUG
)
```

上面代码源于日志模块，定义了一组代表日志级别的常量，`常量类型为Priority`，`实际为int类`型。

参考答案：

`iota初始值为0`，也即LOG_EMERG值为0，下面`每个常量递增1。`



下面代码取自Go源码，请问每个常量值是多少？

```go
const (
    mutexLocked = 1 << iota // mutex is locked
    mutexWoken
    mutexStarving
    mutexWaiterShift = iota
    starvationThresholdNs = 1e6
)
```

以上代码取自Go互斥锁Mutex的实现，用于指示各种状态位的地址偏移。

参考答案：

mutexLocked == 1；mutexWoken == 2；mutexStarving == 4；mutexWaiterShift == 3；starvationThresholdNs == 1000000。



请问每个常量值是多少？

```go
const (
    bit0, mask0 = 1 << iota, 1<<iota - 1
    bit1, mask1
    _, _
    bit3, mask3
)
```

参考答案：

bit0 == 1， mask0 == 0， bit1 == 2， mask1 == 1， bit3 == 8， mask3 == 7



# 规则

其实规则只有一条：

- `iota代表了const声明块的行索引（下标从0开始）`



```go
const (
    bit0, mask0 = 1 << iota, 1<<iota - 1   //const声明第0行，即iota==0
    bit1, mask1                            //const声明第1行，即iota==1, 表达式继承上面的语句
    _, _                                   //const声明第2行，即iota==2
    bit3, mask3                            //const声明第3行，即iota==3
)
```



# 编译原理

const块中`每一行在GO中使用spec数据结构描述`，spec声明如下：

```go
            // A ValueSpec node represents a constant or variable declaration
            // (ConstSpec or VarSpec production).
            //
            ValueSpec struct {
                Doc     *CommentGroup // associated documentation; or nil
                Names   []*Ident      // value names (len(Names) > 0)
                Type    Expr          // value type; or nil
                Values  []Expr        // initial values; or nil
                Comment *CommentGroup // line comments; or nil
            }
```

我们只关注ValueSpec.Names， 这个切片中保存了一行中定义的常量，`如果一行定义N个常量，那么ValueSpec.Names切片长度即为N。`



`const块实际上是spec类型的切片，用于表示const中的多行。`



所以编译期间构造常量时的伪算法如下：

```go
    for iota, spec := range ValueSpecs {
        for i, name := range spec.Names {
            obj := NewConst(name, iota...) //此处将iota传入，用于构造常量
            ...
        }
    }
```

从上面可以更清晰的看出`iota实际上是遍历const块的索引`，每行中即便多次使用iota，其值也不会递增。