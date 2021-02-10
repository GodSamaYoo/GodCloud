package main

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

func RegisterRoutes(e *echo.Echo) {

	//文件目录

	e.GET("/api/file/:path", func(ctx echo.Context) error {
		path, _ := url.QueryUnescape(ctx.Param("path"))
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		return ctx.JSON(http.StatusOK, GetData(email, path))
	})

	//文件获取下载

	e.GET("/api/file/:filename/:base64path", func(ctx echo.Context) error {
		filename, _ := url.QueryUnescape(ctx.Param("filename"))
		base64path := ctx.Param("base64path")
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		path, _ := base64.StdEncoding.DecodeString(base64path)
		if string(path) == "/" {
			return ctx.Attachment(root+email+"/"+filename, filename)
		}
		return ctx.Attachment(root+email+string(path)+"/"+filename, filename)
	})

	//文件缩略图
	e.GET("/api/thumbnail/:fileid", func(ctx echo.Context) error {
		fileid := ctx.Param("fileid")
		return ctx.Attachment("Thumbnail/"+fileid+".jpg",fileid+".jpg")
	})

	//文件上传 / 文件夹创建

	e.POST("/api/file/", func(ctx echo.Context) error {

		err := FileUploadAndCreate(ctx)
		if !err  {
			fmt.Println("上传/创建失败")
			return ctx.String(http.StatusOK, "failed")
		}
		return ctx.String(http.StatusOK, "succeed")
	})

	//重命名文件(夹)

	e.PUT("/api/file/", func(ctx echo.Context) error {
		tmp := new(UpdateType)
		_ = ctx.Bind(tmp)
		if UpdateFile(tmp) == "succeed" {
			UpdateFiles(tmp)
			return ctx.JSON(200, "succeed")
		}
		return ctx.JSON(200, "failed")
	})

	//移动文件 文件夹
	e.PUT("/api/move/", func(ctx echo.Context) error {
		tmp := new(MoveStruct)
		_ = ctx.Bind(tmp)
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		tmp.Email = email
		if MoveFiles(tmp) != "succeed" {
			return ctx.JSON(200, "failed")
		}
		MoveData(tmp)
		return ctx.JSON(200, "succeed")

	})

	//删除文件 文件夹

	e.DELETE("/api/file/", func(ctx echo.Context) error {
		tmp := new(DeleteFile)
		_ = ctx.Bind(tmp)
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if IsOneFile(tmp.FileID) {
			tt := DatasInfoQuery(&Datas{
				FileID: tmp.FileID,
			})
			if DeleteOneDriveFile(tt.ItemID,tt.StoreID) == "succeed" {
				DeleteData(tmp,email)
				return ctx.JSON(200, "succeed")
			}
		}else {
			if DeleteFiles(tmp,email) == "succeed" {
				DeleteData(tmp,email)
				return ctx.JSON(200, "succeed")
			}
		}
		return ctx.JSON(200, "failed")
	})

	//用户登录验证

	e.POST("/api/login", func(ctx echo.Context) error {
		email := ctx.FormValue("email")
		pw := ctx.FormValue("pw")
		pw = md5_(pw)
		rememberme := ctx.FormValue("checked")
		status := GetLogin(email, pw)
		if status == 0 {
			return ctx.JSON(http.StatusOK, "failed")
		}
		cookie := new (http.Cookie)
		cookie.Name = "GODKEY"
		cookie.Value,_ = DesEncrypt(email,deskey)
		cookie.Path = "/"
		if rememberme == "true" {
			cookie.Expires = time.Now().Add(7 * 24 * time.Hour)
		}
		ctx.SetCookie(cookie)
		return ctx.JSON(http.StatusOK, "succeed")
	})

	//用户容量查询
	e.POST("/api/cookies", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		userinfo_,_ := GetUserInfo(email)
		return ctx.JSON(200,userinfo_)
	})

	//解压缩文件

	e.PUT("/api/UnArchive", func(ctx echo.Context) error {
		tmp := new(UnArchiveFile)
		_ = ctx.Bind(tmp)
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !UnArchiver(tmp,email) {
			return ctx.JSON(200, "failed")
		}
		return ctx.JSON(200, "succeed")
	})
	
	//aria2配置设置
	e.POST("/api/aria2config", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := new(aria2config)
		_ = ctx.Bind(tmp)
		ModifyIni("aria2","enable",tmp.Status)
		ModifyIni("aria2","port",tmp.Port)
		ModifyIni("aria2","token",tmp.Secret)
		ModifyIni("aria2","tmpdownpath",tmp.TmpDownPath)
		ModifyIni("aria2","time",tmp.Time)
		if aria2client != nil {
			_ = aria2client.Close()
		}
		aria2client = aria2begin()
		return ctx.JSON(200,"succeed")
	})

	//读取aria2配置
	e.GET("/api/aria2config", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := aria2config{
			Status:      ReadIni("aria2","enable"),
			Port:        ReadIni("aria2","port"),
			Secret:      ReadIni("aria2","token"),
			TmpDownPath: ReadIni("aria2","tmpdownpath"),
			Time:        ReadIni("aria2","time"),
		}
		return ctx.JSON(200,tmp)
	})

	e.GET("/api/aria2test", func(ctx echo.Context) error {
		if  aria2client == nil {
			return ctx.JSON(200,"failed")
		}
		version,err := aria2client.GetVersion()
		if err != nil {
			return ctx.JSON(200,"failed")
		}
		return ctx.JSON(200,version.Version)
	})


	//添加aria2下载
	e.POST("/api/aria2", func(ctx echo.Context) error {
		tmp := new(aria2accepturl)
		_ = ctx.Bind(tmp)
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		a := aria2download(tmp.Url,tmp.Path,email)
		return ctx.JSON(200, len(a))
	})

	//获取aria2下载状态
	e.GET("/api/aria2", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		infos := aria2status(email)

		return ctx.JSON(200,infos)
	})


	//修改aria2下载状态
	e.PUT("/api/aria2", func(ctx echo.Context) error {
		tmp := new(aria2change)
		_ = ctx.Bind(tmp)
		err := aria2taskchange(tmp)
		if err != nil {
			fmt.Println(err)
			return ctx.JSON(200,"failed")
		}
		return ctx.JSON(200,"succeed")
	})

	//储存策略新增
	e.POST("/api/store", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return err
		}
		tmp := new(storemodel)
		_ = ctx.Bind(tmp)
		if tmp.Type == "onedrive" {
			err = FirstGetRefreshToken(tmp)
			if err != nil {
				i :=0
				for {
					err = FirstGetRefreshToken(tmp)
					if err == nil || i > 4 {
						break
					}
					i++
				}
			}
			RefreshToken2()
		}else {
			StoreAdd(tmp)
		}
		return ctx.JSON(200,"succeed")
	})

	//查询储存策略
	e.GET("/api/store", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := StoreQuery()
		return ctx.JSON(200,tmp)
	})

	//查询用户
	e.GET("/api/user", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := UserQuery()
		return ctx.JSON(200,tmp)
	})

	//增加用户
	e.POST("/api/user", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := new(useradd)
		_ = ctx.Bind(tmp)
		if AddUser(tmp) {
			return ctx.JSON(200,"succeed")
		}
		return  ctx.JSON(400,"failed")
	})

	//删除用户
	e.DELETE("/api/user", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		email_ := ctx.FormValue("email")
		if DeleteUser(email_) {
			return ctx.JSON(200,"succeed")
		}
		return  ctx.JSON(400,"failed")
	})


	//查询用户组
	e.GET("/api/usergroup", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := UserGroupInfo()
		return ctx.JSON(200,tmp)
	})

	//增加用户组
	e.POST("/api/usergroup", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		tmp := new(usergroupadd)
		_ = ctx.Bind(tmp)
		if  AddUserGroup(tmp) {
			return ctx.JSON(200,"succeed")
		}
		return  ctx.JSON(400,"failed")
	})

	//删除用户组
	e.DELETE("/api/usergroup", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		if !IsAdmin(email) {
			return ctx.JSON(400,"failed")
		}
		groupid, _ := strconv.Atoi(ctx.FormValue("groupid"))
		if DeleteUserGroup(groupid) {
			return ctx.JSON(200,"succeed")
		}
		return  ctx.JSON(400,"failed")
	})

	//onedrive上传地址获取
	e.POST("/api/onedriveuploadaddress", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		path := ctx.FormValue("path")
		filename := ctx.FormValue("filename")
		filesize,_ := strconv.Atoi(ctx.FormValue("filesize"))
		tmp := GetOneDriveAdd(email,path,filename,filesize)
		if tmp.Address != "" {
			return ctx.JSON(200,tmp)
		}
		return ctx.JSON(200,"failed")
	})

	//onedrive上传回调
	e.POST("/api/onedrivestoredata", func(ctx echo.Context) error {
		enkey,err := ctx.Cookie("GODKEY")
		if err != nil {
			return err
		}
		email,_ := DesDecrypt(enkey.Value,deskey)
		tmp := new(OnedriveDatas)
		_ = ctx.Bind(tmp)
		tmp.UserEmail = email
		picture := []string{"jpg", "jpeg", "bmp", "gif", "png", "tif"}
		tp := path.Ext(tmp.Name)
		url_ := ""
		if tp != "" {
			for _,v_ := range picture{
				if tp[1:] == v_{
					url_ = GetThumbnail(tmp.ItemID,tmp.StoreID)
					break
				}
			}
		}
		fileid := md5_(tmp.Name + tmp.Path + time.Now().String())
		if url_ != "" {
			os.Mkdir("Thumbnail",0777)
			file, _ := http.Get(url_)
			defer file.Body.Close()
			files,_ := ioutil.ReadAll(file.Body)
			_ = ioutil.WriteFile("Thumbnail/"+fileid+".jpg", files, 0644)
		}
		DatasAdd(&Datas{
			FileID:    fileid,
			Name:      tmp.Name,
			Type:      tmp.Type,
			Path:      tmp.Path,
			UserEmail: tmp.UserEmail,
			Size:      tmp.Size,
			StoreID:   tmp.StoreID,
			ItemID:    tmp.ItemID,
		})
		return ctx.JSON(200,"succeed")
	})


}
