package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type File struct {
	Fid      string `json:"fid" gorm:"primarykey;type:string not null;"`
	OwnerId  string `json:"-" gorm:"type:string not null;"`
	Filename string `json:"filename" gorm:"type:string not null;"`
	Filesize uint64 `json:"file_size" gorm:"type:integer not null;"`
	// -1 为根目录
	ParentDiD    string `json:"parent_did" gorm:"type:string not null;default:''"`
	FileRealPath string `json:"-" gorm:"type:string not null;"`
	CreateDate   string `json:"create_date" gorm:"type:DATETIME not null;default:CURRENT_TIMESTAMP"`
}

//
func AddFile(owner_id string, dir_id string, filename string, file_size uint64, file_path string) (string, error) {

	originDir := GetFileByName(filename, owner_id, dir_id)
	if originDir != nil { // 已有同名
		return "", errors.New(constants.TIPS_HAS_SAME_FILE)
	}
	fid := utils.GenerateFid()
	result := DB.Create(&File{
		Fid:          fid,
		OwnerId:      owner_id,
		Filename:     filename,
		ParentDiD:    dir_id,
		FileRealPath: file_path,
		Filesize:     file_size,
	})
	err := result.Error

	if err != nil {
		log.Print("err: ", err)
	}

	return fid, nil
}

func GetFileByName(filename string, owner_id string, parent_did string) *File {
	var file File
	result := DB.Where(&File{
		OwnerId:   owner_id,
		Filename:  filename,
		ParentDiD: parent_did,
	}).Take(&file)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &file
}
func GetFile(fid string, owner_id string) *File {
	var file File
	result := DB.Where(&File{
		OwnerId: owner_id,
		Fid:     fid,
	}).Take(&file)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}

	return &file
}

func GetFileList(parent_id, owner_id string) *[]File {
	var files []File
	result := DB.Where(&File{
		OwnerId:   owner_id,
		ParentDiD: parent_id,
	}).Find(&files)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &files
}

func SearchFileList(owner_id, name string) *[]File {
	var files []File
	result := DB.Where(
		"owner_id = ? and filename like ?",
		owner_id,
		"%"+name+"%",
	).Find(&files)
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	return &files
}

func DeleteFile(fid string, owner_id string) bool {
	result := DB.Where(
		"fid = ? and owner_id = ?",
		fid,
		owner_id,
	).Delete(&File{})
	err := result.Error
	if err != nil {
		log.Fatal(err)
	}
	affectedLines := result.RowsAffected
	return affectedLines == 1
}
func RenameFile(owner_id, fid, name string) bool {
	result := DB.Model(&File{}).Where(
		"owner_id = ? and fid = ?",
		owner_id,
		fid,
	).Updates(&File{
		Filename: name,
	})
	err := result.Error

	if err != nil {
		return false
	}
	affectedLines := result.RowsAffected
	return affectedLines == 1
}
func MoveFile(owner_id, fid, new_parent_did string) bool {
	if new_parent_did != "" && GetDir(new_parent_did, owner_id) == nil {
		return false
	}
	result := DB.Model(&File{}).Where(
		"owner_id = ? and fid = ?",
		owner_id,
		fid,
	).Updates(&File{
		ParentDiD: new_parent_did,
	})
	err := result.Error
	utils.CheckErr(err)
	affectedLines := result.RowsAffected
	return affectedLines == 1
}
func GetFileById(fid string) *File {
	var file File
	result := DB.Where(&File{
		Fid: fid,
	}).Take(&file)
	err := result.Error
	if err != nil {
		return nil
	}
	return &file
}
