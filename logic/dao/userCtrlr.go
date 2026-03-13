package dao

import (
	"gochat/db"
	"gochat/model"
)

func GetUserByID(ID int) (*model.User, error) {
	var user model.User
	err := db.DB.Where("id = ?", ID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}
func UpdateUser(user *model.User) error {
	return db.DB.Save(user).Error
}
func DeleteUser(user *model.User) error {
	return db.DB.Delete(user).Error
}
