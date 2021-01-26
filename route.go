package main

import (
	"encoding/base64"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"strings"
	"time"
)


func RegisterRoutes(e *echo.Echo) {
	
	//文件目录
	
	e.GET("/api/file/:path", func(ctx echo.Context) error {
		path,_ := url.QueryUnescape(ctx.Param("path"))
		/*cookie_, err := ctx.Cookie("email")
		var email string
		if err != nil {
			return err
		} else {
			email = cookie_.Value
		}*/
		email := "admin@godcloud.com"
		return ctx.JSON(http.StatusOK,GetData(email,path))
	})

	//文件获取下载
	
	e.GET("/api/file/:filename/:base64path", func(ctx echo.Context) error {
		filename,_ := url.QueryUnescape(ctx.Param("filename"))
		base64path := ctx.Param("base64path")
		email := "admin@godcloud.com"
		path,_ := base64.StdEncoding.DecodeString(base64path)

		if string(path)=="/"{
			return ctx.Attachment(root+email+"/"+filename,filename)
		}
		return ctx.Attachment(root+email+string(path)+"/"+filename,filename)
	})
	
	//文件上传 / 文件夹创建
	
	e.POST("/api/file/", func(ctx echo.Context) error {
		form, err := ctx.MultipartForm()

		path := ctx.FormValue("path")
		path,_ = url.QueryUnescape(path)
		email := "admin@godcloud.com"
		if err != nil {
			return err
		}

		files := form.File["file"]

		var SavePath string

		if path == "/"{
			SavePath = root+email+"/"
		}else{
			SavePath = root+email+path+"/"
		}
		_ = os.Mkdir(SavePath, 0777)

		if form.File["file"] == nil {
			var tmp PathData
			NameList := strings.Split(ctx.FormValue("name"),"/")
			for i,name_ := range NameList {
				name := strings.TrimSpace(name_)
				if name == "." || name=="" {
					continue
				}
				tmp.DataName = name_
				tmp.DataPath = path
				tmp.DataType = "dir"
				tmp.DataFileId = md5_(name_+path+time.Now().String())
				for _,name__ := range NameList[0:i] {
					if tmp.DataPath == "/" {
						if name__ == "." {
							continue
						}
						tmp.DataPath = ""
					}
					if name__ == "." {
						continue
					}
					tmp.DataPath += "/"+name__
				}
				if tmp.DataPath == "/"{
					SavePath = root+email+"/"+name_
				}else{
					SavePath = root+email+tmp.DataPath+"/"+name_
				}
				_ = os.Mkdir(SavePath, 0777)
				AddData(tmp)
			}
		}else {

			for _, file := range files {
				// Source
				src, err := file.Open()
				if err != nil {
					return err
				}
				defer src.Close()


				// Destination
				dst, err := os.Create(SavePath+file.Filename)
				if err != nil {
					return err
				}
				defer dst.Close()

				// Copy
				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				DataFileId := md5_(file.Filename+path+string(time.Now().Unix()))

				var tmp PathData
				tmp.DataSize = int(file.Size/1024)
				tmp.DataName = file.Filename
				tmp.DataFileId = DataFileId
				tmp.DataType = "file"
				tmp.DataPath = path

				
				video := []string{".wav", ".avi", ".mov", ".mp4", ".flv", ".rmvb", ".mpeg", ".mpg"}
				for _,v := range video {
					if v == path2.Ext(file.Filename) {
						go videothumbnail(SavePath+file.Filename,DataFileId)
					}
				}
				AddData(tmp)

			}
		}


		return ctx.String(http.StatusOK,"succeed")
	})

	//重命名文件(夹)

	e.PUT("/api/file/", func(ctx echo.Context) error {
		tmp := new(UpdateType)
		_ = ctx.Bind(tmp)
		if UpdateFile(tmp) == "succeed" {
			UpdateFiles(tmp)
			return ctx.JSON(200,"succeed")
		}
		return ctx.JSON(200,"failed")
	})

	//移动文件 文件夹
	e.PUT("/api/move/", func(ctx echo.Context) error {
		tmp := new(MoveStruct)
		_ = ctx.Bind(tmp)
		tmp.Email="admin@godcloud.com"
		if MoveFiles(tmp) != "succeed" {
			return ctx.JSON(200,"failed")
		}
		MoveData(tmp)
		return ctx.JSON(200,"succeed")

	})

	//删除文件 文件夹

	e.DELETE("/api/file/", func(ctx echo.Context) error {
		tmp := new(DeleteFile)
		_ = ctx.Bind(tmp)
		if DeleteFiles(tmp) == "succeed" {
			DeleteData(tmp)
			return ctx.JSON(200,"succeed")
		}
		return ctx.JSON(200,"failed")
	})


	//用户登录

	e.GET("/api/login/:email/:pw", func(ctx echo.Context) error {
		email := ctx.Param("email")
		pw := ctx.Param("pw")
		user,status := GetUserInfo(email,pw)
		if status==0 {
			return ctx.JSON(http.StatusOK,"failed")
		}
		return ctx.JSON(http.StatusOK,user)
	})

	//解压缩文件

	e.PUT("/api/UnArchive", func(ctx echo.Context) error {
		tmp := new(UnArchiveFile)
		_ = ctx.Bind(tmp)
		if !UnArchiver(tmp) {
			return ctx.JSON(200,"failed")
		}
		return ctx.JSON(200,"succeed")
	})



}
