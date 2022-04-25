package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type Dir struct {
	Did      	int
	OwnerId  	int
	Dirname		string
	// -1 为根目录
	ParentDiD 	int
	CreateDate 	string
}

func AddDir(owner_id int, dirname string, parent_did int) (int64, error) {
	originDir := GetDirByName(dirname, owner_id, parent_did)
	if originDir != nil { // 已有同名
		return 0, errors.New(constants.TIPS_HAS_SAME_DIR)
	}
	stmt, err := DB.Prepare("insert into dirs (owner_id, dirname, parent_did) values(?, ?, ?)")
	utils.CheckErr(err)

	result, err := stmt.Exec(owner_id, dirname, parent_did)
	utils.CheckErr(err)

	return result.LastInsertId()
}
func GetDirByName(name string, owner_id int, parent_did int) *Dir {
	rows, err := DB.Query("select * from dirs where owner_id = ? and dirname = ? and parent_did = ?", owner_id, name, parent_did)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dir *Dir
	for rows.Next() {
		dir = new(Dir)
		rows.Scan(&dir.Did, &dir.OwnerId, &dir.Dirname, &dir.ParentDiD, &dir.CreateDate)
		break
	}
	return dir
}
func GetDir(id int, owner_id int) *Dir {
	rows, err := DB.Query("select * from dirs where owner_id = ? and did = ?", owner_id, id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dir *Dir = new(Dir)
	for rows.Next() {
		rows.Scan(&dir.Did, &dir.OwnerId, &dir.Dirname, &dir.ParentDiD, &dir.CreateDate)
		break
	}
	return dir
}

// func GetUserByNameAndPassword(username, password string) *User {
// 	rows, err := DB.Query("select * from users where name = ? and password = ?", username, password)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer rows.Close()

// 	var user *User = new(User)
// 	for rows.Next() {
// 		rows.Scan(&user.Uid, &user.Username, &user.Password)
// 	}
// 	fmt.Print("user: ", user)
// 	return user
// }
