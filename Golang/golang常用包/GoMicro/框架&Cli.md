[TOC]

# 参考文档

<https://github.com/go-micro/go-micro>

Examples参考地址:

<https://github.com/go-micro/examples>

官方项目学习的地址

<https://github.com/GoogleCloudPlatform/microservices-demo>

# 架构学习

参考微信文档: <https://mp.weixin.qq.com/s/MPYxO9uAjVU3AWG-QIEa1Q>

# Go Micro组件架构

> Go Micro为微服务提供了基本的构建模块，其目标是简化分布式系统开发。

因为微服务是一种架构模式，Micro的架构思路是通过工具组件进行拆分，简化我们开发微服务项目的难度；Micro的设计哲学是可插拔的插件化架构。

**Go-micro是微服务的独立RPC框架，也是我们学习go-micro的核心。**

我们看下面的架构图：

![图片](https://mmbiz.qpic.cn/mmbiz_png/lCQLg02gtibsMmmspAnQlpLZGVqIAYAMWOS3OCHWJ1sFEhAFNnz89T0uxNT88tO8X6RvKFJn4ibAM7rgzOKzxRicQ/640?wx_fmt=png\&wxfrom=5\&wx_lazy=1\&wx_co=1 "图片")

*   最顶层是service，代表一个微服务
*   服务下面是两个端：客户端和服务端。（*注意区分Service和Server，我在刚接触的时候混淆了这两个概念，service指的是服务；server服务端包含在service中，和client客户端一起作为service的底层支撑*）

    *   **服务端Server**：==用于构建微服务的接口，提供用于RPC请求的方法。==
    *   **客户端Client**：==提供RPC查询方法，它结合了注册表，选择器，代理和传输。它还提供重试机制，超时机制，使用上下文等，是我们入门阶段的重点。==

架构图中最底层的组件对于初学微服务的同学肯定不熟悉，下面来重点介绍一下：

## Registry注册中心

**注册中心提供可插入的服务发现库，来查找正在运行的服务。** 默认的实现方式是consul。

我们也可以很方便的修改为etcd，kubernetes等。毕竟可插拔是go-micro重要特性

## Selector负载均衡

**Selector选择器实现go-micro的负载均衡机制。**

原理是这样的：**==当客户端向服务器发出请求时，首先查询服务的注册中心，注册中心会返回一个正在运行服务的节点列表，选择器将选择这些节点中的其中一个用于查询请求。==**

多次调用选择器将使用平衡算法。目前的方法是循环法、随机哈希、黑名单。go-micro就是通过这种机制实现负载均衡的。

## Broker事件驱动：发布订阅

**Broker是发布和订阅的可插入接口。**

==微服务是一个事件驱动的架构，发布和订阅事件应该是一流的公民==。目前的实现包括nats，rabbitmq和http。

## Transport消息传输

**传输是通过点对点传输消息的可插拔接口。**

==目前的实现是http，rabbitmq和nats。通过提供这种抽象，运输可以无缝地换出。==

# 通过Example的学习基本使用

主体主要是以下

```go
package main

import (
	"go-micro.dev/v4"
)

func main() {
	// create a new service
	service := micro.NewService(
		micro.Name("helloworld"),	// name of the service
	)

	// initialise flags
	service.Init()		// can use service.Init(micro.option...)

	// start the service
	service.Run()
}

```

## &#x20;加入options

`options` struct

```go
// Options for micro service
type Options struct {
	Auth      auth.Auth
	Broker    broker.Broker
	Cache     cache.Cache
	Cmd       cmd.Cmd
	Config    config.Config
	Client    client.Client
	Server    server.Server
	Store     store.Store
	Registry  registry.Registry
	Runtime   runtime.Runtime
	Transport transport.Transport
	Profile   profile.Profile
	Logger    logger.Logger
	// Before and After funcs
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context

	Signal bool
}
```

```go
package main

import (
	"go-micro.dev/v4"
)

const addr string = ":8080"

func main() {
	// create a new service, Optionally include some options here
	service := micro.NewService(
		// here can use options
		micro.Name("helloworld"),
		micro.Version("v1"),
		micro.Metadata(map[string]string{
			"type": "learn_options",
		}),
		micro.Address(addr),
		// there are lots of options, you can see AfterStart, BeforeStart, Logger
	)
	// there can use options to init service but will overwrite above option
	opts := []micro.Option{
		micro.Name("helloworld2"),
		micro.Version("lastest"),
		micro.Address(addr),
	}

	// initialise flags
	service.Init(opts...)

	// start the service
	service.Run()
}

```

# Go-micro cli使用

最新版本需要 go:1.18

可以使用go-micro 命令行来加快开发我们的项目

参考: <https://github.com/go-micro/cli>

主要使用  new创建服务模板文件&#x20;

```bash
# 安装
go install github.com/go-micro/cli/cmd/go-micro@latest

$ go-micro
NAME:
   C:\Users\Liu\go\bin\go-micro.exe - The Go Micro CLI tool

USAGE:
   go-micro.exe [global options] command [command options] [arguments...]

VERSION:
   v1.1.4

COMMANDS:
   call        Call a service, e.g. C:\Users\Liu\go\bin\go-micro.exe call helloworld Helloworld.Call '{"name": "John"}'     
   completion  Output shell completion code for the specified shell (bash or zsh)
   describe    Describe a resource
   generate    Generate project template files after the fact
   new         Create a project template
   run         Build and run a service continuously, e.g. C:\Users\Liu\go\bin\go-micro.exe run [github.com/auditemarlow/helloworld]
   services    List services in the registry
   stream      Create a service stream
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   ....
   --version, -v                   print the version (default: false)
```

## new 创建项目示例

```bash
# 这里会创建helloworld服务文件夹
go-micro new service helloworld
# 开始初始化&执行mak文件
cd helloworld
make init proto update tidy

---- 可以指定各种插件
go-micro new service --sqlc helloworld
go-micro new service --buildkit helloworld
go-micro new service --grpc helloworld

$ go-micro new service --help
NAME:
   C:\Users\Liu\go\bin\go-micro.exe new service - Create a service template, e.g. C:\Users\Liu\go\bin\go-micro.exe new service [github.com/auditemarlow/]helloworld

USAGE:
   C:\Users\Liu\go\bin\go-micro.exe new service [command options] [arguments...]

OPTIONS:
   --advanced               Generate advanced features in main.go server file (default: false)  
   --buildkit               Use BuildKit features in Dockerfile (default: false)
   --complete               Complete will set the following flags to true; jaeger, health, grpc, sqlc, tern, kustomize, tilt, advanced (default: false)
   --grpc                   Use gRPC as default server and client (default: false)
   --health                 Generate gRPC Health service used for Kubernetes liveliness and readiness probes (default: false)
   --jaeger                 Generate Jaeger tracer files (default: false)
   --kubernetes             Generate Kubernetes resource files (default: false)
   --kustomize              Generate kubernetes resouce files in a kustomize structure (default: false)
   --namespace value        Default namespace for kubernetes resources, defaults to 'default' (default: "default")
   --postgresaddress value  Default postgres address for kubernetes resources, defaults to postgres.database.svc (default: "postgres.database.svc")
   --privaterepo            Amend Dockerfile to build from private repositories (add ssh-agent) 
(default: false)
   --skaffold               Generate Skaffold files (default: false)
   --sqlc                   Generate sqlc resources (default: false)
   --tern                   Generate tern resouces; sql migrations templates (default: false)   
   --tilt                   Generate Tiltfile (default: false
```

### &#x20;Run a Client

```bash
go-micro new client helloworld
go-micro run service_name
```

## call 测试调用其中服务

```bash
#  调用helloworld 服务中的 Call 方法
go-micro call helloworld Helloworld.Call '{"name": "John"}'   
helloworld -- service_name
Helloworld -- 运行的struct
Call	-- function
```

## &#x20;run 运行服务

```bash
# Running A Service
go-micro run
# run with docker
make docker
```

## List services && describe

```bash
go-micro services 
go-micro describe helloworld
# 输出 ,,also can use --format=yaml
{
  "name": "helloworld",
  "version": "latest",
  "metadata": null,
  "endpoints": [
    {
      "name": "Helloworld.Call",
      "request": {
        "name": "CallRequest",
        "type": "CallRequest",
        "values": [
          {
            "name": "name",
            "type": "string",
            "values": null
          }
        ]
      },
      "response": {
        "name": "CallResponse",
        "type": "CallResponse",
        "values": [
          {
            "name": "msg",
            "type": "string",
            "values": null
          }
        ]
      }
    }
  ],
  "nodes": [
    {
      "id": "helloworld-9660f06a-d608-43d9-9f44-e264ff63c554",
      "address": "172.26.165.161:45059",
      "metadata": {
        "broker": "http",
        "protocol": "mucp",
        "registry": "mdns",
        "server": "mucp",
        "transport": "http"
      }
    }
  ]
}
```

# 一个调试Example

测试获取metadata和ip地址

使用到了grpc, proto中的注册RegisterSayHandler

```go
package main

import (
	"fmt"
	"log"
	"time"

	hello "github.com/go-micro/examples/greeter/srv/proto/hello"
	proto "github.com/go-micro/examples/service/proto"
	"go-micro.dev/v4"
	"go-micro.dev/v4/metadata"

	"context"
)

//  rpc实现的方法, 根据proto文件
type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		rsp.Msg = "No metadata received"
		return nil
	}
	log.Printf("Received metadata %v\n", md)
	rsp.Msg = fmt.Sprintf("Hello %s thanks for this %v", req.Name, md)
	return nil
}

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.Request, rsp *proto.Response) error {
	md, _ := metadata.FromContext(ctx)

	// local ip of service
	fmt.Println("local ip is", md["Local"])

	// remote ip of caller
	fmt.Println("remote ip is", md["Remote"])

	rsp.Greeting = "Hello " + req.Name
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	service.Init()

	hello.RegisterSayHandler(service.Server(), new(Say))
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

```

```go
// 启动命令
go-micro run GoMicroPro

// 测试命令
go-micro call greeter Greeter.Hello '{"name":"liukaitao"}'
```

# graceful启动

使用`server.Wait(nil) server.Option`

```go
	service.Server().Init(
		server.Wait(nil),
	)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
```

# 使用Plugins

文档: <https://github.com/go-micro/plugins>

## 使用logrus

插件源码: <https://github.com/go-micro/plugins/tree/main/v4/logger/logrus>

原理: 是插件里的`NewLogger` 初始化`logrus` ,` logrus`实现了`go-micro`的`defaultLogger`接口方法

```go
import (
	...
	"github.com/go-micro/plugins/v4/logger/logrus"
	"go-micro.dev/v4/logger"
)

func main() {
...
	// 使用文件写入
	f, err := os.OpenFile("./log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("END: ", err)
		return
	}
	defer f.Close()

	l := logrus.NewLogger(logger.WithOutput(f))
	flog := l.Fields(map[string]interface{}{
		"name": "Liukaitao",
	})
	// 服务初始化
	service.Init([]micro.Option{
		micro.Logger(flog),
	}...)

...
}
```

## 启用http

```go
package main

import (
	"context"

	helloworld "go-micro-pro/helloworld/proto"

	"github.com/gin-gonic/gin"
	httpServer "github.com/go-micro/plugins/v4/proxy/http"
	"go-micro.dev/v4"
	"go-micro.dev/v4/logger"
	"go-micro.dev/v4/registry"
)

const (
	ServerName = "go.micro.web.DemoHTTP"
)

func main() {
	srv := httpServer.NewService(
		micro.Name(ServerName),
		httpServer.WithBackend("http:localhost:8080"),
	)
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.Use(gin.Recovery())

	// register router
	demo := NewDemo()
	demo.InitRouter(eng)

	hd := srv.Server().NewHandler(eng)
	if err := srv.Server().Handle(hd); err != nil {
		logger.Fatal(err)
	}

	//Create service
	service := micro.NewService(
		micro.Registry(registry.NewRegistry()),
		micro.Server(srv.Server()),
	)
	service.Init()

	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

type demo struct{}

func NewDemo() *demo {
	return &demo{}
}

func (a *demo) InitRouter(router *gin.Engine) {
	router.POST("/demo", a.demo)
}

func (a *demo) demo(c *gin.Context) {
	// create a service
	service := micro.NewService()
	service.Init()

	client := helloworld.NewHelloworldService("go.micro.srv.HelloWorld", service.Client())

	rsp, err := client.Call(context.Background(), &helloworld.CallRequest{
		Name: "world!",
	})
	if err != nil {
		c.JSON(200, gin.H{"code": 500, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"code": 200, "msg": rsp.Msg})
}

```

