package main

import (
	"github.com/zyxar/argo/rpc"
)

//数据库模型

//用户模型
type Users struct {
	UserID int `gorm:"primaryKey;autoIncrement"`
	Email string `gorm:"unique"`
	Password string
	GroupID int
	Volume int
	Time int `gorm:"autoCreateTime"`
}

//文件 文件夹 模型

type Datas struct {
	FileID string `gorm:"primaryKey"`
	Name string
	Type string
	Time int `gorm:"autoCreateTime;autoUpdateTime"`
	Path string
	UserEmail string `gorm:"index"`
	Size int
	StoreID int
	ItemID string
}

//用户组

type Usergroups struct {
	GroupID int `gorm:"primaryKey;autoIncrement"`
	Name string
	Volume int
	StoreID string
}

//下载任务模型
type Task struct {
	Gid string `gorm:"primaryKey"`
	UserEmail string
	Path string
	Time int `gorm:"autoCreateTime;autoUpdateTime"`
	TmpPath string
}

//存储模型
type Store struct {
	ID int `gorm:"primaryKey;autoIncrement"`
	Name string
	Type string
	Path string
	Volume int
	RefreshToken string
	ClientSecret string
	ClientID string
	Time int `gorm:"autoCreateTime;autoUpdateTime"`
	Used int
}

//-----------------------------------------------------------------------------------

//接口模型
type PathData struct {
	DataFileId string
	DataName string
	DataType string
	DataTime int
	DataPath string
	DataSize int
}

type UserInfo struct {
	Volume int
	Used int
	GroupName string
	GroupID int
	RegisterTime int
	StoreType int
}

type UpdateType struct {
	OldName string `json:"OldName" form:"OldName" query:"OldName"`
	NewName string `json:"NewName" form:"NewName" query:"NewName"`
	Path string `json:"Path" form:"Path" query:"Path"`
	FileID string `json:"FileID" form:"FileID" query:"FileID"`
	Type string `json:"Type" form:"Type" query:"Type"`
	Email string `json:"Email" form:"Email" query:"Email"`
}

type MoveStruct struct {
	FileID string `json:"FileID" form:"FileID" query:"FileID"`
	FileName string `json:"FileName" form:"FileName" query:"FileName"`
	OldPath string `json:"OldPath" form:"OldPath" query:"OldPath"`
	NewPath string `json:"NewPath" form:"NewPath" query:"NewPath"`
	Type string `json:"Type" form:"Type" query:"Type"`
	Email string `json:"Email" form:"Email" query:"Email"`
}

type DeleteFile struct {
	FileID string `json:"FileID" form:"FileID" query:"FileID"`
	FileName string `json:"FileName" form:"FileName" query:"FileName"`
	Path string `json:"Path" form:"Path" query:"Path"`
	Type string `json:"Type" form:"Type" query:"Type"`
}

type UnArchiveFile struct {
	FileName string `json:"FileName" form:"FileName" query:"FileName"`
	Path string `json:"Path" form:"Path" query:"Path"`
	NewPath string `json:"NewPath" form:"NewPath" query:"NewPath"`
	PassWord string `json:"PassWord" form:"PassWord" query:"PassWord"`
	Email string `json:"Email" form:"Email" query:"Email"`
}

type aria2accepturl struct {
	Url  []string `json:"url" form:"url" query:"url"`
	Path string   `json:"path" form:"path" query:"path"`
}

type aria2downloadinfo struct {
	Gid string
	Status string
	FileNums int
	CompletedLength string
	DownloadSpeed string
	TotalLength string
	BeginTime int
	Path string
	Files []rpc.FileInfo
}

type aria2change struct {
	Gid string `json:"gid" form:"gid" query:"gid"`
	Type int `json:"type" form:"type" query:"type"`
}

type storemodel struct {
	Name string `json:"name" form:"name" query:"name"`
	Type string `json:"type" form:"type" query:"type"`
	Volume int `json:"volume" form:"volume" query:"volume"`
	Path string `json:"path" form:"path" query:"path"`
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
	ClientSecret string `json:"client_secret" form:"client_secret" query:"client_secret"`
	ClientID string `json:"client_id" form:"client_id" query:"client_id"`
	Code string `json:"code" form:"code" query:"code"`
	RedirectUrl string `json:"redirect_uri" form:"redirect_uri" query:"redirect_uri"`
}

type onedrivererefresh struct {
	RefreshToken string `json:"refresh_token" form:"refresh_token" query:"refresh_token"`
	Token string `json:"access_token" form:"access_token" query:"access_token"`
}

type storeinfo struct {
	ID int
	Name string
	Type string
	Used int
	Volume int
}

type userinfo struct {
	ID int
	Email string
	Name string
	Volume int
	Time int
	Used int
}

type useradd struct {
	Email string `json:"Email" form:"Email" query:"Email"`
	Password string `json:"pw" form:"pw" query:"pw"`
	GroupID int `json:"GroupID" form:"GroupID" query:"GroupID"`
	Volume int `json:"Volume" form:"Volume" query:"Volume"`
}

type usergroupadd struct {
	Name string `json:"Name" form:"Name" query:"Name"`
	Volume int  `json:"Volume" form:"Volume" query:"Volume"`
	StoreID string `json:"StoreID" form:"StoreID" query:"StoreID"`
}

type aria2config struct {
	Status string `json:"Status" form:"Status" query:"Status"`
	Port string `json:"Port" form:"Port" query:"Port"`
	Secret string `json:"Secret" form:"Secret" query:"Secret"`
	TmpDownPath string `json:"TmpDownPath" form:"TmpDownPath" query:"TmpDownPath"`
	Time string `json:"Time" form:"Time" query:"Time"`
}

type onedriveAddress struct {
	UploadAddress string `json:"uploadUrl" form:"uploadUrl" query:"uploadUrl"`
}

type OnedriveDatas struct {
	FileID string `json:"FileID" form:"FileID" query:"FileID"`
	Name string `json:"Name" form:"Name" query:"Name"`
	Type string `json:"Type" form:"Type" query:"Type"`
	Path string `json:"Path" form:"Path" query:"Path"`
	UserEmail string `json:"UserEmail" form:"UserEmail" query:"UserEmail"`
	Size int `json:"Size" form:"Size" query:"Size"`
	StoreID int `json:"StoreID" form:"StoreID" query:"StoreID"`
	ItemID string `json:"ItemID" form:"ItemID" query:"ItemID"`
}

type OnedriveAddReturn struct {
	Address string
	StoreID int
}
