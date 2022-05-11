package models

import (
	"errors"
	"log"

	"gitee.com/tzxhy/web/constants"
	"gitee.com/tzxhy/web/utils"
)

type File struct {
	Fid      string `json:"fid"`
	OwnerId  string
	Filename string `json:"filename"`
	Filesize int    `json:"file_size"`
	// -1 为根目录
	ParentDiD    string `json:"parent_did"`
	FileRealPath string
	CreateDate   string `json:"create_date"`
}

//
func AddFile(owner_id string, dir_id string, filename string, file_size uint64, file_path string) (string, error) {

	originDir := GetFileByName(filename, owner_id, dir_id)
	if originDir != nil { // 已有同名
		return "", errors.New(constants.TIPS_HAS_SAME_FILE)
	}

	stmt, err := DB.Prepare("insert into files (fid, owner_id, filename, parent_did, file_real_path, file_size) values(?, ?, ?, ?, ?, ?)")
	utils.CheckErr(err)
	defer stmt.Close()
	fid := utils.RandStringBytesMaskImprSrc(8)
	_, err = stmt.Exec(fid, owner_id, filename, dir_id, file_path, file_size)
	utils.CheckErr(err)
	return fid, nil
}

func GetFileByName(filename string, owner_id string, parent_did string) *File {
	rows, err := DB.Query("select * from files where owner_id = ? and filename = ? and parent_did = ?", owner_id, filename, parent_did)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var file *File
	for rows.Next() {
		file = new(File)
		rows.Scan(
			&file.Fid,
			&file.OwnerId,
			&file.Filename,
			&file.Filesize,
			&file.ParentDiD,
			&file.FileRealPath,
			&file.CreateDate,
		)
		break
	}
	return file
}
func GetFile(fid string, owner_id string) *File {
	rows, err := DB.Query("select * from files where owner_id = ? and fid = ?", owner_id, fid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var file *File
	for rows.Next() {
		file = new(File)
		rows.Scan(
			&file.Fid,
			&file.OwnerId,
			&file.Filename,
			&file.Filesize,
			&file.ParentDiD,
			&file.FileRealPath,
			&file.CreateDate,
		)
		break
	}
	return file
}

func GetFileList(parent_id, owner_id string) *[]File {
	rows, err := DB.Query("select * from files where owner_id = ? and parent_did = ?", owner_id, parent_id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		file := new(File)
		rows.Scan(
			&file.Fid,
			&file.OwnerId,
			&file.Filename,
			&file.Filesize,
			&file.ParentDiD,
			&file.FileRealPath,
			&file.CreateDate,
		)
		files = append(files, *file)
	}
	return &files
}

func SearchFileList(owner_id, name string) *[]File {
	rows, err := DB.Query("select * from files where owner_id = ? and filename like ?", owner_id, "%"+name+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		file := new(File)
		rows.Scan(
			&file.Fid,
			&file.OwnerId,
			&file.Filename,
			&file.Filesize,
			&file.ParentDiD,
			&file.FileRealPath,
			&file.CreateDate,
		)
		files = append(files, *file)
	}
	return &files
}

func DeleteFile(fid string, owner_id string) bool {
	stmt, err := DB.Prepare("delete from files where fid = ? and owner_id = ?")
	utils.CheckErr(err)
	defer stmt.Close()
	result, err := stmt.Exec(fid, owner_id)
	utils.CheckErr(err)
	affectedLines, err := result.RowsAffected()
	utils.CheckErr(err)
	return affectedLines == 1
}
func RenameFile(owner_id, fid, name string) bool {
	stmt, err := DB.Prepare("update files set filename = ? where owner_id = ? and fid = ?")
	utils.CheckErr(err)
	defer stmt.Close()
	result, err := stmt.Exec(name, owner_id, fid)
	utils.CheckErr(err)
	affectedLines, err := result.RowsAffected()
	utils.CheckErr(err)
	return affectedLines == 1
}
func MoveFile(owner_id, fid, new_parent_did string) bool {
	if new_parent_did != "" && GetDir(new_parent_did, owner_id) == nil {
		return false
	}
	stmt, err := DB.Prepare("update files set parent_did = ? where owner_id = ? and fid = ?")
	utils.CheckErr(err)
	defer stmt.Close()
	result, err := stmt.Exec(new_parent_did, owner_id, fid)
	utils.CheckErr(err)
	affectedLines, err := result.RowsAffected()
	utils.CheckErr(err)
	return affectedLines == 1
}
