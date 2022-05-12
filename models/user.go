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
	uid := utils.GenerateUid()
	_, err = stmt.Exec(uid, username, password)
	utils.CheckErr(err)

	return uid, nil
}

func AddUserWithId(uid, username, password string) (string, error) {
	stmt, err := DB.Prepare("insert into users (uid, name, password) values(?, ?, ?)")
	utils.CheckErr(err)
	_, err = stmt.Exec(uid, username, password)
	utils.CheckErr(err)

	return uid, nil
}

func GetUserById(id string) *User {
	row := DB.QueryRow("select uid, name from users where uid = ?", id)

	var user *User = new(User)
	err := row.Scan(&user.Uid, &user.Username)
	if err != nil {
		return nil
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

func GetAdminUser() *[]User {
	rows, err := DB.Query("select uid from admin")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var admins []string
	for rows.Next() {
		admin := ""
		rows.Scan(&admin)
		admins = append(admins, admin)
	}
	log.Print("admins: ", admins)
	var users []User
	for _, uid := range admins {
		user := GetUserById(uid)
		if user != nil {
			users = append(users, *user)
		}
	}
	log.Print("users: ", users)
	return &users
}
