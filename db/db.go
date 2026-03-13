package db

import (
	"gochat/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func Init() error {
	var err error
	DB, err = gorm.Open("mysql", viper.GetString("common-mysql.dsn"))
	if err != nil {
		return err
	}
	err = autoMigrate()
	if err != nil {
		return err
	}
	return nil
}
func autoMigrate() error {
	// 自动迁移数据库表
	err := model.AutoMigrate(DB)
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
