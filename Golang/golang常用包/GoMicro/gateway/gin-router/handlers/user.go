package handlers

import (
	"context"
	"errors"
	"net/http"

	services "api-gateway/services/proto"

	"github.com/gin-gonic/gin"
	"github.com/go-micro/plugins/v4/logger/logrus"
	"go-micro.dev/v4/logger"
)

var log logger.Logger

func init() {
	// 使用logrus
	log = logrus.NewLogger()
}

func UserRegister(c *gin.Context) {
	var userReq services.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	// 从gin.Key中取出服务实例
	userService := c.Keys["userService"].(services.UserService)
	userResp, err := userService.UserRegister(context.Background(), &userReq)
	PanicIfUserError(err)
	c.JSON(http.StatusOK, gin.H{"data": userResp})
}

func UserLogin(c *gin.Context) {
	var userReq services.UserRequest
	PanicIfUserError(c.Bind(&userReq))
	// 从gin.Key中取出服务实例
	userService := c.Keys["userService"].(services.UserService)
	// 调用userService前后 会执行wrapper
	userResp, err := userService.UserLogin(context.Background(), &userReq)
	PanicIfUserError(err)
	c.JSON(http.StatusOK, gin.H{"data": userResp})
}

// 包装错误
func PanicIfUserError(err error) {
	if err != nil {
		err = errors.New("userService--" + err.Error())
		userLog := log.Fields(map[string]interface{}{
			"service": "user",
		})
		userLog.Log(logger.InfoLevel, err)
		panic(err)
	}
}

func PanicIfTaskError(err error) {
	if err != nil {
		err = errors.New("taskService--" + err.Error())
		log.Log(logger.InfoLevel, err)
		panic(err)
	}
}
