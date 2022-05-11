package models

import (
	"database/sql"
	"log"
	"strings"

	"gitee.com/tzxhy/web/utils"
)

const (
	GroupTypeCommon       = 0
	GroupTypeVisibleByUid = 1
)

type ResourceGroupItem struct {
	Gid        string   `json:"gid"`
	Name       string   `json:"name"`
	UserIds    []string `json:"-"`
	CreateDate string   `json:"create_date"`
	GroupType  uint8    `json:"-"`
}

func GetResourceGroup(uid string) *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, create_date from user_group where group_type = 0 or user_ids like ?", "%"+uid+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items = make([]ResourceGroupItem, 0)
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

func GetGroupById(gid string) *ResourceGroupItem {
	row := DB.QueryRow("select * from user_group where gid = ?", gid)

	var item = &ResourceGroupItem{}
	userIds := ""
	err := row.Scan(
		&item.Gid,
		&item.Name,
		&userIds,
		&item.CreateDate,
		&item.GroupType,
	)
	if err != nil {
		log.Print("err: ", err)
		return nil
	}
	item.UserIds = strings.Split(userIds, ";")
	return item
}

func GetGroupByIdAndUid(gid, uid string) *ResourceGroupItem {
	row := DB.QueryRow("select * from user_group where gid = ? and user_ids like ?", gid, "%"+uid+"%")

	var item = &ResourceGroupItem{}
	userIds := ""
	err := row.Scan(
		&item.Gid,
		&item.Name,
		&userIds,
		&item.CreateDate,
	)
	if err != nil {
		return nil
	}
	item.UserIds = strings.Split(userIds, ";")
	return item
}

func GetAllResourceGroup() *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, user_ids, create_date from user_group")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items = make([]ResourceGroupItem, 0)
	for rows.Next() {
		item := &ResourceGroupItem{}
		userIds := ""
		rows.Scan(
			&item.Gid,
			&item.Name,
			&userIds,
			&item.CreateDate,
		)
		item.UserIds = strings.Split(userIds, ";")
		items = append(items, *item)
	}
	return &items
}

const (
	GROUP_RESOURCE_DIR  = 1
	GROUP_RESOURCE_FILE = 2
)

type ResourceGroupDirItem struct {
	Gid        string `json:"gid"`
	Rid        string `json:"rid"`
	Fid        string `json:"fid"`
	Did        string `json:"did"`
	Name       string `json:"name"`
	ParentDid  string `json:"parent_did"`
	RType      uint8  `json:"r_type"`
	Uid        string `json:"-"`
	CreateDate string `json:"create_date"`
	ExpireDate string `json:"expire_date"`
}

func GetGroupDir(gid, parent_did string) *[]ResourceGroupDirItem {
	rows, err := DB.Query("select * from user_group_resource where gid = ? and parent_did = ?", gid, parent_did)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var items = make([]ResourceGroupDirItem, 0)
	for rows.Next() {
		item := new(ResourceGroupDirItem)
		rows.Scan(
			&item.Gid,
			&item.Rid,
			&item.Fid,
			&item.Did,
			&item.Name,
			&item.ParentDid,
			&item.RType,
			&item.Uid,
			&item.CreateDate,
			&item.ExpireDate,
		)
		items = append(items, *item)
	}
	return &items
}

func SetUidResourceGroup(gid string, user_ids []string) (succ bool, err error) {
	stmt, _ := DB.Prepare("update user_group set user_ids = ? where gid = ?")
	defer stmt.Close()

	ret, err := stmt.Exec(strings.Join(user_ids, ";"), gid)
	if err != nil {
		return false, err
	}
	rows, _ := ret.RowsAffected()
	return rows == 1, nil
}
func CreateGroup(name string, groupType uint8) (gid string, err error) {
	gid = utils.GenerateGid()
	stmt, _ := DB.Prepare("insert into user_group (gid, name, group_type) values (?, ?, ?)")
	defer stmt.Close()
	_, err = stmt.Exec(gid, name, groupType)
	if err != nil {
		return "", err
	}
	return gid, nil
}

func CreateGroupDir(gid, parent_did, name, uid string) (rid string, err error) {
	stmt, _ := DB.Prepare("insert into user_group_resource (gid, rid, parent_did, name, uid, rtype) values (?, ?, ?, ?, ?, ?)")
	defer stmt.Close()
	rid = utils.GenerateRid()
	ret, err := stmt.Exec(gid, rid, parent_did, name, uid, GROUP_RESOURCE_DIR)
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

func ShareDirOrFileToGroup(gid, fid, did, name, uid, parent_did, expire_date string, rtype uint8) (rid string, err error) {
	rid = utils.GenerateRid()
	stmt, err := DB.Prepare(
		`insert into user_group_resource (gid, rid, fid, did, name, parent_did, rtype, uid, expire_date)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	utils.CheckErr(err)
	result, err := stmt.Exec(gid, rid, fid, did, name, parent_did, rtype, uid, expire_date)
	utils.CheckErr(err)
	_, err = result.RowsAffected()
	if err != nil {
		return "", err
	}
	return rid, nil
}
