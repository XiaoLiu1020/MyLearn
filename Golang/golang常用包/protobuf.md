# `protobuf`介绍
`Protobuf`是`Protocol Buffer`简称，是`Google`公司于2008年开源的一种高效的与平台无关，语言无关，可扩展的数据格式，一般用于`RPC`接口的基础工具

# `protobuf`使用
需要先编写`IDL`文件然后借助专门工具生成指定语言代码，从而实现数据序列化与反序列化过程

大致开发流程：
* `IDL`编写
* 生成指定语言代码
* 序列化与反序列化

# `protobuf`语法
[protobuf3语法指南](https://colobu.com/2017/03/16/Protobuf3-language-guide/)

# 编译器安装 
## `ptotoc`
`protobuf`协议编译器是用c++编写的,根据自己的操作系统下载对应版本的`protoc` 编译器[https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases), 解压后拷贝到`GOPATH/bin`目录下

## `protoc-gen-go`
安装生成`Go`语言代码工具
```golang
go get -u github.com/golang/protobuf/protoc-gen-go
```

## 编写`IDL`代码
在`protobuf_demo/address`目录下新建`person.proto`文件具体内容:
```golang
// 指定使用protobuf版本
// 此处使用v3版本
syntax = "proto3";

// 包名，通过protoc生成go文件
package address;

// 性别类型
// 枚举类型第一个字段必须为0
enum GenderType {
    SECRET = 0;
    FEMALE = 1;
    MALE = 2;
}

// 人
message Person {
    int64 id = 1;
    string name = 2;
    GenderType gender = 3;
    string number = 4;
}

// 联系簿
message ContactBook {
    repeated Person persons = 1;    
}
```

## 生成`go`语言代码
在`protobuf_demo/address`目录下执行以下命令
```bash
address $ protoc --go_out=. ./person.proto
```
当前目录下生成`person.pb.go`文件,`Go`语言代码里使用这个文件。 在`protobuf_demo/main.go` 文件中:
```golang
package main

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"github.com/Q1mi/studygo/code_demo/protobuf_demo/address"
)

// protobuf demo

func main() {
	var cb address.ContactBook

	p1 := address.Person{
		Name:   "小王子",
		Gender: address.GenderType_MALE,
		Number: "7878778",
	}
	fmt.Println(p1)
	cb.Persons = append(cb.Persons, &p1)
	// 序列化
	data, err := proto.Marshal(&p1)
	if err != nil {
		fmt.Printf("marshal failed,err:%v\n", err)
		return
	}
	ioutil.WriteFile("./proto.dat", data, 0644)

	data2, err := ioutil.ReadFile("./proto.dat")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}
	var p2 address.Person
	proto.Unmarshal(data2, &p2)
	fmt.Println(p2)
}
```
转化列表里字段
```
package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"protobuf_try/m/address"
)

// protobuf demo

func main() {
	//var cb address.ContactBook
	//
	//p1 := address.Person{
	//	Name:   "小王子",
	//	Gender: address.GenderType_MALE,
	//	Number: "7878778",
	//}
	//cb.Persons = append(cb.Persons, &p1)
	cb := &address.ContactBook{
		Persons: []*address.Person{
			{
				Name:   "小王子",
				Gender: address.GenderType_MALE,
				Number: "7878778",
			},
			{
				Name:   "liukaitao",
				Gender: address.GenderType_SECRET,
				Number: "123456",
			},
		},
	}
	fmt.Println("cb: ", cb)
	// 序列化
	data, err := proto.Marshal(cb)
	if err != nil {
		fmt.Printf("marshal failed,err:%v\n", err)
		return
	}
	_ = ioutil.WriteFile("./proto.dat", data, 0644)
	data2, err := ioutil.ReadFile("./proto.dat")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return
	}
	cb2 := &address.ContactBook{}

	_ = proto.Unmarshal(data2, cb2)
	fmt.Println(cb2.Persons)
	fmt.Printf("%v \n", cb2)
	for _, person := range cb2.Persons {
		fmt.Println("=========================")
		fmt.Println(person.Name)

	}

}
```

## 定义服务`Service`
如果想要将消息类型用在`RPC`(远程方法调用)系统中,可以在`.proto`文件中定义一个`RPC`服务接口;

例如:定义一个`RPC`服务并具有一个方法,:
```bash
service SearchService {
    rpc Search(SearchRequest) returns (SearchResponse) {}
}