package main

type PathData struct {
	DataFileId string
	DataName string
	DataType string
	DataTime string
	DataPath string
	DataSize int
	DataEmail string
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
