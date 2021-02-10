package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
	"strings"
)

func CheckSqlite() {
	_, err := os.Stat("data.db")
	if err != nil {
		dbs, err_ := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err_ != nil {
			fmt.Println("创建Sqlite数据库失败")
		}
		db = *dbs

		_ = db.AutoMigrate(&Users{})
		_ = db.AutoMigrate(&Datas{})
		_ = db.AutoMigrate(&Usergroups{})
		_ = db.AutoMigrate(&Task{})
		_ = db.AutoMigrate(&Store{})
		db.Create(&Users{
			UserID:   1,
			Email:    "admin@godcloud.com",
			Password: "49ba59abbe56e057",
			GroupID:  1,
			Volume:   1048576,
		})
		db.Create(&Usergroups{
			GroupID: 1,
			Name:    "管理员",
			Volume:  1048576,
			StoreID: "1",
		})
		db.Create(&Store{
			ID:           1,
			Name:         "本地存储",
			Path:         "./",
			Volume:       1048576,
			Type:         "local",
		})
		fmt.Println("用户名：admin@godcloud.com\n密码：123456")
	}else {
		ConnectSqlite()
	}
}

func ConnectSqlite() {
	dbs, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("连接Sqlite数据库失败")
	}
	db = *dbs
}

func GetData(UserEmail string, Path string) (data_ []PathData) {
	var data []Datas
	db.Where(&Datas{UserEmail: UserEmail, Path: Path}).Find(&data)
	for _, v := range data {
		tmp := PathData{
			DataFileId: v.FileID,
			DataName:   v.Name,
			DataType:   v.Type,
			DataTime:   v.Time,
			DataPath:   v.Path,
			DataSize:   v.Size,
		}
		data_ = append(data_, tmp)
	}
	return data_
}

func GetUserInfo(email string) (UserInfo, int) {
	status := 0
	var user Users
	row := db.Where(&Users{Email: email}).Find(&user)
	userinfo_ := UserInfo{
		Volume:  user.Volume,
		RegisterTime: user.Time,
		GroupID: user.GroupID,
	}
	if row.RowsAffected != 1 {
		return UserInfo{}, status
	}
	status = 1
	db.Raw("SELECT SUM(size) AS USED FROM `datas` WHERE user_email=? GROUP BY user_email", email).Scan(&userinfo_.Used)
	db.Raw("SELECT name FROM usergroups WHERE group_id =?",user.GroupID).Scan(&userinfo_.GroupName)
	return userinfo_, status
}

//验证登录
func GetLogin(email string, pw string) int {
	status := 0
	var user Users
	row := db.Where(&Users{Email: email, Password: pw}).Find(&user)
	if row.RowsAffected != 1 {
		return status
	}
	status = 1
	return status
}

//验证是否为管理员
func IsAdmin(email string) bool {
	var user Users
	row := db.Where(&Users{Email: email,GroupID: 1}).Find(&user).RowsAffected
	if row == 1 {
		return true
	}
	return false
}

//创建新的文件信息
func AddData(Data_ PathData,email string) bool {
	var datas Datas
	num := db.Where(&Datas{UserEmail: email,Path: Data_.DataPath,Name: Data_.DataName}).Find(&datas).RowsAffected
	if num == 1 {
		return false
	}
	db.Create(&Datas{
		FileID:    Data_.DataFileId,
		Name:      Data_.DataName,
		Type:      Data_.DataType,
		Path:      Data_.DataPath,
		UserEmail: email,
		Size:      Data_.DataSize,
	})
	return true
}

//重命名文件 文件夹
func UpdateFile(tmp *UpdateType) string {
	num := db.Where(&Datas{UserEmail: tmp.Email, Name: tmp.NewName, Type: tmp.Type, Path: tmp.Path}).Find(&Datas{}).RowsAffected
	if num != 0 {
		fmt.Println("同目录下存在相同的文件(夹)名")
		return "failed"
	}
	if tmp.Type == "dir" {
		var oldpath, newpath string
		if tmp.Path == "/" {
			oldpath = tmp.Path + tmp.OldName + "%"
			newpath = tmp.Path + tmp.NewName
		} else {
			oldpath = tmp.Path + "/" + tmp.OldName + "%"
			newpath = tmp.Path + "/" + tmp.NewName
		}
		db.Model(&Datas{}).Where("path LIKE ? AND user_email = ?", oldpath, tmp.Email).Update("path", newpath)
	}
	db.Model(&Datas{}).Where("file_id = ?", tmp.FileID).Updates(Datas{
		Name: tmp.NewName,
	})
	return "succeed"
}

//删除文件 文件夹
func DeleteData(files *DeleteFile,email string) {
	db.Delete(&Datas{FileID: files.FileID})
	if files.Type == "dir" {
		var path string
		if files.Path == "/" {
			path = files.Path + files.FileName + "%"
		} else {
			path = files.Path + "/" + files.FileName + "%"
		}
		db.Delete(Datas{UserEmail: email}, "path LIKE ?", path)
	}
}

//移动文件 文件夹
func MoveData(tmp *MoveStruct) string {
	num := db.Where(&Datas{UserEmail: tmp.Email, Path: tmp.NewPath, Type: tmp.Type, Name: tmp.FileName}).RowsAffected
	if num != 0 {
		fmt.Println("目录下存在相同的文件")
		return "failed"
	}
	if tmp.Type == "dir" {
		sql_ := fmt.Sprintf("UPDATE `datas` SET `PATH` = CONCAT('%v',SUBSTRING(`PATH`,%v)) WHERE USEREMAIL='%v' AND PATH LIKE '", tmp.NewPath+"/"+tmp.FileName, len(tmp.OldPath+"/"+tmp.FileName)+1, tmp.Email)
		sql_ += tmp.OldPath + "/" + tmp.FileName + "%'"
		db.Raw(sql_)
	}
	db.Model(&Datas{}).Where("file_id = ?", tmp.FileID).Updates(Datas{Path: tmp.NewPath})

	return "succeed"
}

//aria2创建下载相关信息

func Aria2DataBegin(gid string,path string,email string,tmp string)  {
	db.Create(&Task{Gid: gid,Path: path,UserEmail:email,TmpPath: tmp})
}

//aria2查询信息

//通过邮件
func Aria2DataInfo(email string) []Task{
	var task []Task
	db.Where("user_email = ?",email).Find(&task)
	return task
}
//通过地址
func Aria2DataByAdd(tmp string) Task{
	var task Task
	db.Where("tmp_path=?",tmp).Find(&task)
	return task
}

//aria2删除任务

func Aria2DeleteTask(tmp string)  {
	db.Delete(&Task{},"tmp_path = ?",tmp)
}

//用户修改密码

func Changepw(email string,pw string) bool {
	num := db.Model(&Users{}).Where("email = ?",email).Updates(Users{Password: md5_(pw)}).RowsAffected
	if num == 1 {
		return true
	}
	return false
}
//查询用户

func UserQuery() []userinfo {
	var users []userinfo
	db.Raw("SELECT users.user_id AS ID,users.email,usergroups.name,users.volume,users.time,SUM(size) AS used FROM (users LEFT JOIN usergroups ON users.group_id = usergroups.group_id) LEFT JOIN datas ON datas.user_email = users.email GROUP BY users.email ORDER BY user_id").Scan(&users)
	return users
}

//增加用户

func AddUser(tmp *useradd) bool {
	var usergroup Usergroups
	db.Where(&Usergroups{
		GroupID: tmp.GroupID,
	}).Find(&usergroup)
	user := Users{
		Email:    tmp.Email,
		Password: md5_(tmp.Password),
		GroupID:  tmp.GroupID,
		Volume:   usergroup.Volume,
	}
	num := db.Create(&user).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

//删除用户

func DeleteUser(email string) bool {
	db.Delete(&Datas{},"user_email = ?",email)
	num := db.Delete(&Users{},"email = ?",email).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

//增加用户组

func AddUserGroup(tmp *usergroupadd) bool {
	num := db.Create(&Usergroups{
		Name: tmp.Name,
		Volume: tmp.Volume,
		StoreID: tmp.StoreID,
	}).RowsAffected
	if num == 1 {
		return true
	}
	return false
}
//用户组查询
func UserGroupInfo() []Usergroups {
	var tmp []Usergroups
	db.Find(&tmp)
	return tmp
}

//删除用户组

func DeleteUserGroup(groupid int) bool {
	var user []Users
	db.Where(&Users{
		GroupID: groupid,
	}).Find(&user)
	for _,v := range user {
		DeleteUser(v.Email)
	}
	num := db.Delete(&Usergroups{},"group_id = ?",groupid).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

//修改用户组名字容量

func ChangeUserGroup(name string,volume int,groupid int) bool {
	num := db.Model(&Usergroups{}).Where("group_id = ?",groupid).Updates(Usergroups{Name: name,Volume: volume}).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

//储存策略增加
func StoreAdd(tmp *storemodel) bool {
	store := Store{
		Name:         tmp.Name,
		Type:         tmp.Type,
		Path:         tmp.Path,
		Volume:       tmp.Volume,
		RefreshToken: tmp.RefreshToken,
		ClientSecret: tmp.ClientSecret,
		ClientID:     tmp.ClientID,
	}
	num := db.Create(&store).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

//储存策略查询

func StoreQuery() []storeinfo {
	var stores []Store
	var tmp []storeinfo
	db.Find(&stores)
	for _,v := range stores {
		tmp_ := storeinfo{
			ID:     v.ID,
			Name:   v.Name,
			Type:   v.Type,
			Used:   v.Used,
			Volume: v.Volume,
		}
		tmp = append(tmp, tmp_)
	}
	return tmp
}

//储存策略容量查询

func StoreInfoQuery(tmp *Store) *Store {
	tmp_ := new(Store)
	db.Where(tmp).Find(&tmp_)
	return tmp_
}

func UserInfoQuery(tmp *Users) *Users {
	tmp_ := new(Users)
	db.Where(tmp).Find(&tmp_)
	return tmp_
}

func UserGroupQuery(tmp *Usergroups) *Usergroups {
	tmp_ := new(Usergroups)
	db.Where(tmp).Find(&tmp_)
	return tmp_
}

func DatasInfoQuery(tmp *Datas) *Datas {
	tmp_ := new(Datas)
	db.Where(tmp).Find(&tmp_)
	return tmp_
}

func OneDriveRefreshUpdate(tmp *Store) {
	db.Model(&Store{}).Where("id = ?",tmp.ID).Updates(tmp)
}

func DatasAdd(tmp *Datas)  {
	tmp_ := new(Datas)
	num := db.Where(&Datas{
		UserEmail: tmp.UserEmail,
		Path: tmp.Path,
		Name: tmp.Name,
	}).Find(&tmp_).RowsAffected
	if num > 0 {
		db.Model(&Datas{}).Where("file_id = ?",tmp_.FileID).Updates(tmp)
	}else {
		db.Create(&tmp)
	}
}

func IsOnedriveFile(storeid int) int {
	var tmp Store
	num := db.Where(&Store{
		ID: storeid,
		Type: "onedrive",
	}).Find(&tmp).RowsAffected

	if num == 1 {
		return storeid
	}
	return -1
}

func IsOneFile(fileid string) bool {
	var tmp Datas
	db.Where("file_id=?",fileid).Find(&tmp)
	num := IsOnedriveFile(tmp.StoreID)
	if num == -1 {
		return false
	}
	return true
}

func IsOneDriveStore(email string) bool {
	var a Users
	var b Usergroups
	var e Store
	db.Where("email = ?",email).Find(&a)
	db.Where(&Usergroups{
		GroupID: a.GroupID,
	}).Find(&b)
	c :=strings.Split(b.StoreID,",")[0]
	d ,_ := strconv.Atoi(c)
	num := db.Where(&Store{
		ID: d,
		Type: "onedrive",
	}).Find(&e).RowsAffected
	if num == 1 {
		return true
	}
	return false
}

