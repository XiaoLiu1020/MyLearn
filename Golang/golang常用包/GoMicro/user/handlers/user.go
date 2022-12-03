package handlers

import (
	"context"
	"errors"
	"time"
	services "user/services/proto" //  导入proto
)

/* 实现proto Handler的方法
type UserServiceHandler interface {
	UserLogin(context.Context, *UserRequest, *UserDetailResponse) error
	UserRegister(context.Context, *UserRequest, *UserDetailResponse) error
}
*/

type UserService struct{}

func (*UserService) UserRegister(ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	// 模板, 可以自己调入gorm使用
	resp.Code = 200
	now := time.Now().Unix()
	resp.UserDetail = &services.UserModel{
		ID:        uint32(1),
		UserName:  "XiaoLiu",
		CreatedAt: now,
		UpdatedAt: now,
	}
	return nil
}

func (*UserService) UserLogin(ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	/* resp.Code = 400
	now := time.Now().Unix()
	resp.UserDetail = &services.UserModel{
		ID:        uint32(1),
		UserName:  "XiaoLiu",
		CreatedAt: now,
		UpdatedAt: now,
	} */
	// 报错误过多可以触发熔断
	return errors.New("some problem in UserLogin service")
}
