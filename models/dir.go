package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type Dir struct {
	Did     string `json:"did"`
	OwnerId string `json:"owner_id"`
	Dirname string `json:"dirname"`
	// -1 为根目录
	ParentDiD  string `json:"parent_did"`
	CreateDate string `json:"create_date"`
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

func GetDirList(parent_id, owner_id string) *[]Dir {
	rows, err := DB.Query("select * from dirs where owner_id = ? and parent_did = ?", owner_id, parent_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var dirs []Dir
	for rows.Next() {
		dir := new(Dir)
		rows.Scan(&dir.Did, &dir.OwnerId, &dir.Dirname, &dir.ParentDiD, &dir.CreateDate)
		dirs = append(dirs, *dir)
	}
	return &dirs
}

func MoveDir(owner_id string, did string, new_parent_did string) bool {
	if new_parent_did != "" && GetDir(new_parent_did, owner_id) == nil {
		return false
	}
	stmt, err := DB.Prepare("update dirs set parent_did = ? where owner_id = ? and did = ?")
	utils.CheckErr(err)
	defer stmt.Close()
	result, err := stmt.Exec(new_parent_did, owner_id, did)
	utils.CheckErr(err)
	lines, _ := result.RowsAffected()
	return lines == 1
}

func RenameDir(owner_id string, did string, new_name string) bool {
	stmt, err := DB.Prepare("update dirs set dirname = ? where owner_id = ? and did = ?")
	utils.CheckErr(err)
	defer stmt.Close()
	result, err := stmt.Exec(new_name, owner_id, did)
	utils.CheckErr(err)
	lines, _ := result.RowsAffected()
	return lines == 1
}
