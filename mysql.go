package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)
func GetData(UserEmail string ,Path string) (data_ [] PathData) {
	db := ConnectSql()
	defer db.Close()
	sql_ := fmt.Sprintf("select fileid,name,type,time,path,size from data where useremail='%v' and path='%v'",UserEmail,Path)
	rows, err := db.Query(sql_)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		var tmp PathData
		_ = rows.Scan(&tmp.DataFileId, &tmp.DataName, &tmp.DataType, &tmp.DataTime, &tmp.DataPath, &tmp.DataSize)
		data_=append(data_,tmp)
	}
	return data_
}

func ConnectSql() *sql.DB {
	cfg := OpenIni(IniPath)
	url := cfg.Section("mysql").Key("user").String()+":"+
		cfg.Section("mysql").Key("pass").String()+"@tcp("+
		cfg.Section("mysql").Key("url").String()+":"+
		cfg.Section("mysql").Key("port").String()+")/"+
		cfg.Section("mysql").Key("database").String()
	db, err := sql.Open("mysql",url)
	if err != nil {
		panic(err)
		return nil
	}
	return db
}

func GetUserInfo(email string,pw string) (UserInfo,int){
	db := ConnectSql()
	defer db.Close()
	status := 0

	var user UserInfo
	sql_ := fmt.Sprintf("select volume,groupid from user where email='%v' and pw='%v'",email,pw)
	_ = db.QueryRow(sql_).Scan(&user.Volume, &user.GroupID)
	if user.GroupID == 0 {
		return user,status
	}
	sql__ := fmt.Sprintf("SELECT SUM(SIZE) AS USED FROM `data` WHERE USEREMAIL='%v' GROUP BY USEREMAIL",email)

	_ = db.QueryRow(sql__).Scan(&user.Used)
	status = 1
	return user,status
}

//创建新的文件信息

func AddData(Data_ PathData)  {
	db := ConnectSql()
	defer db.Close()
	sql_ := fmt.Sprintf("INSERT INTO `data`(FILEID,`NAME`,TYPE,TIME,PATH,SIZE,USEREMAIL) VALUES('%v','%v','%v','%v','%v',%v,'%v')",
		Data_.DataFileId,Data_.DataName,Data_.DataType,Data_.DataTime,Data_.DataPath,Data_.DataSize,Data_.DataEmail)
	_, err := db.Exec(sql_)
	if err != nil {
		fmt.Printf("insert data error: %v\n", err)
		return
	}

}

//重命名文件 文件夹
func UpdateFile(tmp *UpdateType) string {
	db := ConnectSql()
	defer db.Close()
	num := "0"
	sql_ := fmt.Sprintf("SELECT COUNT(*) AS NUM FROM `data` WHERE USEREMAIL='%v' AND `NAME`='%v' AND TYPE='%v'",tmp.Email,tmp.NewName,tmp.Type)
	_ = db.QueryRow(sql_).Scan(&num)
	if num != "0" {
		fmt.Println("同目录下存在相同的文件(夹)名")
		return "failed"
	}
	sql_ = fmt.Sprintf("UPDATE `data` SET `NAME` = '%v',TIME = '%v' WHERE FILEID = '%v'",tmp.NewName,time.Now().Format("2006-01-02 15:04:05"),tmp.FileID)
	_, err := db.Exec(sql_)
	if err != nil {
		fmt.Printf("update data error: %v\n", err)
		return "failed"
	}
	return "succeed"
}

//删除文件 文件夹
func DeleteData(files *DeleteFile)  {
	email := "123456@qq.com"
	db := ConnectSql()
	defer db.Close()
	var sql_ string
	sql_ = fmt.Sprintf("DELETE FROM `data` WHERE FILEID='%v'",files.FileID)
	_, err := db.Exec(sql_)
	if err != nil {
		fmt.Printf("delete data error: %v\n", err)
		return
	}
	if files.Type == "dir" {
		sql_ = "DELETE FROM `data` WHERE PATH LIKE '"
		if files.Path == "/" {
			sql_ += files.Path+files.FileName+"%' "
		} else {
			sql_ += files.Path+"/"+files.FileName+"%' "
		}
		sql_+="AND USEREMAIL='" + email +"'"
		_,_ = db.Exec(sql_)
	}


}

//移动文件 文件夹
func MoveData(tmp *MoveStruct) string {
	db := ConnectSql()
	defer db.Close()
	num := "0"
	sql_ := fmt.Sprintf("SELECT COUNT(*) AS NUM FROM `data` WHERE USEREMAIL='%v' AND `PATH`='%v' AND TYPE='%v' AND NAME='%v'",tmp.Email,tmp.NewPath,tmp.Type,tmp.FileName)
	_ = db.QueryRow(sql_).Scan(&num)
	if num != "0" {
		fmt.Println("目录下存在相同的文件")
		return "failed"
	}
	if tmp.Type == "dir" {
		sql_ = fmt.Sprintf("UPDATE `data` SET `PATH` = CONCAT('%v',SUBSTRING(`PATH`,%v)) WHERE USEREMAIL='%v' AND PATH LIKE '",tmp.NewPath+"/"+tmp.FileName, len(tmp.OldPath+"/"+tmp.FileName)+1,tmp.Email)
		sql_ += tmp.OldPath+"/"+tmp.FileName+"%'"
		_, err := db.Exec(sql_)
		if err != nil {
			fmt.Printf("update data error: %v\n", err)
			return "failed"
		}
	}
	sql_ = fmt.Sprintf("UPDATE `data` SET `PATH` = '%v',TIME = '%v' WHERE FILEID = '%v'",tmp.NewPath,time.Now().Format("2006-01-02 15:04:05"),tmp.FileID)
	_, err := db.Exec(sql_)
	if err != nil {
		fmt.Printf("移动失败: %v\n", err)
		return "failed"
	}

	return "succeed"
}


