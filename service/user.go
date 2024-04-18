package service

import (
	"errors"
	"jassue-gin/entity"
	"jassue-gin/entity/request"
	"jassue-gin/repository"
	"jassue-gin/util"
)

type UserService interface {
	Register(params request.Register) (error, entity.User)
	Login(params request.Login) (error, entity.User)
	GetUserInfo(id string) (error, entity.User)
}

type userServiceImpl struct {
	repo repository.UserRepository
}

// NewUserService 创建一个新的UserService实现的实例
func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (userServiceImpl *userServiceImpl) Register(params request.Register) (error, entity.User) {
	// 使用repository检查用户是否存在
	_, err := userServiceImpl.repo.FindByMobile(params.Mobile)
	if err == nil { // 如果没有错误，说明找到了用户，手机号已被注册
		return errors.New("手机号已存在"), entity.User{} // 明确返回错误和空的User
	}

	secretKey, err := util.GenerateSecureKey() // 生成的Hex字符串长度会是64
	if err != nil {
		return errors.New("Error generating secure secret key"), entity.User{} // 明确返回错误和空的User
	}

	user := entity.User{Name: params.Name, Mobile: params.Mobile, Password: util.BcryptMake([]byte(params.Password)), AccessKey: util.GenerateAccessKey(), SecretKey: secretKey, UserRole: "user"}
	// 创建用户
	if err := userServiceImpl.repo.Create(&user); err != nil {
		return err, entity.User{} // 创建用户时出错，返回错误和User的零值
	}
	return nil, user // 成功创建用户，返回nil错误和用户信息
}

func (userServiceImpl *userServiceImpl) Login(params request.Login) (error, entity.User) {
	user, err := userServiceImpl.repo.FindByMobile(params.Mobile)
	if err != nil {
		return errors.New("用户不存在"), entity.User{}
	}
	if !util.BcryptMakeCheck([]byte(params.Password), user.Password) {
		return errors.New("密码错误"), entity.User{}
	}
	return nil, entity.User{}
}

func (userServiceImpl *userServiceImpl) GetUserInfo(id string) (error, entity.User) {
	user, err := userServiceImpl.repo.FindById(id)
	if err != nil {
		return errors.New("用户不存在"), entity.User{}
	}
	return nil, *user
}
