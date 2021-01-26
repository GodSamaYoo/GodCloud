package main

//数据库模型
type Users struct {
	UserID int `gorm:"primaryKey"`
	Email string
	Password string
	GroupID int
	Volume int
	Time int `gorm:"autoCreateTime"`
}

type Datas struct {
	FileID string `gorm:"primaryKey"`
	Name string
	Type string
	Time int `gorm:"autoCreateTime;autoUpdateTime"`
	Path string
	UserEmail string `gorm:"index"`
	Size int
}

type Usergroups struct {
	GroupID int
	name string
	volume int
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
