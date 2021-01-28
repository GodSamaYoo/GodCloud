package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zyxar/argo/rpc"
)

var root string
var aria2client rpc.Client
var option  map[string]string
var deskey string

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
	aria2enable := ReadIni("aria2", "enable")
	deskey = ReadIni("Des","key")
	if aria2enable == "yes" {
		aria2url := ReadIni("aria2", "url")
		aria2token := ReadIni("aria2", "token")
		/*dir,_ := os.Getwd()
		option["dir"] = dir+"/tmp"*/
		aria2client = aria2begin(aria2url, aria2token)
	}
	e.HideBanner = true
	fmt.Print("\n   __________  ____     ________    ____  __  ______ \n  / ____/ __ \\/ __ \\   / ____/ /   / __ \\/ / / / __ \\\n / / __/ / / / / / /  / /   / /   / / / / / / / / / /\n/ /_/ / /_/ / /_/ /  / /___/ /___/ /_/ / /_/ / /_/ / \n\\____/\\____/_____/   \\____/_____/\\____/\\____/_____/  \n                                                     \n")
	e.Logger.Fatal(e.Start(":" + ServicePort))
}
