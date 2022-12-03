- [`range`](#range)
- [热身](#热身)
  - [动态遍历](#动态遍历)
- [实现原理](#实现原理)
- [`range for slice`](#range-for-slice)
- [`range for map`](#range-for-map)
- [`range for channel`](#range-for-channel)
- [总结:](#总结)

# `range`

range是Golang提供的一种迭代遍历手段，可操作的类型有数组、切片、Map、channel等，实际使用频率非常高。

# 热身

## 动态遍历

请问如下程序是否能正常结束？

```go
func main() {
    v := []int{1, 2, 3}
    for i:= range v {
        v = append(v, i)
    }
}
```

程序解释： main()函数中定义一个切片v，通过range遍历v，遍历过程中不断向v中添加新的元素。

参考答案：` 能够正常结束`。循环内改变切片的长度，`不影响循环次数，循环次数在循环开始前就已经确定了。`



# 实现原理

对于for-range语句的实现，可以从编译器源码中找到答案。 编译器源码`gofrontend/go/statements.cc/For_range_statement::do_lower()`方法中有如下注释。

```go
// Arrange to do a loop appropriate for the type.  We will produce
//   for INIT ; COND ; POST {
//           ITER_INIT
//           INDEX = INDEX_TEMP
//           VALUE = VALUE_TEMP // If there is a value
//           original statements
//   }
```

可见range实际上是一个C风格的循环结构。`range支持数组、数组指针、切片、map和channel类型，对于不同类型有些细节上的差异。`



# `range for slice`

下面的注释解释了遍历slice的过程：



```go
 The loop we generate:
   for_temp := range
   len_temp := len(for_temp)
   for index_temp = 0; index_temp < len_temp; index_temp++ {
           value_temp = for_temp[index_temp]
           index = index_temp
           value = value_temp
           original body
   }
```

遍历slice前会`先获取slice的长度len_temp作为循环次数`，循环体中，每次循环会先获取元素值，`如果for-range中接收index和value的话，则会对index和value进行一次赋值。`



`由于循环开始前循环次数就已经确定了，所以循环过程中新添加的元素是没办法遍历到的。`

另外，数组与数组指针的遍历过程与slice基本一致，不再赘述。



# `range for map`

下面的注释解释了遍历map的过程：

```
// The loop we generate:
//   var hiter map_iteration_struct
//   for mapiterinit(type, range, &hiter); hiter.key != nil; mapiternext(&hiter) {
//           index_temp = *hiter.key
//           value_temp = *hiter.val
//           index = index_temp
//           value = value_temp
//           original body
//   }
```

`遍历map时没有指定循环次数`，循环体与遍历slice类似

由于map底层实现与slice不同，`map底层使用hash表实现，插入数据位置是随机的`，`所以遍历过程中新插入的数据不能保证遍历到。`



# `range for channel`

`遍历channel`是最特殊的，这是由`channel的实现机制`决定的：

```
// The loop we generate:
//   for {
//           index_temp, ok_temp = <-range
//           if !ok_temp {	// 如果关闭了返回　ok = false , 解除阻塞
//                   break
//           }
//           index = index_temp
//           original body
//   }
```

`channel遍历是依次从channel中读取数据,读取前是不知道里面有多少个元素的`。



`如果channel中没有元素，则会阻塞等待，如果channel已被关闭，则会解除阻塞并退出循环`。



# 总结:

- `使用index,value接收range返回值会发生一次数据拷贝`
- 遍历channel时`，如果channel中没有数据，可能会阻塞 `,　`因此记得关闭channel，这样range就能解阻塞`

