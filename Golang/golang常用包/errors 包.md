## 使用

```golang
go get "github.com/pkg/errors"
```

## 函数

```golang

package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	//"io/ioutil"
)

func main() {
	_, err := ioutil.ReadFile("./nothing.json")
	if err != nil {
		fmt.Println(errors.Wrap(err, "[read failed]"))      // 添加前缀
		return
	}
	fmt.Println("end")

	fmt.Printf("%v\n", err)
	err = fmt.Errorf("second level err\n %w", err) // 也是添加前缀
	fmt.Printf("%+v\n", err)    

	// 可以调用栈信息，并且带有前缀, 传入参数 format 和 arg ，格式化
	// 会引起当前报错，返回堆栈信息
	myerr := errors.Errorf("whoops: %s", err)
	fmt.Printf("%+v", myerr)            // +v 返回详细堆栈信息
	fmt.Println("============")

	fmt.Println(err)
	fmt.Println(errors.Cause(err))      // 获取最根本错误原因

	myerr2 := errors.WithStack(err)
	fmt.Println(myerr2)

}
```

# 最佳实践

场景: 当请求报错时候,想要暴露给前端的只有比较语义化简单的错误, 但是还需要日志中记录详细的错误日志,这样既能优美返回错误,又能给开发人员找到根本问题

因此可以使用自定义错误返回,捕获报错返回时候进行判断, 把详细错误记录下

参考文档: [https://github.com/Mikaelemmmm/go-zero-looklook/blob/main/doc/chinese/10-%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86.md]()

需要用到`errors.Wrapf(err, s string, arg...) 和 errors.Cause(err error)`

`Cause` 会返回最底层的`err`

```go
// main.go
package main

import (
	"os"

	LiuErr "liu/errors"

	"log"

	"github.com/pkg/errors"
)

func Open() error {
	f, err := os.OpenFile("./unknown", os.O_RDONLY, 0444)
	if err != nil {
		err = os.ErrNotExist
		return errors.Wrapf(err, "There are something wrong in file, %s", "haha")
	}
	f.Close()
	return LiuErr.NewErrCodeMsg(LiuErr.SERVER_COMMON_ERROR, "you cannot open this file")
}

func main() {
	err := Open()
	if err != nil {
		causeErr := errors.Cause(err)
		if _, ok := causeErr.(*LiuErr.CodeError); ok { //自定义错误类型
			log.Printf("[LiuErr.CodeError] %+v", err) // [LiuErr.CodeError] ErrCode:100001，ErrMsg:you cannot open this file
		} else {
			// 这里可以把err记录 留下, 把causeErr反馈就行
			log.Print(err)      // There are something wrong in file, haha: file does not exist
			log.Print(causeErr) // file does not exist
		}
	}

}

// errors/errors.go
package LiuErr

import (
	"fmt"
)

//全局错误码
const SERVER_COMMON_ERROR uint32 = 100001
const REUQEST_PARAM_ERROR uint32 = 100002
const TOKEN_EXPIRE_ERROR uint32 = 100003
const TOKEN_GENERATE_ERROR uint32 = 100004
const DB_ERROR uint32 = 100005

/**
常用通用固定错误
*/

type CodeError struct {
	errCode uint32
	errMsg  string
}

//返回给前端的错误码
func (e *CodeError) GetErrCode() uint32 {
	return e.errCode
}

//返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: SERVER_COMMON_ERROR, errMsg: errMsg}
}

```



