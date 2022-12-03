package main

import (
	"fmt"
	"strings"
	"time"
	"user/handlers"
	"user/model"

	services "user/services/proto"

	"github.com/go-micro/plugins/v4/registry/etcd"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"gopkg.in/ini.v1"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

// Init DB
func InitDB() error {
	file, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("配置文件出错: Err: ", err)
		return err
	}
	LoadMysqlData(file)
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	return model.Database(path)

}

func LoadMysqlData(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func main() {
	// 示例
	/* if err := InitDB(); err != nil {
		fmt.Println("There is some problem in InitDB")
		panic(err)
	} */
	// etcd注册, 需要用账号密码
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
		etcd.Auth("root", "password"),
	)

	// 得到微服务
	userService := micro.NewService(
		micro.Name("rpcUserService"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(etcdReg),
		micro.RegisterTTL(30*time.Second),
		micro.RegisterInterval(10*time.Second),
	)
	userService.Init()

	// 服务注册
	services.RegisterUserServiceHandler(
		userService.Server(),
		new(handlers.UserService),
	)

	//
	userService.Run()
}
