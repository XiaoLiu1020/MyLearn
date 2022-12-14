- [`Go`的文件操作](#go的文件操作)
- [打开和关闭文件](#打开和关闭文件)
- [读取文件](#读取文件)
	- [`file.Read()`](#fileread)
		- [基本使用](#基本使用)
		- [循环读取](#循环读取)
	- [`bufio`读取文件--一行一行](#bufio读取文件--一行一行)
	- [`ioutil`读取整个文件](#ioutil读取整个文件)
- [文件写入操作](#文件写入操作)
	- [`Write`和`WriteString`--固定长度写](#write和writestring--固定长度写)
	- [`bufio.NewWriter`--通过缓存](#bufionewwriter--通过缓存)
	- [`ioutil.WriteFile`--一次性全写入](#ioutilwritefile--一次性全写入)
- [练习](#练习)
	- [`copyFile`](#copyfile)
	- [实现一个`cat`命令](#实现一个cat命令)


# `Go`的文件操作

文件分为文本文件和二进制文件

# 打开和关闭文件

`os.Open()`--打开一个文件,返回一个`*File`和一个`err`,对得到的文件实例调用`close()`方法关闭。

```golang
package main

import (
    "fmt"
    "os"
    )

func main() {
    //只读方式打开当前目录下main.go文件
    file, err := os.Open("./main.go")
    if err != nil {
        fmt.Println("open file failed:", err)
        return
    }
    //关闭文件
    file.Close()        //通常使用defer注册文件关闭语句
}
```

# 读取文件

## `file.Read()`

### 基本使用

```golang
func (f *File) Read(b []byte) (n int, err error)
// 接受字节切片, 返回读取字节数和错误, 读到文件末尾会返回0和io.EOF

func main() {
    // 只读方式打开当前目录下的main.go文件
	file, err := os.Open("./main.go")
	if err != nil {
		fmt.Println("open file failed!, err:", err)
		return
	}
	defer file.Close()
	// 使用Read方法读取数据
	
	var tmp := make([]byte, 128)
	n, err := file.Read(tmp)
	if err == io.EOF {
	    fmt.Println("读取完毕")
	    return
	}
	
	if err != nil {
	    fmt.Println("read file failed, err:", err)
		return
	}
	
	fmt.Printf("读取了%d字节数据\n", n)
	fmt.Println(string(tmp[:n]))
}
```

### 循环读取

使用`for`循环读取文件中所有数据

```golang
func main() {
    //只读方式打开当前目录下main.go文件
    file, err := os.Open("./main.go")
    if err != nil {
        fmt.Println("open file failed!, err: ", err)
        return
    }
    defer file.Close()
    //循环读取文件
    var content []byte
    var tmp = make([]byte, 128)
    for {
        n, err := file.Read(tmp)
        if err == io.EOF {
            fmt.Println("文件读完了")
            break
        }
        if err != nil {
            fmt.Println("read file failed, err:", err)
            return
        }
        content = append(content, tmp[:n]...)
    }
    fmt.Rrintln(string(content))
}
```

## `bufio`读取文件--一行一行

`bufio`是在`file`基础上封装多一层`API`,支持更多功能

```golang
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)


//bufio按行读取示例
func main() {
    file, err := os.Open("./main.go")
    if err != nil {
        fmt.Println("open file failed, err:", err)
        return
    }
    defer file.Close()
    reader := bufio.NewReader(file)
    for {
        line, err := reader.ReadString("\n")    //注意是字符,每行读取
        if err == io.EOF {
            fmt.Println("文件读完了")
            break
        }
        if err != nil {
			fmt.Println("read file failed, err:", err)
			return
		}
		fmt.Print(line)
        
    }
}
```

## `ioutil`读取整个文件

`io/ioutil`包的`ReadFile`方法能够读取完整文件,只需要将文件名作为参数传入.

```golang
package main

import (
    "fmt"
    "io/ioutil"
    )
    
//ioutil.ReadFile读取整个文件
func main() {
    content, err ::= ioutil.ReadFile("./main.go")
    if err != nil {
        fmt.Println("read file failed, err:", err)
		return
    }
    fmt.Println(string(content))
}
```

# 文件写入操作

`os.OpenFile()` 函数能够以指定模式打开文件,从而实现文件写入功能

```golang
func OpenFile(name string, flag int, perm FileMode) (*File, error) {
    ...
}
```

其中:

`name`:为打开文件名,`flag`:打开文件模式

| 模式          | 含义     |
| ------------- | -------- |
| `os.O_WRONLY` | 只写     |
| `os.O_CREATE` | 创建文件 |
| `os.O_RDONLY` | 只读     |
| `os.O_RDWR`   | 读写     |
| `os.O_TRUNC`  | 清空     |
| `os.O_APPEND` | 追加     |

`perm`: 文件权限,一个八进制数, `r`(读) 04, `w`(写)02, `x`(执行) 01

## `Write`和`WriteString`--固定长度写

```golang
func main() {
    file, err := os.OpenFile("xx.txt", os.O_CREATE | os.O_TRUNC | os.WRONLY, 0666)
    if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	str := "hello liukaitao"
	file.Write([]byte(str))             // 写入字节切片数据
	file.WriteString("hello 小王子")    //直接写入字符串数据
}
```

## `bufio.NewWriter`--通过缓存

```golang
func main() {
    file, err := os.OpenFile("xx.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)         //定义一个writer
	for i := 0; i < 10; i++ {
		writer.WriteString("hello沙河\n") //将数据先写入缓存
	}
	writer.Flush() //将缓存中的内容写入文件
}
```

## `ioutil.WriteFile`--一次性全写入

```golang
func main() {
    str := "hello 沙河"
    err := ioutil.WriteFile("./xx.txt", []byte(str),0666)
    if err != nil {
		fmt.Println("write file failed, err:", err)
		return
	}
}
```

# 练习

## `copyFile`

借助`io.Copy()`实现拷贝文件函数

```golang
//CopyFile 拷贝文件函数
func CopyFile(dstName, srcName string)(wirtten int64, err error) {
    //以读方式打开源文件
    src, err := os.Open(srcName)
    if err != nil {
        fmt.Printf("open %s failed, err:%v.\n", srcName, err)
		return err
    }
    defer src.Close()
    // 以写|创建的方式打开目标文件
    dst, err := os.OpenFile(dstName, os.WRONLY | os.O_CREATE, 0644)
    if err != nil {
		fmt.Printf("open %s failed, err:%v.\n", dstName, err)
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)        //调用io.Copy()拷贝内容
}
func main() {
	_, err := CopyFile("dst.txt", "src.txt")
	if err != nil {
		fmt.Println("copy file failed, err:", err)
		return err
	}
	fmt.Println("copy done!")
```

## 实现一个`cat`命令

```golang
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

// cat命令实现
func cat(r *bufio.Reader) {
	for {
		buf, err := r.ReadBytes('\n') //注意是字符
		if err == io.EOF {
			break
		}
		fmt.Fprintf(os.Stdout, "%s", buf)
	}
}

func main() {
	flag.Parse() // 解析命令行参数
	if flag.NArg() == 0 {
		// 如果没有参数默认从标准输入读取内容
		cat(bufio.NewReader(os.Stdin))
	}
	// 依次读取每个指定文件的内容并打印到终端
	for i := 0; i < flag.NArg(); i++ {
		f, err := os.Open(flag.Arg(i))
		if err != nil {
			fmt.Fprintf(os.Stdout, "reading from %s failed, err:%v\n", flag.Arg(i), err)
			continue
		}
		cat(bufio.NewReader(f))
	}
}
```

