package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zyxar/argo/rpc"
	"gorm.io/gorm"
	"github.com/robfig/cron/v3"
)

var root string
var aria2client rpc.Client
var option  map[string]string
var deskey string
var db gorm.DB
var tmppath string
var onedrivetokens map[int]string

func main() {
	e := echo.New()
	CheckSqlite()
	CheckIni()
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 3}))
	e.Use(middleware.CORS())
	e.Static("/", "public")
	RegisterRoutes(e)
	root = ReadIni("filepath", "path")
	ServicePort := ReadIni("service", "port")
	deskey = ReadIni("Des","key")
	aria2client = aria2begin()
	c := cron.New()
	onedrivetokens = make(map[int]string)
	RefreshToken2()
	_, _ = c.AddFunc("*/50 * * * *", RefreshToken2)
	c.Start()
	e.HideBanner = true
	fmt.Print("\n   __________  ____     ________    ____  __  ______ \n  / ____/ __ \\/ __ \\   / ____/ /   / __ \\/ / / / __ \\\n / / __/ / / / / / /  / /   / /   / / / / / / / / / /\n/ /_/ / /_/ / /_/ /  / /___/ /___/ /_/ / /_/ / /_/ / \n\\____/\\____/_____/   \\____/_____/\\____/\\____/_____/  \n                                                     \n")
	e.Logger.Fatal(e.Start(":" + ServicePort))
}
