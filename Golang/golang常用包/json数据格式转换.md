- [`Go json`包](#go-json包)
- [构建`json`数据](#构建json数据)
- [解析数据](#解析数据)
  - [解析`json`数据到`struct`---(结构已知)](#解析json数据到struct---结构已知)
  - [解析`json`到`interface`--(结构未知)](#解析json到interface--结构未知)
  - [还可以解析，创建`json`流式数据](#还可以解析创建json流式数据)

# `Go json`包
```golang
func Marshal(v interface{}) ([]byte, error)
func Unmarshal(data []byte, v interface{}) error
```

# 构建`json`数据
* `struct, slice, array, map`都可以转换成`json`
* 只有字段首字母大写的才会被转换
* `map`转换时候，`key`必须为`string`
* 封装时候，如果是指针，会追踪指针指向的对象进行封装

```golang
package main

import (
	"encoding/json"
	"fmt"
)

//只有struct可以使用tag ``
type Post struct {
	Id			int			`json:"ID"`
	Content	string		`json:"content"`
	Author		string		`json:"author"`
	Label		[]string		`json:"label"`
}

func main() {
	post := &Post{
		Id: 1,
		Content: "Hello World",
		Author: "liukaitao",
		Label: []string{"linux", "shell"},
	}
	//或者使用 json.Marshal(post)
	b, err := json.MarshalIndent(post, "hi", "\t")
	if err != nil {
		fmt.Println("There is something wrong")
		return
	}
	// Marshal()返回是 []byte
	fmt.Println(string(b))

	// slice -> json
	s := []string{"a", "b", "c"}
	d, _ := json.MarshalIndent(s, "", "\t")
	fmt.Println(string(d))

	//map -> json
	m := map[string]string{
		"Author": "liukaitao",
		"age": "18",
		"address": "guangzhou",
	}
	jm, _ := json.MarshalIndent(m, "", "\t")
	fmt.Println(string(jm))
}
```

使用`struct tag`时候，注意点：
* `tag` 可以设置为`json:"-"` 表示本字段不转换为`json`数据，即使是大写字母开头
* `json:"label, omitempty"` omitempty选项，表示这个字段值为0值时候，不转换为`json`格式
* 如果字段的类型为`bool、string、int类、float`类，而tag中又带有,`string`选项，那么这个字段的值将转换成`json`字符串
```golang
type Post struct {
    Id      int      `json:"ID,string"`
    Content string   `json:"content"`
    Author  string   `json:"author"`
    Label   []string `json:"label,omitempty"`
}
```

# 解析数据
## 解析`json`数据到`struct`---(结构已知)
`json`数据可以解析到`struct`或空接口`interface{}`中(也可以`是slice、map等`)。
```golang
{
    "id": 1,
    "content": "hello world",
    "author": {
        "id": 2,
        "name": "userA"
    },
    "published": true,
    "label": [],
    "nextPost": null,
    "comments": [{
            "id": 3,
            "content": "good post1",
            "author": "userB"
        },
        {
            "id": 4,
            "content": "good post2",
            "author": "userC"
        }
    ]
}
//需要根据里面数据分析
type Post struct {
    ID        int64         `json:"id"`       
    Content   string        `json:"content"`  
    Author    Author        `json:"author"`   
    Published bool          `json:"published"`
    Label     []string      `json:"label"`    
    NextPost  *Post         `json:"nextPost"` 
    Comments  []*Comment    `json:"comments"` 
}

type Author struct {
    ID   int64  `json:"id"`  
    Name string `json:"name"`
}

type Comment struct {
    ID      int64  `json:"id"`     
    Content string `json:"content"`
    Author  string `json:"author"` 
}

//解析过程
func main() {
    // 打开json文件
    fh, err := os.Open("a.json")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer fh.Close()
    // 读取json文件，保存到jsonData中
    jsonData, err := ioutil.ReadAll(fh)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    var post Post
    // 解析json数据到post中
    err = json.Unmarshal(jsonData, &post)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(post)
}

//输出结果
{1 hello world {2 userA} true [] <nil> [0xc042072300 0xc0420723c0]}
```

## 解析`json`到`interface`--(结构未知)
如果`json`结构是未知的或者结构可能会发生改变的情况，则解析到`struct`是不合理的。这时可以解析到空接口`interface{}`或`map[string]interface{}`类型上，这两种类型的结果是完全一致的。

解析到`interface{}`上时，`Go`类型和`json`类型的对应关系:
```
JSON类型             Go类型                
---------------------------------------------
JSON objects    <-->  map[string]interface{} 
JSON arrays     <-->  []interface{}          
JSON booleans   <-->  bool                   
JSON numbers    <-->  float64                
JSON strings    <-->  string                 
JSON null       <-->  nil          
```

```golang
func main() {
    // 读取json文件
    fh, err := os.Open("a.json")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer fh.Close()
    jsonData, err := ioutil.ReadAll(fh)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    // 定义空接口接收解析后的json数据
    var unknown interface{}
    // 或：map[string]interface{} 结果是完全一样的
    err = json.Unmarshal(jsonData, &unknown)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(unknown)
}

//输出结果经过格式化
map[
    nextPost:<nil>
    comments:[
        map[
            id:3
            content:good post1
            author:userB
        ]
        map[
            id:4
            content:good post2
            author:userC
        ]
    ]
    id:1
    content:hello world
    author:map[
        id:2
        name:userA
    ]
    published:true
    label:[]
]

//现在可以从这个map去判断类型，取得对应的值，使用类型断言
// 进行断言，并switch匹配
    m := unknown.(map[string]interface{})
    for k, v := range m {
        switch vv := v.(type) {
        case string:
            fmt.Println(k, "type: string\nvalue: ", vv)
            fmt.Println("------------------")
        case float64:
            fmt.Println(k, "type: float64\nvalue: ", vv)
            fmt.Println("------------------")
        case bool:
            fmt.Println(k, "type: bool\nvalue: ", vv)
            fmt.Println("------------------")
        case map[string]interface{}:
            fmt.Println(k, "type: map[string]interface{}\nvalue: ", vv)
            for i, j := range vv {
                fmt.Println(i,": ",j)
            }
            fmt.Println("------------------")
        case []interface{}:
            fmt.Println(k, "type: []interface{}\nvalue: ", vv)
            for key, value := range vv {
                fmt.Println(key, ": ", value)
            }
            fmt.Println("------------------")
        default:
            fmt.Println(k, "type: nil\nvalue: ", vv)
            fmt.Println("------------------")
        }
    }
```

## 还可以解析，创建`json`流式数据
* `type Decoder`解码`json`到`Go`数据结构
* `type Encoder`编码`Go`数据结构到`json

```golang
const jsonStream = `
    {"Name": "Ed", "Text": "Knock knock."}
    {"Name": "Sam", "Text": "Who's there?"}
    {"Name": "Ed", "Text": "Go fmt."}
    {"Name": "Sam", "Text": "Go fmt who?"}
    {"Name": "Ed", "Text": "Go fmt yourself!"}
`
type Message struct {
    Name, Text string
}
func main() {
    dec := json.NewDecoder(strings.NewReader(jsonStream))
    for {
        var m Message
        if err := dec.Decode(&m); err == io.EOF {
            break
        } else if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s: %s\n", m.Name, m.Text)
    }
}

//输出
Ed: Knock knock.
Sam: Who's there?
Ed: Go fmt.
Sam: Go fmt who?
Ed: Go fmt yourself!

//再例如，从标准输入读json数据，解码后处理删除，最后编码重新返回
func main() {
    dec := json.NewDecoder(os.Stdin)
    enc := json.NewEncoder(os.Stdout)
    for {
        var v map[string]interface{}
        if err := dec.Decode(&v); err != nil {
            log.Println(err)
            return
        }
        for k := range v {
            if k != "Name" {
                delete(v, k)
            }
        }
        if err := enc.Encode(&v); err != nil {
            log.Println(err)
        }
    }
}
```