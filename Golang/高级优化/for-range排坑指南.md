- [`for-range`](#for-range)
- [遍历取不到所有元素指针](#遍历取不到所有元素指针)
- [这样遍历中起`goroutine`可以么](#这样遍历中起goroutine可以么)
- [对`map`遍历时删除元素能遍历到吗?](#对map遍历时删除元素能遍历到吗)
- [对`map`遍历时新增元素能遍历到吗?](#对map遍历时新增元素能遍历到吗)
- [遍历会停止吗?](#遍历会停止吗)

# `for-range`
`for-range` 其实是语法糖，内部调用还是 `for` 循环，初始化会拷贝带遍历的列表`（如 array，slice，map）`，然后每次遍历的`v`都是对同一个元素的遍历赋值。也就是说如果直接对`v`取地址，最终只会拿到一个地址，而对应的值就是最后遍历的那个元素所附给`v`的值。

# 遍历取不到所有元素指针
```golang
arr := [2]int{1, 2}
res := []*int{}
for _, v := range arr {
    res = append(res, &v)
}
//expect: 1 2
fmt.Println(*res[0],*res[1])
//but output: 2 2
```
* 取不到, 同样代码对切片`[]int{1,2}`或`map[int]int{1:1, 2:2}`也不符合预期
* `for-range`语法糖 最终只会拿到一个地址,就是最后遍历那个元素所赋给`v`的值
* 就相当于python闭包延迟性一样

伪代码 `for-range`
```golang
 //len_temp := len(range)
// range_temp := range
// for index_temp = 0; index_temp < len_temp; index_temp++ {
//     value_temp = range_temp[index_temp]
//     index = index_temp   //这里不是取值, 是划等号
//     value = value_temp
//     original body
//   }
```

# 这样遍历中起`goroutine`可以么
```golang
var m = []int{1, 2, 3}
for i := range m {
    go func() {
        fmt.Print(i)
    }()
}
//block main 1ms to wait goroutine finished
time.Sleep(time.Millisecond)
```
* 解释如上, 只会获得最后`i`的值
* 以参数方式传入
```golang
for i := range m {
    go func(i int) {
        fmt.Print(i)
    }(i)
}
```
* 使用局部变量拷贝
```golang
for i := range m {
    i := i  # 这里不一定使用i
    go func() {
        fmt.Print(i)
    }()
}
```

# 对`map`遍历时删除元素能遍历到吗?
```golang
var m = map[int]int{1: 1, 2: 2, 3: 3}
//only del key once, and not del the current iteration key
var o sync.Once
for i := range m {
    o.Do(func() {
        for _, key := range []int{1, 2, 3} {
            if key != i {
                fmt.Printf("when iteration key %d, del key %d\n", i, key)
                delete(m, key)
                break
            }
        }
    })
    fmt.Printf("%d%d ", i, m[i])
}
```
* **不会**, `map`内部实现是链式`hash`表,初始化会**随机从一个遍历开始的位置**

# 对`map`遍历时新增元素能遍历到吗?
```golang
var m = map[int]int{1:1, 2:2, 3:3}
for i, _ := range m {
    m[4] = 4
    fmt.Printf("%d%d ", i, m[i])
}
```
* **可能会**, 因为位置也是从随机位置开始, 输出中可能有`44`,

# 遍历会停止吗?
```golang
v := []int{1, 2, 3}
for i := range v {
    v = append(v, i)
}
```
* **会** , 遍历前会对`v`进行拷贝,期间对原来`v`修改不会反映到遍历中