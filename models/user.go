package models

import (
	"fmt"
	"log"

	"gitee.com/tzxhy/web/utils"
)

type User struct {
	Uid        string `gorm:"primaryKey"`
	Username   string `gorm:"index;type:string not null"`
	Password   string `gorm:"type:string not null"`
	CreateDate string `gorm:"type:datetime not null default CURRENT_TIMESTAMP;"`
}

func IsUserOnline(username, password string) bool {
	return GetKey(fmt.Sprintf("%s-%s", username, password)).(bool)
}

func HasUsername(username string) bool {
	var count int64
	DB.Where(&User{
		Username: username,
	}).Limit(1).Find(&User{}).Count(&count)

	return count > 0
}

func AddUser(username, password string) (string, error) {
	uid := utils.GenerateUid()
	result := DB.Create(&User{
		Uid:      uid,
		Username: username,
		Password: password,
	})

	err := result.Error
	if err != nil {
		return "", err
	}

	return uid, nil
}

func AddUserWithId(uid, username, password string) (string, error) {
	result := DB.Create(&User{
		Uid:      uid,
		Username: username,
		Password: password,
	})
	err := result.Error

	utils.CheckErr(err)

	return uid, nil
}

func GetUserById(id string) *User {
	var user User
	result := DB.Find(&user, User{
		Uid: id,
	})

	err := result.Error

	if err != nil {
		return nil
	}
	return &user
}

func GetUserByNameAndPassword(username, password string) *User {
	var user User
	result := DB.Where(&User{
		Username: username,
		Password: password,
	}).Limit(1).Find(&user)

	err := result.Error

	if err != nil {
		return nil
	}
	if result.RowsAffected < 1 {
		return nil
	}
	return &user
}

func GetAdminUser() *[]User {
	var admins []Admin
	result := DB.Find(&admins)

	err := result.Error

	if err != nil {
		log.Fatal(err)
	}

	var users []User
	for _, admin := range admins {
		user := GetUserById(admin.Uid)
		if user != nil {
			users = append(users, *user)
		}
	}
	return &users
}

func GetUserByIds(ids []string) *[]User {
	var users []User
	err := DB.Find(&users, ids).Error
	if err != nil {
		log.Print("GetUserByIds err: ", err)
	}
	return &users
}
