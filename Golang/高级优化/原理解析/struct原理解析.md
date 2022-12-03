- [`struct`结构体分析](#struct结构体分析)
- [前言](#前言)
- [`Tag`本质](#tag本质)
  - [`Tag`规则](#tag规则)
  - [`Tag`是`struct`的一部分](#tag是struct的一部分)
  - [获取`Tag`](#获取tag)
  - [`Tag`存在意义](#tag存在意义)

# `struct`结构体分析

# 前言

struct声明允许字段附带`Tag`来对字段做一些标记。

该`Tag`不仅仅是一个字符串那么简单，因为其主要用于反射场景，

`reflect`包中提供了操作`Tag`的方法，所以`Tag`写法也要遵循一定的规则。



# `Tag`本质

## `Tag`规则

`Tag`本身是一个字符串，但字符串中却是：`以空格分隔的 key:value 对`。

- `key`: 必须是非空字符串，字符串不能包含控制字符、空格、引号、冒号。
- `value`: 以双引号标记的字符串
- 注意：冒号前后不能有空格



如下代码所示，如此写没有实际意义，仅用于说明`Tag`规则

```
type Server struct {
    ServerName string `key1: "value1" key11:"value11"`
    ServerIP   string `key2: "value2"`
}
```

上述代码`ServerName`字段的`Tag`包含两个key-value对。`ServerIP`字段的`Tag`只包含一个key-value对。



## `Tag`是`struct`的一部分

`Tag`只有在反射场景中才有用，而反射包中提供了操作`Tag`的方法。

在说方法前，有必要先了解一下`Go是如何管理struct字段`的。



以下是`reflect`包中的类型声明，省略了部分与本文无关的字段。

```go
// A StructField describes a single field in a struct.
type StructField struct {
    // Name is the field name.
    Name string
    ...
    Type      Type      // field type
    Tag       StructTag // field tag string
    ...
}

type StructTag string
```

可见，描述一个结构体成员的结构中包含了`StructTag`，而其本身是一个`string`。

也就是说`Tag`其实是`结构体字段的一个组成部分`。



## 获取`Tag`

`StructTag`提供了`Get(key string) string`方法来获取`Tag`，示例如下：

```go
package main

import (
    "reflect"
    "fmt"
)

type Server struct {
    ServerName string `key1:"value1" key11:"value11"`
    ServerIP   string `key2:"value2"`
}

func main() {
    s := Server{}
    st := reflect.TypeOf(s)

    field1 := st.Field(0)			// 获取struct 的第一个属性的Tag
    fmt.Printf("key1:%v\n", field1.Tag.Get("key1"))		// Tag的第一个key的值
    fmt.Printf("key11:%v\n", field1.Tag.Get("key11"))

    filed2 := st.Field(1)
    fmt.Printf("key2:%v\n", filed2.Tag.Get("key2"))
}
```

输出如下：

```bash
key1:value1
key11:value11
key2:value2
```



## `Tag`存在意义

本文示例中tag没有任何实际意义，这是为了阐述tag的定义与操作方法，也为了避免与你之前见过的诸如`json:xxx`混淆。



使用反射可以动态的给结构体成员赋值，正是因为有tag，`在赋值前可以使用tag来决定赋值的动作。`

官方的`encoding/json`包，可以将一个JSON数据`Unmarshal`进一个结构体，此过程中就使用了Tag。该包定义一些规则，只要参考该规则设置tag就可以将不同的JSON数据转换成结构体。



总之：正是基于struct的tag特性，才有了诸如json、orm等等的应用。理解这个关系是至关重要的。或许，你可以定义另一种tag规则，来处理你特有的数据。