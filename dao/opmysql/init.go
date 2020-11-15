package opmysql

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	// 读取配置文件
	_ "github.com/qnfnypen/mumori/public"
	// 连接数据库
	_ "github.com/go-sql-driver/mysql"
)

var db *gorm.DB

func init() {
	var err error

	conn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`,
		viper.GetString("Driver.MySQL.User"),
		viper.GetString("Driver.MySQL.Password"),
		viper.GetString("Driver.MySQL.Host"),
		viper.GetString("Driver.MySQL.Database"))

	db, err = gorm.Open("mysql", conn)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("连接MySQL数据库失败")
	}

	// 进行MySQL的相关设置
	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(viper.GetInt("Driver.MySQL.MaxIdleConns"))
	db.DB().SetMaxOpenConns(viper.GetInt("Driver.MySQL.MaxOpenConns"))
}
