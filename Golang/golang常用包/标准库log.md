- [标准库`log`介绍](#标准库log介绍)
  - [`Logger`使用](#logger使用)
- [配置`logger`](#配置logger)
  - [标准的`logger`的配置](#标准的logger的配置)
  - [`flag`选项](#flag选项)
  - [配置日志前缀](#配置日志前缀)
  - [配置日志输出位置](#配置日志输出位置)
  - [标准设置操作--放进去`init`中](#标准设置操作--放进去init中)
- [创建`logger`](#创建logger)
- [总结](#总结)

# 标准库`log`介绍

`Go`语言内置的`log`包实现简单的日志服务。

## `Logger`使用

`log`包定义了`Logger`类型，提供了一些格式化输出方法。本包也提供了一个预定义的“标准”`logger`，可以通过调用函数`Print系列(Print|Printf|Println）、Fatal系列（Fatal|Fatalf|Fatalln）、和Panic系列（Panic|Panicf|Panicln）`来使用，比自行创建一个`logger`对象更容易使用。

下面例子，默认他们会将日志信息打印到终端界面

    package main

    import (
    	"log"
    )

    func main() {
    	log.Println("这是一个简单的日志，已经带有\\n")
    	v := "很普通的"
    	log.Printf("这是一条%s日志。\n", v)
    	log.Panicln("这是一条会触发panic的日志")	//写入日志信息后，触发panic
    	//log.Fatalln("这是一条会触发fatal的日志")		// 写入日志信息后，调用os.Exit(1)退出
    }

`Fatal`系列函数写入日志信息后调用`os.Exit(1)`，`Panic`系列函数写入日志信息后报错`panic`。

# 配置`logger`

## 标准的`logger`的配置

默认情况下**只会提供日志的时间信息**,需要添加其他信息.

`log`标准库中`Flags`函数会返回标准`logger`的输出配置, 而`SetFlags`函数用来设置`logger`的输出配置.

    func Flags() int        //返回logger输出配置
    func SetFlags(flag int) //设置logger的输出配置,传入flag

## `flag`选项

如下的`flag`选项, 他们是一系列定义好的常量。

    const (
        // 控制输出日志信息细节,不能控制输出顺序和格式
        // 输出日志在每一项后都有一个冒号分割: 例如2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
        Ldate       =1 << iota  //日期:2009/01/23
        Ltime                   //时间:01:23:23
        Lmicroseconds           //微妙级别时间:01:23:23.123123(增加Ltime位)
        Llonglife               //文件全路径名+行号:    /a/b/c/d.go:23
        Lshortfile              //文件名+行号 (会覆盖掉Llonglife)
        LUTC                    //使用UTC时间
        LstdFlags       = Ldate|Ltime   //标准logger的初始值
    )

    //例子
    log.SetFlags(log.Lshortfile | log.Lmicroseconds | log.Ldate)
    log.Println("尝试一下新设置，很普通的日志")

    //输出
    2019/11/21 11:19:31.664236 use_log.go:16: 尝试一下新设置，很普通的日志

## 配置日志前缀

前缀方便检索,配置提供了两个两个方法:

    func Prefix() string        //查看标准logger的输出前缀
    func SetPrefix(prefix string) //设置输出前缀  

    fmt.Printf("前缀为%s\n", log.Prefix())	// 空白
    log.SetPrefix("[小王子]")
    log.Println("这是一条很普通的日志")
    fmt.Printf("前缀为%s\n", log.Prefix())	//前缀为[小王子]

    前缀为
    [小王子]2019/11/21 11:24:35.617695 use_log.go:21: 这是一条很普通的日志
    前缀为[小王子]

## 配置日志输出位置

```golang
func SetOutput(w io.Writer)     //用来设置标准logger输出目的地,默认是标准错误输出。

//例子: 输出到同目录下的 xx.log文件中

// 设置保存在Logfile中
	logFile, err := os.OpenFile("./xx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}

	log.SetOutput(logFile)
	log.Println("这是一条很普通的日志,你信吗")	//保存在文件中
```

## 标准设置操作--放进去`init`中

```golang
func init() {
	//初始化,设置保存在file中
	logFile, err := os.OpenFile("./xx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile|log.Lmicroseconds|log.Ldate)
	log.SetPrefix("[liukaitao]")
}
```

# 创建`logger`

提供一个创建新的`logger`对象的构造函数-`New`,签名如下:

```golang
func New(out io.Writer, prefix string, flag int) *Logger    //返回*Logger类型指针
```

`New`创建一个`Logger`对象,参数`out`设置**日志信息写入的目的地**,参数`prefix`会添加到生成每一条日志前面,参数`flag`定义日志属性。

```golang
func main() {
    logger := log.New(os.Stout, "<New>", log.Lshortfile|log.Ldate|log.Ltime)
    logger.Println("这是自定义的logger记录日志")
}
```

# 总结

实际项目中会选择使用第三方日志库,如`logrus`,`zap`等
