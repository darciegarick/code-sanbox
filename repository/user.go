package repository

import (
	"fmt"
	"jassue-gin/entity"
	"jassue-gin/global"
)

type UserRepository interface {
	FindById(id string) (*entity.User, error)
	FindByMobile(mobile string) (*entity.User, error)
	Create(user *entity.User) error
}

type userRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &userRepositoryImpl{}
}

func (repo *userRepositoryImpl) FindByMobile(mobile string) (*entity.User, error) {
	var user entity.User
	result := global.App.DB.Where("mobile = ?", mobile).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *userRepositoryImpl) Create(user *entity.User) error {
	return global.App.DB.Create(user).Error
}

func (repo *userRepositoryImpl) FindById(id string) (*entity.User, error) {
	var user entity.User
	result := global.App.DB.Where("id = ?", id).First(&user)
	fmt.Println(result)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
