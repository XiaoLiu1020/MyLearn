- [背景介绍](#背景介绍)
- [`Go Logger`优势和劣势](#go-logger优势和劣势)
		- [优势:](#优势)
		- [劣势:](#劣势)
- [`Uber-go Zap`](#uber-go-zap)
	- [为什么选择它](#为什么选择它)
	- [安装](#安装)
	- [配置`Zap Logger`](#配置zap-logger)
		- [`Logger`](#logger)
		- [`Sugared Logger`](#sugared-logger)
- [定制`logger`](#定制logger)
	- [将日志文件写入文件中](#将日志文件写入文件中)
	- [更改时间编码并添加调用者详细信息](#更改时间编码并添加调用者详细信息)
	- [使用`Lumberjack`进行日志切割归档](#使用lumberjack进行日志切割归档)
		- [安装](#安装-1)
	- [测试所有功能](#测试所有功能)


转载: <https://www.liwenzhou.com/posts/Go/zap/>

# 背景介绍

在许多`Go`项目中,一般是需要一个好的日志记录器能够提供下面这些功能:

*   能够将事件记录到文件中,而不是应用程序控制台
*   日志切割---根据文件大小,时间,或间隔等来切割日志文件
*   支持不同的日志级别。例如`INFO,DEBUG,ERROR`等
*   能够打印更加详细的基本信息

# `Go Logger`优势和劣势

### 优势:

使用简单,可以设置任何`io.Writer`作为日志记录输出饼向其发送要写入的日志

### 劣势:

*   仅限基本日志级别
*   对于错误日志,只有`Fatal`和`Panic`
    *   缺少一个`ERROR`日志级别,可以在不抛出`panic`或退出程序情况下继续记录错误
*   缺乏日志格式化能力
*   不提供日志切割能力

# `Uber-go Zap`

`Zap`是非常快,结构化,分日志级别的`Go`日志库

## 为什么选择它

*   提供结构化日志记录和`printf`风格日志记录
*   快

## 安装

```shell
go get -u go.uber.org/zap
```

## 配置`Zap Logger`

`Zap`提供两种类型记录器-`Sugared Logger`和`Logger`

*   **在性能很好但不是很关键的上下文中**,使用`SugaredLogger`, 比其他结构化日志记录包快4-10倍
*   **在每一微秒和每一次内存分配都很重要的上下文中**,使用`Logger`,内存分配次数更少,但只支持强类型的结构化日志记录

### `Logger`

*   通过调用`zap.NewProduction()/zap.NewDevelopment()`或者`zap.Example()`创建一个`Logger`
*   **上面每一个函数都会创建一个**`Logger`, 唯一区别是它将记录的信息不同
*   通过`Logger`调用`Info/Error`等
*   默认情况下日志都会打印到应用程序的`console`界面

**日志记录器语法:**

```golang
func (log *Logger) MethodXXX(msg string, fields ...Field)
```

*   `MethodXXX` 可变参数函数,可以是`Info/Error/Debug/Panic`等, 每个方法都接受一个消息字符串和任意数量的`zapcore.Field`场参数
*   `zapcore.Field` 就是一组键值对参数

**例子:**

```golang
var logger *zap.Logger

func main() {
    InitLogger()
    defer logger.Sync()
        simpleHttpGet("www.google.com")
        simpleHttpGet("http://www.google.com")
}

func InitLogger() {
    logger, _ = zap.NewProduction()
}

func simpleHttpGet(url string) {
    resp, err := http.Get(url)          //发出请求
    if err != nil {
        logger.Error(
            "Error fetching url...",    //输出信息
            zap.String("url", url),     
            zap.Error(err))
        )
    } else {
        logger.Info(
            "Success..",
            zap.String("statusCode", resp.Status),
            zap.String("url", url))
        resp.Body.Close()
    }
}

// 输出:
{"level":"error","ts":1572159218.912792,"caller":"zap_demo/temp.go:25","msg":"Error fetching url..","url":"www.sogo.com","error":"Get www.sogo.com: unsupported protocol scheme \"\"","stacktrace":"main.simpleHttpGet\n\t/Users/q1mi/zap_demo/temp.go:25\nmain.main\n\t/Users/q1mi/zap_demo/temp.go:14\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

{"level":"info","ts":1572159219.1227388,"caller":"zap_demo/temp.go:30","msg":"Success..","statusCode":"200 OK","url":"http://www.sogo.com"}
```

### `Sugared Logger`

用`Sugared Logger`来实现相同功能

*   唯一区别: 我们通过调用主`Logger`的`.Sugar()`方法获取一个`SugaredLogger`
*   然后使用`SugaredLogger`以`printf`格式记录

下面是修改过后使用`SugaredLogger`代替`Logger`代码:

```golang
var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.google.com")
	simpleHttpGet("http://www.google.com")
}

func InitLogger() {
  logger, _ := zap.NewProduction()
	sugarLogger = logger.Sugar()
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

//输出
{"level":"error","ts":1572159149.923002,"caller":"logic/temp2.go:27","msg":"Error fetching URL www.sogo.com : Error = Get www.sogo.com: unsupported protocol scheme \"\"","stacktrace":"main.simpleHttpGet\n\t/Users/q1mi/zap_demo/logic/temp2.go:27\nmain.main\n\t/Users/q1mi/zap_demo/logic/temp2.go:14\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:203"}

{"level":"info","ts":1572159150.192585,"caller":"logic/temp2.go:29","msg":"Success! statusCode = 200 OK for URL http://www.sogo.com"}
```

\*\*目前为止,两个`logger`都打印输出`JSON`结构格式

# 定制`logger`

## 将日志文件写入文件中

将使用`zap.New(...)`方法手动传递所有配置,而不使用像`zap.NewProduction()` 这样的预置方法创建`logger`

```golang
func New(core zapcore.Core, options ...Option) *Logger

zapcore.Core 也是需要定义的
core := zapcore.NewCore(zapcore.Encoder, zapcore.WriteSyncer, zapcore.LogLevel)
```

`zapcaore.Core`需要三个配置---`Encoder`, `WriteSyncer`, `Loglevel`

1.  **Encoder**: 编码器(**如何写入日志**) 使用开箱即用的`NewJSONEncode()`, 并使用预先设置的`ProductionEncoderConfig()`

```golang
zapcore.NewJSONEncode(zap.NewProductionEncodeConfig())
```

1.  **WriterSyncer**: **指定日志去哪里**, 我们使用`zapcore.AddSync()`函数并且将打开文件句柄传进去

```golang
file, _ := os.Create("./test.log")
writeSyncer := zapcore.AddSync(file)
```

1.  **LogLevel**: 哪种级别日志将被写入

我们将修改上述部分中`Logger`代码,并重写`InitLogger()`方法, 其他方法 `main()/ SimpleHttpGet()`保持不变

```golang
func InitLogger() {
    writeSyncer := getLogWriter()
    encoder := getEncoder()
    core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)   //定义core
    
    logger := zap.New(core)     //使用core定义logger
    sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
    return zapcore.NewJsonEncoder(zap.NewProductionEncoderConfig())     //设定Encoder,日志格式, 也可以选择NewConsoleEncoder
}

func getLogWriter() zapcore.WriteSyncer {
    file, _ := os.Create("./test.log")
    return zapcore.AddSync(file)        //生成WriteSyncer, 决定输出到哪里
}
```

当使用这些修改过的logger配置调用上述部分的`main()`函数时，以下输出将打印在文件——`test.log`中。

```golang
{"level":"debug","ts":1572160754.994731,"msg":"Trying to hit GET request for www.sogo.com"}
{"level":"error","ts":1572160754.994982,"msg":"Error fetching URL www.sogo.com : Error = Get www.sogo.com: unsupported protocol scheme \"\""}
{"level":"debug","ts":1572160754.994996,"msg":"Trying to hit GET request for http://www.sogo.com"}
{"level":"info","ts":1572160757.3755069,"msg":"Success! statusCode = 200 OK for URL http://www.sogo.com"}
```

## 更改时间编码并添加调用者详细信息

*   时间需要用人类可读方式展示
*   调用方函数详细信息没有显示在日志中

覆盖默认的`productionConfig()`,**修改时间编码器,在日志文件中使用大写字母记录日志级别**

```golang
func getEncoder() zapcore.Encoder {
    encoderConfig := zap.NewProductionEncoderConfig()
    //重写
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    encoderConfig.EncoderLevel = zapcore.CapitalLevelEncoder
    return zapcore.NewConsoleEncoder(encoderConfig)
}
```

修改`zap logger`代码,添加将**调用函数信息记录**到日志中的功能, 在`zap.New(...)`函数中添加一个`Option`

```golang
logger := zap.New(core, zap.AddCaller())
```

## 使用`Lumberjack`进行日志切割归档

`Zap`本身不支持切割归档日志文件,需要第三方库`Lumberjack`实现

### 安装

```shell
go get -u github.com/natefinch/lumberjack
```

\###`zap logger`中加入`Lumberjack`
需要修改`WriteSyncer`diamante

```golang
func getLogWriter() zapcore.WriteSyncer{
    lumberJackLogger := &lumberjack.Logger{
        Filename: "./test.log",
        MaxSize:    10,
        MaxBackups: 5,
        MaxAge:     30,
        Compress:   false,
    }
    return zapcore.AddSync(lumberJackLogger)
}
```

`Lumberjack Logger`采用以下属性输入:

*   `Filename`: 日志文件位置
*   `MaxSize`: 在进行切割前,日志文件最大大小(MB单位)
*   `MaxBackups`: 保留旧文件最大个数
*   `MaxAges`: 保留旧文件最大天数
*   `Compress`:是否压缩/归档旧文件

## 测试所有功能

最终,使用`Zap/Lumberjack logger`完整实例代码如下:

```golang
package main

import (
	"net/http"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger

func main() {
	InitLogger()
	defer sugarLogger.Sync()
	simpleHttpGet("www.sogo.com")
	simpleHttpGet("http://www.sogo.com")
}

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func simpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error = %s", url, err)
	} else {
		sugarLogger.Infof("Success! statusCode = %s for URL %s", resp.Status, url)
		resp.Body.Close()
	}
}

//输出到文件,结果如下: test.log
2019-10-27T15:50:32.944+0800	DEBUG	logic/temp2.go:48	Trying to hit GET request for www.sogo.com
2019-10-27T15:50:32.944+0800	ERROR	logic/temp2.go:51	Error fetching URL www.sogo.com : Error = Get www.sogo.com: unsupported protocol scheme ""
2019-10-27T15:50:32.944+0800	DEBUG	logic/temp2.go:48	Trying to hit GET request for http://www.sogo.com
2019-10-27T15:50:33.165+0800	INFO	logic/temp2.go:53	Success! statusCode = 200 OK for URL http://www.sogo.com
```

