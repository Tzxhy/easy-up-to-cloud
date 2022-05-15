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

type ResourceGroupDirItem struct {
	Gid        string         `json:"gid"`
	Rid        string         `json:"rid"`
	Fid        string         `json:"fid"`
	Did        string         `json:"did"`
	Name       string         `json:"name"`
	ParentDid  string         `json:"parent_did"`
	RType      uint8          `json:"r_type"`
	Uid        string         `json:"-"`
	CreateDate string         `json:"create_date"`
	ExpireDate sql.NullString `json:"expire_date"`
}

func GetResourceGroup(uid string) *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, create_date from user_group where group_type = ? or user_ids like ?", GroupTypeCommon, "%"+uid+"%")
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
	accounts := GetAdminAccount()
	if utils.Has(accounts, uid) { // 管理员
		uid = "%" // 清空查询
	}
	row := DB.QueryRow("select * from user_group where gid = ? and (user_ids like ? or group_type = ?)", gid, "%"+uid+"%", GroupTypeCommon)

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
		log.Print("GetGroupByIdAndUid err: ", err)
		return nil
	}
	item.UserIds = strings.Split(userIds, ";")
	return item
}
func GetGroupResourceById(rid string) *ResourceGroupDirItem {
	row := DB.QueryRow("select * from user_group_resource where rid = ?", rid)

	var item = &ResourceGroupDirItem{}
	err := row.Scan(
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
	if err != nil {
		log.Print("err: ", err)
		return nil
	}
	return item
}

func GetAllResourceGroup() *[]ResourceGroupItem {
	rows, err := DB.Query("select gid, name, user_ids, create_date, group_type from user_group")
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
			&item.GroupType,
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

func GetGroupDir(parent_did string) *[]ResourceGroupDirItem {
	rows, err := DB.Query("select * from user_group_resource where parent_did = ?", parent_did)
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
	stmt, _ := DB.Prepare("update user_group set user_ids = ? where gid = ? and group_type <> ?")
	defer stmt.Close()

	ret, err := stmt.Exec(strings.Join(*utils.Filter(&user_ids, func(s string) bool {
		return s != ""
	}), ";"), gid, GroupTypeCommon)
	if err != nil {
		return false, err
	}
	rows, _ := ret.RowsAffected()
	refreshLocalKeyCache()
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
	refreshLocalKeyCache()
	CreateGroupDir(gid, "", name, "")
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

func GetAdminAccount() *[]string {
	rows, _ := DB.Query("select * from admin")
	defer rows.Close()

	var accounts = make([]string, 1)
	for rows.Next() {
		a := ""
		rows.Scan(&a)
		accounts = append(accounts, a)
	}
	return &accounts
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
		AdminAccount = append(AdminAccount, uid)
	}

	row, _ := result.RowsAffected()
	return row == 1
}

func ShareFileToGroup(gid, fid, name, uid, parent_did, expire_date string, rtype uint8) (rid string, err error) {
	rid = utils.GenerateRid()
	stmt, err := DB.Prepare(
		`insert into user_group_resource (gid, rid, fid, name, parent_did, rtype, uid, expire_date)
		values (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	utils.CheckErr(err)
	result, err := stmt.Exec(gid, rid, fid, name, parent_did, rtype, uid, expire_date)
	utils.CheckErr(err)
	_, err = result.RowsAffected()
	if err != nil {
		return "", err
	}
	return rid, nil
}

func SearchResourceByName(name string) *[]ResourceGroupDirItem {
	rows, err := DB.Query("select * from user_group_resource where name like ?", "%"+name+"%")
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

func DeleteResourceByUidAndRid(uid, rid string) (bool, error) {
	stmt, _ := DB.Prepare("delete from user_group_resource where uid = ? and rid = ?")
	defer stmt.Close()

	ret, err := stmt.Exec(uid, rid)
	if err != nil {
		return false, err
	}
	row, _ := ret.RowsAffected()
	return row == 1, nil
}

func MoveResourceByUidAndRid(uid, rid, new_parent_did string) (bool, error) {
	stmt, _ := DB.Prepare("update user_group_resource set parent_did = ? where uid = ? and rid = ?")
	defer stmt.Close()

	ret, err := stmt.Exec(new_parent_did, uid, rid)
	if err != nil {
		return false, err
	}
	row, _ := ret.RowsAffected()
	return row == 1, nil
}

func RenameResourceByUidAndRid(uid, rid, newName string) (bool, error) {
	stmt, _ := DB.Prepare("update user_group_resource set name = ? where uid = ? and rid = ?")
	defer stmt.Close()

	ret, err := stmt.Exec(newName, uid, rid)
	if err != nil {
		return false, err
	}
	row, _ := ret.RowsAffected()
	return row == 1, nil
}
func ExpireChangeResourceByUidAndRid(uid, rid, newExpireDate string) (bool, error) {
	stmt, _ := DB.Prepare("update user_group_resource set expire_date = ? where uid = ? and rid = ?")
	defer stmt.Close()

	ret, err := stmt.Exec(newExpireDate, uid, rid)
	if err != nil {
		return false, err
	}
	row, _ := ret.RowsAffected()
	return row == 1, nil
}

// 直接通过fid删除文件分享记录，慎用
func DeleteResourceByFid(fid string) {
	DB.Exec("delete from user_group_resource where fid = ?", fid)
}

// 直接通过rid删除文件分享记录，慎用
func DeleteResourceByRid(rid string) {
	DB.Exec("delete from user_group_resource where rid = ?", rid)
}

func MoveMultiGroups(rids []string, new_parent_did string) uint8 {

	newRids := *utils.Map(&rids, func(st string) string {
		return "'" + st + "'"
	})
	sql := "update user_group_resource set parent_did = '" +
		new_parent_did +
		"' where rid in (" + strings.Join(
		newRids,
		", ",
	) + ")"
	ret, err := DB.Exec(sql)
	if err != nil {
		log.Print("MoveMultiGroups exec err: ", err)
		return 0
	}
	row, _ := ret.RowsAffected()
	return uint8(row)
}
