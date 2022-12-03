package main

import (
	"api-gateway/gin-router/router"
	services "api-gateway/services/proto"
	"api-gateway/wrappers"
	"time"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

func main() {
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
}
