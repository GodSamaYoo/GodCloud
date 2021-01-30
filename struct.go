package main

import (
	"github.com/zyxar/argo/rpc"
)

//数据库模型

//用户模型
type Users struct {
	UserID int `gorm:"primaryKey"`
	Email string
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
}

//用户组

type Usergroups struct {
	GroupID int `gorm:"primaryKey"`
	Name string
	Volume int
}

//下载任务模型
type Task struct {
	Gid string `gorm:"primaryKey"`
	UserEmail string
	Path string
	Time int `gorm:"autoCreateTime;autoUpdateTime"`
	TmpPath string
}


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
	GroupID int
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
