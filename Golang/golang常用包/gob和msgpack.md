
- [二进制协议`gob`和`msgpack`](#二进制协议gob和msgpack)
	- [`json`序列化问题](#json序列化问题)
- [`gob`序列化实例](#gob序列化实例)
- [`msgpack`](#msgpack)
	- [安装](#安装)
	- [示例](#示例)

# 二进制协议`gob`和`msgpack`
有个问题： `go`语言中`json`包在序列化空接口存放的数字类型(整型，浮点型等)都序列化成`float64`类型

我们先构造一个结构体
```golang
type s struct {
    data map[string]interface{}
}
```
## `json`序列化问题
```golang
func jsonDemo() {
    var s1 = s {
        data : make(map[string]interface{}, 8),
    }
    s1.data["count"] = 1
    ret, err := json.Marshal(s1.data)
    if err != nil {
        fmt.Println("marshal failed ",err)
    }
    fmt.Printf("%#v\n", string(ret))
    
    var s2 = s {
        data : make(map[string]interface{}, 8),
    }
    err = json.Unmarshal(ret, &s2.data)
    if err != nil {
        fmt.Println("unmarshal failed ", err)
    }
    fmt.Println(s2)
    for _, v := range s2.data {
        fmt.Printf("value:%v, type:%T\n", v, v)
    }
}

//输出结果
"{\"count\":1}"     //s1
{map[count:1]}      //s2
value:1, type:float64   // 1变为float64类型了

```
# `gob`序列化实例
`gob`是`golang`提供的"私有"的编解码方式,效率会比`json,xml`高,适合`go`**语言程序间传递数据**
```golang
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

func main() {
	gobDemo()
}

type s struct {
	data map[string]interface{}
}

func gobDemo() {
	var s1 = s{
		data : make(map[string]interface{}, 8),
	}
	s1.data["count"] = 1
	//encode
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)  //gob编码器需要一个临时的缓存区域
	err := enc.Encode(s1.data)	//把s1.data数据 使用编码器，编码，数据会存到buf缓存区
	if err != nil {
		fmt.Println("gob encode failed, err", err)
		return
	}
	b := buf.Bytes()		//放入缓存区的数据，转化成byte，赋值给b
	fmt.Println("b: ", b)
	var s2 = s{
		data : make(map[string]interface{}, 8),
	}
	//decode
	/*
	或者
	byteEn := buf.Bytes()       //把上面数据赋值byteEn
	decoder := gob.NewDecoder(bytes.NewReader(byteEn))  //起一个新读取器，放入新的解码器
	err := decoder.Decode(&s2.data)     //调用解吗方法
	*/
	dec := gob.NewDecoder(bytes.NewBuffer((b)))		//创建解码器， 再起一个缓存区存入b数据，交给Decoder
	err = dec.Decode(&s2.data)							// Decoder中对b新建的缓存区 按照s2.data格式进行解码
	if err != nil {
		fmt.Println("go decode failed, err: ", err)
		return
	}
	fmt.Println("s2.data : ", s2.data)
	for _, v := range s2.data {
		fmt.Printf("value: %v, type: %T\n", v, v)
	}
}
```

# `msgpack`
`MessagePack`是一种高效的二进制序列化格式, 允许你在多种语言(如`json)之间交换数据,但更快更小
## 安装
```bash
go get -u github.com/vmihailenco/msgpack
```
## 示例
```golang
package main

import (
	"fmt"
	"github.com/vmihailenco/msgpack"
	"log"
)

type Person struct {
	Name  string
	Age 	int
	Gender	string
}

func main() {
	p1 := Person{
		Name: 		"刘凯涛",
		Age:		18,
		Gender:		"boy",
	}
	//marshal
	b, err := msgpack.Marshal(p1)
	checkError("Marshal", err)
	fmt.Printf("b : %v\n", b)   // [131 164 78 97 109 101 169 229 136 152 229 135 175 230 182 155 163 65 103
 101 211 0 0 0 0 0 0 0 18 166 71 101 110 100 101 114 163 98 111 121]
	//unmarshal
	var p2 Person
	err = msgpack.Unmarshal(b, &p2)
	checkError("Unmarshal", err)
	fmt.Printf("p2: %#v\n", p2) //p2: main.Person{Name:"刘凯涛", Age:18, Gender:"boy"}
	fmt.Printf("P2.Name: %s\n", p2.Name)    //P2.Name: 刘凯涛
}

func checkError(prefix string, e error) {
	if e != nil {
		fmt.Printf("%v failed, err:%v\n",prefix, e)
		log.Fatal()
	}
}
```