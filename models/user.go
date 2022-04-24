package models

import (
	"fmt"
	"log"

	"gitee.com/tzxhy/web/utils"
)

type User struct {
	Uid      int
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

func AddUser(username, password string) (int64, error) {
	stmt, err := DB.Prepare("insert into users (name, password) values(?, ?)")
	utils.CheckErr(err)

	result, err := stmt.Exec(username, password)
	utils.CheckErr(err)

	return result.LastInsertId()
}

func GetUserById(id int) *User {
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
	rows, err := DB.Query("select * from users where name = ? and password = ?", username, password)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var user *User = new(User)
	for rows.Next() {
		rows.Scan(&user.Uid, &user.Username, &user.Password)
	}
	fmt.Print("user: ", user)
	return user
}
