package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("连接Sqlite数据库失败")
	}
	return db
}

func GetData(UserEmail string ,Path string) (data_ [] PathData) {
	db := ConnectSqlite()
	var data []Datas
	db.Where(&Datas{UserEmail:UserEmail,Path: Path}).Find(&data)
	for _,v := range data {
		tmp := PathData{
			DataFileId: v.FileID,
			DataName:   v.Name,
			DataType:   v.Type,
			DataTime:   v.Time,
			DataPath:   v.Path,
			DataSize:   v.Size,
		}
		data_=append(data_, tmp)
	}
	return data_
}

func GetUserInfo(email string,pw string) (UserInfo,int){
	db := ConnectSqlite()
	status := 0
	//var userinfo UserInfo
	var user Users
	row := db.Where(&Users{Email: email,Password: pw}).Find(&user)
	userinfo := UserInfo{
		Volume:  user.Volume,
		GroupID: user.GroupID,
	}
	if row.RowsAffected !=1 {
		return UserInfo{},status
	}
	status = 1
	db.Raw("SELECT SUM(size) AS USED FROM `datas` WHERE user_email='?' GROUP BY user_email",email).Scan(&userinfo.Used)
	return userinfo,status
}

//创建新的文件信息
func AddData(Data_ PathData)  {
	email := "admin@godcloud.com"
	db := ConnectSqlite()
	db.Create(&Datas{
		FileID:    Data_.DataFileId,
		Name:      Data_.DataName,
		Type:      Data_.DataType,
		Path:      Data_.DataPath,
		UserEmail: email,
		Size:      Data_.DataSize,
	})
}

//重命名文件 文件夹
func UpdateFile(tmp *UpdateType) string {
	db := ConnectSqlite()
	num := db.Where(&Datas{UserEmail: tmp.Email,Name: tmp.NewName,Type: tmp.Type}).Find(&Datas{}).RowsAffected
	if num != 0 {
		fmt.Println("同目录下存在相同的文件(夹)名")
		return "failed"
	}
	if tmp.Type == "dir" {
		var oldpath,newpath string
		if tmp.Path == "/" {
			oldpath = tmp.Path + tmp.OldName+"%"
			newpath = tmp.Path + tmp.NewName
		} else {
			oldpath= tmp.Path+"/" + tmp.OldName+"%"
			newpath = tmp.Path+"/" + tmp.NewName
		}
		db.Model(&Datas{}).Where("path LIKE ? AND user_email",oldpath,tmp.Email).Update("path",newpath)
	}
	db.Model(&Datas{}).Where("file_id = ?",tmp.FileID).Updates(Datas{
		Name:      tmp.NewName,
	})
	return "succeed"
}


//删除文件 文件夹
func DeleteData(files *DeleteFile)  {
	email := "admin@godcloud.com"
	db := ConnectSqlite()
	db.Delete(&Datas{FileID: files.FileID})
	if files.Type == "dir" {
		var path string
		if files.Path == "/" {
			path = files.Path+files.FileName+"%"
		} else {
			path= files.Path+"/"+files.FileName+"%"
		}
		db.Delete(Datas{UserEmail: email},"path LIKE ?",path)
		fmt.Println(email+path)
	}
}

//移动文件 文件夹
func MoveData(tmp *MoveStruct) string {
	db := ConnectSqlite()
	num := db.Where(&Datas{UserEmail: tmp.Email,Path: tmp.NewPath,Type: tmp.Type,Name: tmp.FileName}).RowsAffected
	if num != 0 {
		fmt.Println("目录下存在相同的文件")
		return "failed"
	}
	if tmp.Type == "dir" {
		sql_ := fmt.Sprintf("UPDATE `datas` SET `PATH` = CONCAT('%v',SUBSTRING(`PATH`,%v)) WHERE USEREMAIL='%v' AND PATH LIKE '",tmp.NewPath+"/"+tmp.FileName, len(tmp.OldPath+"/"+tmp.FileName)+1,tmp.Email)
		sql_ += tmp.OldPath+"/"+tmp.FileName+"%'"
		db.Raw(sql_)
	}
	db.Model(&Datas{}).Where("file_id = ?",tmp.FileID).Updates(Datas{Path: tmp.NewPath})

	return "succeed"
}