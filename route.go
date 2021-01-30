package main

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"net/url"
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

	//文件上传 / 文件夹创建

	e.POST("/api/file/", func(ctx echo.Context) error {

		err := FileUploadAndCreate(ctx)
		if err != nil {
			fmt.Println(err)
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
		if DeleteFiles(tmp,email) == "succeed" {
			DeleteData(tmp,email)
			return ctx.JSON(200, "succeed")
		}
		return ctx.JSON(200, "failed")
	})

	//用户登录

	e.POST("/api/login", func(ctx echo.Context) error {
		email := ctx.FormValue("email")
		pw := ctx.FormValue("pw")
		user, status := GetUserInfo(email, pw)
		if status == 0 {
			return ctx.JSON(http.StatusOK, "failed")
		}
		cookie := new (http.Cookie)
		cookie.Name = "GODKEY"
		cookie.Value,_ = DesEncrypt(email,deskey)
		ctx.SetCookie(cookie)
		return ctx.JSON(http.StatusOK, user)
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
}
