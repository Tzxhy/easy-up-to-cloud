package models

import (
	"database/sql"
	"log"
	"strings"

	"gitee.com/tzxhy/web/utils"
	"gorm.io/gorm"
)

type Admin struct {
	Uid string `gorm:"primaryKey;type:string not null;"`
}

const (
	// 与 下面 GroupType tag中保持一致
	GroupTypeCommon       = 0
	GroupTypeVisibleByUid = 1
)

const (
	RTYPE_DIR  = 1
	RTYPE_FILE = 2
)

type ResourceGroupType uint
type RType uint
type ResourceGroupItem struct {
	Gid        string            `json:"gid" gorm:"primaryKey;type:string not null;"`
	Name       string            `json:"name" gorm:"index;type:string not null;"`
	UserIds    string            `json:"-" gorm:"type:string not null;default:''"`
	CreateDate string            `json:"create_date" gorm:"type:datetime not null default CURRENT_TIMESTAMP;"`
	GroupType  ResourceGroupType `json:"-" gorm:"type:integer not null;default:0;check: group_type >= 0;"`
}

type ResourceGroupDirItem struct {
	Gid        string         `json:"gid" gorm:"index:resource_unique;type:string not null;"`
	Rid        string         `json:"rid" gorm:"primaryKey;type:string not null;"`
	Fid        string         `json:"fid" gorm:"type:string;default:''"`
	Did        string         `json:"did" gorm:"type:string;default:''"`
	Name       string         `json:"name" gorm:"index:resource_unique;type:string not null;"`
	ParentDid  string         `json:"parent_did" gorm:"index:resource_unique;type:string not null;default:'ROOT'"`
	RType      RType          `json:"r_type" gorm:"type:integer not null;"`
	Uid        string         `json:"-" gorm:"type:string not null;"`
	CreateDate string         `json:"create_date" gorm:"type:datetime not null default CURRENT_TIMESTAMP;"`
	ExpireDate sql.NullString `json:"expire_date" gorm:"type:integer default 0;"`
}

func GetResourceGroup(uid string) *[]ResourceGroupItem {
	var groups []ResourceGroupItem
	result := DB.Where(
		"group_type = ? or user_ids like ?",
		GroupTypeCommon,
		"%"+uid+"%",
	).Find(&groups)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &groups
}

func GetGroupById(gid string) *ResourceGroupItem {
	var group ResourceGroupItem
	result := DB.Where("gid = ?", gid).Take(&group)
	err := result.Error
	if err != nil {
		log.Print("err: ", err)
		return nil
	}
	// group.UserIds = strings.Split(userIds, ";")
	return &group
}

func GetGroupByIdAndUid(gid, uid string) *ResourceGroupItem {
	accounts := GetAdminAccount()
	if utils.Has(accounts, uid) { // 管理员
		uid = "%" // 清空查询
	}
	var item ResourceGroupItem
	result := DB.Where(
		"gid = ? and (user_ids like ? or group_type = ?)",
		gid,
		"%"+uid+"%",
		GroupTypeCommon,
	).Take(&item)

	err := result.Error
	if err != nil {
		log.Print("GetGroupByIdAndUid err: ", err)
		return nil
	}
	// item.UserIds = strings.Split(userIds, ";")
	return &item
}
func GetGroupResourceById(rid string) *ResourceGroupDirItem {
	var item ResourceGroupDirItem
	result := DB.Where("rid = ?", rid).Take(&item)
	err := result.Error
	if err != nil {
		log.Print("err: ", err)
		return nil
	}
	return &item
}

func GetAllResourceGroup() *[]ResourceGroupItem {
	var items []ResourceGroupItem
	result := DB.Find(&items)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &items
}

const (
	GROUP_RESOURCE_DIR  = 1
	GROUP_RESOURCE_FILE = 2
)

func GetGroupDir(parent_did string) *[]ResourceGroupDirItem {
	var items []ResourceGroupDirItem
	result := DB.Where("parent_did = ?", parent_did).Find(&items)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &items
}

func SetUidResourceGroup(gid string, user_ids []string) (succ bool, err error) {

	result := DB.Model(&ResourceGroupItem{}).Where(
		"gid = ? and group_type <> ?",
		gid,
		GroupTypeCommon,
	).Update("user_ids", strings.Join(*utils.Filter(&user_ids, func(s string) bool {
		return s != ""
	}), ";"))

	err = result.Error

	if err != nil {
		return false, err
	}
	rows := result.RowsAffected
	refreshLocalKeyCache()
	return rows == 1, nil
}
func CreateGroup(name string, groupType ResourceGroupType) (gid string, err error) {
	gid = utils.GenerateGid()
	result := DB.Select("Gid", "Name", "GroupType").Create(&ResourceGroupItem{
		Gid:       gid,
		Name:      name,
		GroupType: groupType,
	})
	sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Select("Gid", "Name", "GroupType").Create(&ResourceGroupItem{
			Gid:       gid,
			Name:      name,
			GroupType: groupType,
		})
	})
	log.Println("sql: ", sql)
	err = result.Error

	if err != nil {
		return "", err
	}
	refreshLocalKeyCache()
	CreateGroupDir(gid, "", name, "")
	return gid, nil
}

func CreateGroupDir(gid, parent_did, name, uid string) (rid string, err error) {
	rid = utils.GenerateRid()
	result := DB.Select("Gid", "Rid", "ParentDid", "Name", "Uid", "RType").Create(&ResourceGroupDirItem{
		Gid:       gid,
		Rid:       rid,
		ParentDid: parent_did,
		Name:      name,
		Uid:       uid,
		RType:     GROUP_RESOURCE_DIR,
	})
	err = result.Error

	if err != nil {
		return "", err
	}
	affected := result.RowsAffected
	if affected == 1 {
		return rid, nil
	}
	return "", err
}

func GetAdminAccount() *[]string {
	var admins []Admin
	DB.Find(&admins)

	return utils.Map(&admins, func(item Admin) string {
		return item.Uid
	})
}

//
func DeleteOrInsertAdminAccount(uid string, isAdmin bool) bool {
	var result sql.Result
	if !isAdmin {
		DB.Exec("delete from admin where uid = ? ", uid)
	} else {
		var err error
		result := DB.Exec("insert into admin values(?) ", uid)
		err = result.Error
		if err != nil {
			return false
		}
		AdminAccount = append(AdminAccount, Admin{
			Uid: uid,
		})
	}

	row, _ := result.RowsAffected()
	return row == 1
}

func ShareFileToGroup(gid, fid, name, uid, parent_did, expire_date string, rtype uint8) (rid string, err error) {
	rid = utils.GenerateRid()
	ret := DB.Exec(
		`insert into user_group_resource (gid, rid, fid, name, parent_did, rtype, uid, expire_date)
		values (?, ?, ?, ?, ?, ?, ?, ?)
	`, gid, rid, fid, name, parent_did, rtype, uid, expire_date)
	err = ret.Error

	if err != nil {
		return "", err
	}
	return rid, nil
}

func SearchResourceByName(name string) *[]ResourceGroupDirItem {
	var items []ResourceGroupDirItem
	ret := DB.Where(
		"name like ?",
		"%"+name+"%",
	).Find(&items)
	err := ret.Error
	if err != nil {
		log.Fatal(err)
	}

	return &items
}

func DeleteResourceByUidAndRid(uid, rid string) (bool, error) {
	ret := DB.Where("uid = ? and rid = ?", uid, rid).Delete(&ResourceGroupDirItem{})
	err := ret.Error
	if err != nil {
		return false, err
	}
	row := ret.RowsAffected
	return row == 1, nil
}

func MoveResourceByUidAndRid(uid, rid, new_parent_did string) (bool, error) {
	ret := DB.Model(&ResourceGroupDirItem{}).Where(
		"uid = ? and rid = ?",
		uid,
		rid,
	).Update("parent_did", new_parent_did)
	err := ret.Error
	if err != nil {
		return false, err
	}
	row := ret.RowsAffected
	return row == 1, nil
}

func RenameResourceByUidAndRid(uid, rid, newName string) (bool, error) {
	ret := DB.Model(&ResourceGroupDirItem{}).Where(
		"uid = ? and rid = ?",
		uid,
		rid,
	).Update("name", newName)
	err := ret.Error

	if err != nil {
		return false, err
	}
	row := ret.RowsAffected
	return row == 1, nil
}
func ExpireChangeResourceByUidAndRid(uid, rid, newExpireDate string) (bool, error) {
	result := DB.Model(&ResourceGroupDirItem{}).Where(
		"uid = ? and rid = ?",
		uid,
		rid,
	).Update("expire_date", newExpireDate)

	err := result.Error

	if err != nil {
		return false, err
	}
	row := result.RowsAffected
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
	result := DB.Raw(sql)
	err := result.Error
	if err != nil {
		log.Print("MoveMultiGroups exec err: ", err)
		return 0
	}
	row := result.RowsAffected
	return uint8(row)
}
