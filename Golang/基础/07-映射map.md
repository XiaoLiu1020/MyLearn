# `Go`语言基础之`map`
映射关系的容器：`map`，其内部使用`散列表(hash)`实现

`map`是一种无序的 `key-value`数据结构，为**引用类型**，必须初始化

### `map`定义
```
map[KeyType]ValueType
```
* `KeyType` 表示键类型
* `ValueType` 表示键对应值类型

`map`**类型初始值为**`nil`，需要使用`make()`函数分配内存-容量
```
make(map[Keytype]ValueType, [cap])
// cap非必须，可以初始化时候就指定合适容量
```
### `map`使用
`map`中数据都是成对使用的
```
func main() {
    scoreMap := make(map[string]int, 8])
    scoreMap["张三"] = 90
    scoreMap["小明"] = 100
    fmt.Println(scoreMap)
    fmt.Println(scoreMap["张三"])
    fmt.Printf("type of a:%T\n", scoreMap)
}


//输出
map[小明:100 张三:90]
90
type of a:map[string]int

//也支持声明时候直接填充
func main() {
	userInfo := map[string]string{
		"username": "沙河小王子",
		"password": "123456",
	}
	fmt.Println(userInfo) //
}

//输出
map[username:沙河小王子 password:123456]
```

### 判断某个键是否存在
```
value, ok := map[key]
```
举个例子：
```
func main() {
    scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	// 如果key存在ok为true,v为对应的值；不存在ok为false,v为值类型的零值
	v, ok := scoreMap["张三"]   //返回 value 与 寻找结果
	if ok {                     //判断是否存在
		fmt.Println(v)
	} else {
		fmt.Println("查无此人")
	}
}
```

### `map`的遍历
使用 `for range` 遍历 `map`
```
func main() {
	scoreMap := make(map[string]int)
	scoreMap["张三"] = 90
	scoreMap["小明"] = 100
	scoreMap["娜扎"] = 60
	for k, v := range scoreMap {        // for key, value := range map { }
		fmt.Println(k, v)
	}
	// 也可以只遍历key
	for k := range scoreMap{
	    fmt.Println(k)
	}
}
```
### 使用 `delete`删除键值对
```
delete(map, key)
```
* `map`要删除键值对的 `map`
* `key`删除键值对的键

#### 按照指定顺序遍历`map`
```
func main() {
	rand.Seed(time.Now().UnixNano()) //初始化随机数种子

	var scoreMap = make(map[string]int, 200)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("stu%02d", i) //生成stu开头的字符串
		value := rand.Intn(100)          //生成0~99的随机整数
		scoreMap[key] = value
	}
	//取出map中的所有key存入切片keys
	var keys = make([]string, 0, 200)
	for key := range scoreMap {
		keys = append(keys, key)
	}
	//对切片进行排序
	sort.Strings(keys)
	//按照排序后的key遍历map
	for _, key := range keys {
		fmt.Println(key, scoreMap[key])
	}
}
```

## 元素为`map`类型的`Slice`-- 切片`Slice`嵌套字典`map`
```
func main() {
    //先定义 切片Slice
	var mapSlice = make([]map[string]string, 3)
	for index, value := range mapSlice {
		fmt.Printf("index:%d value:%v\n", index, value)
	}
	fmt.Println("after init")
	// 对切片中的map元素进行初始化
	mapSlice[0] = make(map[string]string, 10)
	mapSlice[0]["name"] = "小王子"
	mapSlice[0]["password"] = "123456"
	mapSlice[0]["address"] = "沙河"
	for index, value := range mapSlice {
		fmt.Printf("index:%d value:%v\n", index, value)
	}
}
```
## 值为`Slice`切片类型的`map` -- 字典`map`嵌套切片`Slice`
```
func main() {
    // 定义map字典
	var sliceMap = make(map[string]string, 3)
	fmt.Println(sliceMap)
	fmt.Println("after init")
	key := "中国"
	value, ok := sliceMap[key]
	if !ok {
		value = make([]string, 0, 2)    //值为Slice切片
	}
	value = append(value, "北京", "上海")
	sliceMap[key] = value
	fmt.Println(sliceMap)
}
```