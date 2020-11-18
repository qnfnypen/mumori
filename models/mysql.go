package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qnfnypen/mumori/utils"
)

// User 用户信息
type User struct {
	ID        int        `gorm:"column:id"`
	UserName  string     `gorm:"column:user_name"`
	Phone     string     `gorm:"column:phone"`
	Password  string     `gorm:"column:password"`
	Status    int8       `gorm:"column:status"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	DeletedAt *time.Time `gorm:"column:deleated_at"`
}

// BeforeSave 在保存前回调，对密码进行加密
func (u *User) BeforeSave(scope *gorm.Scope) error {
	if u.Password != "" {
		hp, err := utils.EncryptPassword(u.Password)
		// 密码加密失败，则明文存储
		if err != nil {
			return nil
		}
		if scope.SetColumn("password", hp); err != nil {
			return err
		}
	}

	return nil
}
