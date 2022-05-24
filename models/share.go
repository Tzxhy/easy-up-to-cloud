package models

import "gitee.com/tzxhy/web/utils"

// 只操作数据库，不涉及到实际文件
// 分享
type ShareItem struct {
	Sid string `json:"sid" gorm:"primaryKey;type:string not null;"`
	Fid string `json:"fid" gorm:"type:string;default:''"`
	// 是分享文件夹的话，那么查看分享分容时，直接索引到该did下所有文件（夹）
	Did string `json:"did" gorm:"type:string;default:''"`
	// 分享名称
	Name string `json:"name" gorm:"type:string not null;"`
	// ParentDid string `json:"parent_did" gorm:"index:resource_unique;type:string not null;default:'ROOT'"`
	// SType      RType  `json:"r_type" gorm:"type:integer not null;"`
	User User `json:"-" gorm:"references:Uid"`
	// 分享人
	UserId string `json:"-" gorm:"type:string not null;"`
	// 分享密码，可以为空
	Password   string `json:"-" gorm:"type:string default '';"`
	CreateDate int64  `json:"create_date" gorm:"autoUpdateTime:milli"`
	ExpireDate int64  `json:"expire_date" gorm:"type:integer default -1;"`
}

func AddShareItem(fid, did, name, userId, password string, expireDate int64) (bool, error) {

	sid := utils.GenerateRid()
	res := DB.Create(&ShareItem{
		Sid:        sid,
		Fid:        fid,
		Did:        did,
		Name:       name,
		UserId:     userId,
		Password:   password,
		ExpireDate: expireDate,
	})
	err := res.Error
	if err != nil {
		return false, err
	}
	return res.RowsAffected == 1, nil
}

func GetShareItem(sid string) *ShareItem {
	var shareItem ShareItem
	err := DB.Where(&ShareItem{
		Sid: sid,
	}).Take(&shareItem).Error
	if err == nil {
		return &shareItem
	}
	return nil
}
func GetAllShareItems() *[]ShareItem {
	var shareItems []ShareItem
	err := DB.Find(&shareItems).Error
	if err == nil {
		return &shareItems
	}
	return nil
}

func DeleteShare(sid, uid string) bool {
	ret := DB.Where(&ShareItem{
		Sid:    sid,
		UserId: uid,
	}).Delete(&ShareItem{})
	err := ret.Error

	if err != nil {
		return false
	}
	return ret.RowsAffected == 1
}
func DeleteShareByFid(fid string) uint8 {
	if fid == "" {
		return 0
	}
	ret := DB.Where(&ShareItem{
		Fid: fid,
	}).Delete(&ShareItem{})
	err := ret.Error
	if err != nil {
		return 0
	}
	return uint8(ret.RowsAffected)
}

func DeleteShareByDid(did string) uint8 {
	if did == "" {
		return 0
	}
	ret := DB.Where(&ShareItem{
		Did: did,
	}).Delete(&ShareItem{})
	err := ret.Error
	if err != nil {
		return 0
	}
	return uint8(ret.RowsAffected)
}
