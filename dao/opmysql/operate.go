package opmysql

import (
	"errors"

	"github.com/qnfnypen/mumori/models"
	"github.com/qnfnypen/mumori/utils"
)

// StoreUserInfo 存储用户信息
func StoreUserInfo(user models.User) error {
	return db.Create(&user).Error
}

// CheckUserInfo 用户信息匹配与检测
func CheckUserInfo(account, password, category string) error {
	var user models.User
	switch category {
	case "username":
		db.Where("user_name = ?",account).First(&user)
	case "phone":
		db.Where("phone = ?",account).First(&user)
	case "email":
		db.Where("email = ?",account).First(&user)
	}

	if user.UID < 2020 {
		return errors.New("该用户不存在")
	}
	
	return utils.CompareHashAndPassword(password,user.Password)
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(user models.User) error {
	return nil
}
