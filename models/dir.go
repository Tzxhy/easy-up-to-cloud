package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type Dir struct {
	Did     string
	OwnerId string
	Dirname string
	// -1 为根目录
	ParentDiD  string
	CreateDate string
}

func AddDir(owner_id string, dirname string, parent_did string) (string, error) {
	originDir := GetDirByName(dirname, owner_id, parent_did)
	if originDir != nil { // 已有同名
		return "", errors.New(constants.TIPS_HAS_SAME_DIR)
	}
	if parent_did != "" { // non root
		parentDir := GetDir(parent_did, owner_id)
		if parentDir == nil { // no parent dir
			return "", errors.New(constants.TIPS_CREATE_DIR_WITH_NO_EXIST_PARENT)
		}
	}
	stmt, err := DB.Prepare("insert into dirs (did, owner_id, dirname, parent_did) values(?, ?, ?, ?)")
	utils.CheckErr(err)
	did := utils.RandStringBytesMaskImprSrc(5)
	_, err = stmt.Exec(did, owner_id, dirname, parent_did)
	utils.CheckErr(err)
	return did, nil
}

func GetDirByName(name string, owner_id string, parent_did string) *Dir {
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
func GetDir(did string, owner_id string) *Dir {
	rows, err := DB.Query("select * from dirs where owner_id = ? and did = ?", owner_id, did)
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
