package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type Dir struct {
	Did     string `json:"did" gorm:"primaryKey"`
	OwnerId string `gorm:"index:dir_unique;type:string not null"`
	Dirname string `json:"dirname" gorm:"index:dir_unique;type:string not null"`
	// -1 为根目录
	ParentDid  string `json:"parent_did" gorm:"index:dir_unique;type:string not null;default:''"`
	CreateDate string `json:"create_date" gorm:"type:datetime;default:CURRENT_TIMESTAMP"`
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

	did := utils.GenerateDid()
	dirItem := &Dir{
		Did:       did,
		OwnerId:   owner_id,
		Dirname:   dirname,
		ParentDid: parent_did,
	}

	result := DB.Create(&dirItem)
	utils.CheckErr(result.Error)
	return did, nil
}

func GetDirByName(name string, owner_id string, parent_did string) *Dir {
	var dir Dir
	result := DB.Where(&Dir{
		OwnerId:   owner_id,
		ParentDid: parent_did,
		Dirname:   name,
	}).Take(&dir)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}

	return &dir
}
func GetDir(did string, owner_id string) *Dir {
	var dir Dir
	result := DB.Where(&Dir{
		OwnerId: owner_id,
		Did:     did,
	}).Take(&dir)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &dir
}

func GetDirList(parent_id, owner_id string) *[]Dir {
	var dirs []Dir
	result := DB.Where(&Dir{
		OwnerId:   owner_id,
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
		"owner_id = ? and dirname like ?",
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
	result := DB.Model(&Dir{
		OwnerId: owner_id,
		Did:     did,
	}).Updates(&Dir{
		ParentDid: new_parent_did,
	})
	err := result.Error
	utils.CheckErr(err)

	lines := result.RowsAffected
	return lines == 1
}

func RenameDir(owner_id string, did string, new_name string) bool {
	result := DB.Model(&Dir{
		OwnerId: owner_id,
		Did:     did,
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
	result := DB.Delete(&Dir{
		OwnerId: owner_id,
		Did:     did,
	})
	err := result.Error
	if err != nil {
		return false
	}
	lines := result.RowsAffected
	return lines == 1
}
