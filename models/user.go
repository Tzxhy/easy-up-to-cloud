package models

import (
	"fmt"
	"log"

	"gitee.com/tzxhy/web/utils"
)

type User struct {
	Uid      string
	Username string
	Password string
}

func IsUserOnline(username, password string) bool {
	return GetKey(fmt.Sprintf("%s-%s", username, password)).(bool)
}

func HasUsername(username string) bool {
	rows, err := DB.Query("select * from users where name = ?", username)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	hasRow := false
	for rows.Next() {
		hasRow = true
		break
	}

	return hasRow
}

func AddUser(username, password string) (string, error) {
	stmt, err := DB.Prepare("insert into users (uid, name, password) values(?, ?, ?)")
	utils.CheckErr(err)
	uid := utils.RandStringBytesMaskImprSrc(5)
	_, err = stmt.Exec(uid, username, password)
	utils.CheckErr(err)

	return uid, nil
}

func GetUserById(id string) *User {
	rows, err := DB.Query("select * from users where uid = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user *User = new(User)
	for rows.Next() {
		rows.Scan(&user.Uid, &user.Username, &user.Password)
	}
	return user
}

func GetUserByNameAndPassword(username, password string) *User {
	rows, err := DB.Query("select uid, name, password from users where name = ? and password = ?", username, password)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user *User
	for rows.Next() {
		user = new(User)
		rows.Scan(&user.Uid, &user.Username, &user.Password)
		break
	}
	return user
}
