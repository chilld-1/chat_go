package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID         uint      `gorm:"primary_key"`
	Username   string    `gorm:"unique;not null"`
	Password   string    `gorm:"not null;"`
	CreateTime time.Time `gorm:"default:current_timestamp"`
}

func (u *User) TableName() string {
	return "user"
}
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}).Error
}
