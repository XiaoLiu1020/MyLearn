- [简单的转换操作](#简单的转换操作)
- [`strconv`包](#strconv包)
  - [`string`和`int`转换](#string和int转换)
  - [`Parse`类函数](#parse类函数)
  - [`Format`类函数](#format类函数)
  - [`Append`类函数](#append类函数)


# 简单的转换操作
```
valueOfTypeB = typeB(valueOfTypeA)
```
例如：
```
// 浮点数
a := 5.0
//　转换为int类型
b := int(a)
```

`Go`允许在底层结构相同两个类型之间互转：

```
//　IT类型底层是int类型
type IT int

// a类型为ＩＴ
var a IT = 5

//将(IT)转为int, b 现在是int
b := int(5)

//将b(int) 转为为ITz

c := IT(b)
```

# `strconv`包
`strconv`包提供了字符串与简单数据类型之间的类型转换功能

包提供很多函数，大概分几类:
* 字符串转`int`: `Atoi()`
* `int`转字符串: `Itoa()`
* `ParseTp`类函数将`string`转换为`TP`类型：`ParseBool(),ParseFloat(),ParseInt(),ParseUnit()`，返回第二个值为err
* `FormatTp`类函数用于其他类型转`string`: `FormatBool(), FormatFloat(), FormatInt(), FormatUnix()`
* `AppendTp`类函数用于将`Tp`转换成字符串后`append`到一个`slice`中：　`AppendBool(),AppendFloat(),AppendInt(),AppendUint()`

类型无法转换时，报两种错误
```
var ErrRange = errors.New("value out of range")
var ErrSyntax = erros.New("invalid syntax")
```

## `string`和`int`转换
**int转换为字符串: Itoa()**
```
// Itoa(): int -> string
prntln('a' + strconv.Itoa(32))      // a32
```

**string转换为int: Atoi()**
```
func Atoi(s string) (int, error)
```
`string`可能无法转换为`int`所以有两个返回值，第二个是返回值判断是否转换成功

```
//Atoi(): string -> int
i, _ := strconv.Atoi("3")

//Atoi()转换失败
i, err := strconv.Atoi("a")
if err != nil {
    println("converted failed")
}
```

## `Parse`类函数
**转换字符串为给定类型的值**

由于字符串转换为其它类型可能会失败，所以这些函数都有两个返回值，第一个返回值保存转换后的值，第二个返回值判断是否转换成功。
```
b, err := strconv.ParseBool("true")
f, err := strconv.ParseFloat("3.1415", 64)

//ParseInt 和 ParseUint有三个参数
func ParseInt(s string, base int, bitSize int) (i int64, err error)
func ParseUint(s string, base int, bitSize int) (uint64, error)
// bitSize 参数表示转换为什么位的int/uint
// base 参数表示以什么进制方式去解析给定字符串，有效值为0,2-36

i, err := strconv.ParseInt("-42", 10, 64)   // 十进制解析,转为int64
i, err := strconv.ParseInt("23", 5, 64)   // 5进制解析,转为int64 ,结果为13
u, err := strconv.ParseUint("42", 10, 64)
```

## `Format`类函数
**将给定类型格式化为string**: `FormatBool(),FormatFloat(),FormaInt(),FormatUint()`
```
s := strconv.FormatBool(true)

//FormatFloat()参数众多：
func FormatFloat(f float64, fmt byte, prec, bitSize int) string
//fmt 表示格式
//prec 控制精度，表示小数点后数字个数
//bitSize 表示f来源类型，32:float32


s := strconv.FormatFloat(3.1415, 'E', -1, 64)

// FormatInt()和FormatUint()有两个参数：
func FormatInt(i int64, base int) string
func FormatUint(i uint64, base int) string
// base指定将第一个参数转为多少进制

s := strconv.FormatInt(-42, 16)
s := strconv.FormatUint(42, 16)
```

## `Append`类函数
**将tp转换成字符串后append到一个slice中**: `AppendBool(), AppendFloat(), AppendInt(), AppendUint()`
```
package main

import (
    "fmt"
    "strconv"
    )

func main() {
    // 声明一个slice
    b10 := []byte("int (base 10):")
    
    // 将转换为10进制的string, 追加到slice
    b10 = strconv.AppendInt(b10, -42, 10)
    fmt.Println(string(b10))
    
    b16 := []byte("int (base 16):")
    b16 = strconv.AppendInt(b16, -42, 16)   //转换为16进制追加
    fmt.Println(string(b16))
}

//结果
int (base 10):-42
int (base 16):-2a
```