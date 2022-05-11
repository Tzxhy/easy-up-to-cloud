package models

import (
	"database/sql"
	"log"

	"gitee.com/tzxhy/web/utils"
)

type ResourceGroupItem struct {
	Gid        string `json:"gid"`
	Name       string `json:"name"`
	CreateDate string `json:"create_date"`
}

func GetResourceGroup(uid string) *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, create_date from user_group where user_ids like ?", "%"+uid+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []ResourceGroupItem
	for rows.Next() {
		item := new(ResourceGroupItem)
		rows.Scan(
			&item.Gid,
			&item.Name,
			&item.CreateDate,
		)
		items = append(items, *item)
	}
	return &items
}

func GetAllResourceGroup() *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, create_date from user_group")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items []ResourceGroupItem
	for rows.Next() {
		item := new(ResourceGroupItem)
		rows.Scan(
			&item.Gid,
			&item.Name,
			&item.CreateDate,
		)
		items = append(items, *item)
	}
	return &items
}

func CreateGroupDir(gid, parent_did, name, uid string) (rid string, err error) {
	stmt, _ := DB.Prepare("insert into user_group_resource (gid, rid, parent_did, name, uid) values (?, ?, ?, ?, ?)")
	defer stmt.Close()
	rid = utils.GenerateRid()
	ret, err := stmt.Exec(gid, rid, parent_did, name, uid)
	if err != nil {
		return "", err
	}
	affected, err := ret.RowsAffected()
	if affected == 1 {
		return rid, nil
	}
	return "", err
}

//
func DeleteOrInsertAdminAccount(uid string, isAdmin bool) bool {
	var result sql.Result
	if !isAdmin {
		result, _ = DB.Exec("delete from admin where uid = ? ", uid)
	} else {
		var err error
		result, err = DB.Exec("insert into admin values(?) ", uid)
		if err != nil {
			return false
		}
	}

	row, _ := result.RowsAffected()
	return row == 1
}
