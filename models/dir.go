package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type Dir struct {
	Did     string `json:"did" gorm:"primaryKey"`
	User    User   `json:"-" gorm:"references:Uid"`
	UserID  string `json:"-" gorm:"uniqueIndex:dir_unique;type:string not null"`
	Dirname string `json:"dirname" gorm:"uniqueIndex:dir_unique;type:string not null"`
	// -1 为根目录
	ParentDid  string `json:"parent_did" gorm:"uniqueIndex:dir_unique;type:string not null;default:'ROOT'"`
	CreateDate int64  `json:"create_date" gorm:"autoUpdateTime:milli"`
}

func AddDir(owner_id string, dirname string, parent_did string) (string, error) {
	if parent_did != constants.DIR_ROOT_ID { // non root
		parentDir := GetDir(parent_did, owner_id)
		if parentDir == nil { // no parent dir
			return "", errors.New(constants.TIPS_CREATE_DIR_WITH_NO_EXIST_PARENT)
		}
	}

	did := utils.GenerateDid()
	dirItem := &Dir{
		Did:       did,
		UserID:    owner_id,
		Dirname:   dirname,
		ParentDid: parent_did,
	}

	result := DB.Create(&dirItem)
	if result.Error != nil {
		log.Print("error: ", result.Error)
		return "", result.Error
	}
	return did, nil
}

func GetDirByName(name string, owner_id string, parent_did string) *Dir {
	var dir Dir
	result := DB.Where(&Dir{
		UserID:    owner_id,
		ParentDid: parent_did,
		Dirname:   name,
	}).Limit(1).Find(&dir)
	err := result.Error
	if err != nil {
		return nil
	}
	if result.RowsAffected < 1 {
		return nil
	}
	return &dir
}
func GetDir(did string, owner_id string) *Dir {
	var dir Dir
	result := DB.Where(&Dir{
		UserID: owner_id,
		Did:    did,
	}).Limit(1).Find(&dir)
	err := result.Error
	if err != nil {
		log.Print("err: ", err)
		return nil
	}
	if result.RowsAffected < 1 {
		return nil
	}
	return &dir
}

func GetDirList(parent_id, owner_id string) *[]Dir {
	var dirs []Dir
	result := DB.Where(&Dir{
		UserID:    owner_id,
		ParentDid: parent_id,
	}).Find(&dirs)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &dirs
}

func SearchDirList(owner_id, dirname string) *[]Dir {
	var dirs []Dir
	result := DB.Where(
		"user_id = ? and dirname like ?",
		owner_id,
		"%"+dirname+"%",
	).Find(&dirs)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &dirs
}

func MoveDir(owner_id string, did string, new_parent_did string) bool {
	if new_parent_did != "" && GetDir(new_parent_did, owner_id) == nil {
		return false
	}
	result := DB.Where(&Dir{
		UserID: owner_id,
		Did:    did,
	}).Updates(&Dir{
		ParentDid: new_parent_did,
	})

	err := result.Error
	utils.CheckErr(err)

	lines := result.RowsAffected
	return lines == 1
}

func RenameDir(owner_id string, did string, new_name string) bool {
	result := DB.Where(&Dir{
		UserID: owner_id,
		Did:    did,
	}).Updates(&Dir{
		Dirname: new_name,
	})
	err := result.Error

	if err != nil {
		return false
	}
	lines := result.RowsAffected
	return lines == 1
}

// 删除指定did文件夹（不会递归删除其下文件/文件夹）
func DeleteSingleDir(owner_id string, did string) bool {
	result := DB.Delete(&Dir{}, &Dir{
		UserID: owner_id,
		Did:    did,
	})
	err := result.Error
	if err != nil {
		return false
	}
	lines := result.RowsAffected
	return lines == 1
}
