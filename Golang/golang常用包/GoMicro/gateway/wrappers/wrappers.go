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
