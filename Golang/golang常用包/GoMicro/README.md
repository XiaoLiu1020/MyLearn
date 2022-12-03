- [参考项目](#参考项目)
- [背景](#背景)
- [涉及技术](#涉及技术)
- [关键代码](#关键代码)
  - [`main.go`](#maingo)
  - [`router.go`](#routergo)
  - [`wrappers.go`](#wrappersgo)

# 参考项目
https://github.com/CocaineCong/micro-todoList/tree/af1a892cf4ed946e3f1654296f264198c743b350

# 背景
在学习中自己找了部分资料,但是很多都是v3以前的版本, 所以自己也写了个差不多的,只是用了go-micro的v4版本

# 涉及技术
也是参考项目,使用的`go-micro``registry`是`etcd`, 熔断使用的框架是`go-hystrix`, 使用了`gin`的路由作为`gateway`的http路由, 转发到各个rpc服务,比如user

相关文档: 

https://github.com/go-micro/plugins/v4/registry/etcd

https://github.com/afex/hystrix-go/hystrix

# 关键代码
## `main.go`
```go

//  gateway/main.go

    ...
    etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		etcd.Auth("root", "password"),
	)

	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		// use the hystrix in the userWrapper
		// type Wrapper func(Client) Client
		micro.WrapClient(wrappers.NewUserWrapper),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

// 开启user的grpc服务
	userSerivce := services.NewUserService("rpcUserService", userMicroService.Client())

	server := web.NewService(
		web.Name("httpService"),
		web.Address("127.0.0.1:4000"),
		web.Registry(etcdReg),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*10),
		web.Metadata(map[string]string{"protocol": "http"}),
		// 将服务调用实例使用gin处理
		// web.Handler(h http.Handler)
		web.Handler(router.NewRouter(userSerivce)),
	)

	//启动
	server.Init()
	server.Run()
...
```
## `router.go`

```go
// gateway/gin-router/router/router.go

func NewRouter(service ...interface{}) *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware.Cors(), middleware.InitMiddleware(service), middleware.ErrorMiddleware())
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))
	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		userV1 := v1.Group("/user")
		// 用户服务
		userV1.GET("register", handlers.UserRegister)
		userV1.GET("login", handlers.UserLogin)
	}
	return ginRouter
}
```
## `wrappers.go`
可以参考官方文档 https://github.com/go-micro/plugins/tree/main/v4 来使用插件

```go

// gateway/wrappers/wrappers.go

package wrappers

import (
	"context"
	"fmt"

	"github.com/afex/hystrix-go/hystrix"
	"go-micro.dev/v4/client"
)

// Can see https://github.com/go-micro/plugins/tree/main/v4/wrapper#client-wrapper-usage

type userWrapper struct {
	client.Client
}

func NewUserWrapper(c client.Client) client.Client {
	return &userWrapper{c}
}

func (wrapper *userWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	config := hystrix.CommandConfig{
		Timeout:                hystrix.DefaultTimeout,
		RequestVolumeThreshold: 2,    //熔断器请求阈值，默认20，意思是有20个请求才能进行错误百分比计算
		ErrorPercentThreshold:  50,   //错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次 开启，默认50%
		SleepWindow:            5000, //过多长时间，熔断器再次检测是否开启，单位毫秒ms（默认5秒）
	}
	hystrix.ConfigureCommand(cmdName, config)
	return hystrix.Do(cmdName, func() error {
		// 通过熔断器检测执行以下
		fmt.Printf("Run cmdName: %v \n", cmdName)
		return wrapper.Client.Call(ctx, req, resp)
	}, func(err error) error {
		if err != nil {
			// 熔断器触发报错:
			// cmdName: rpcUserService.UserService.UserLogin Err: hystrix: circuit open
			// 5s 之后重新检测,会报服务的错误
			//cmdName: rpcUserService.UserService.UserLogin Err: some problem in UserLogin service
			fmt.Printf("cmdName: %v Err: %v \n", cmdName, err)
		}
		return err
	})
}

```