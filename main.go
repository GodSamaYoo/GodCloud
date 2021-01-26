package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/zyxar/argo/rpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var root string
var aria2client rpc.Client

func main() {
	e := echo.New()
	_, err := os.Stat("data.db")
	if err != nil {
		db, err_ := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
			SkipDefaultTransaction: true,
		})

		if err_ != nil {
			fmt.Println("创建Sqlite数据库失败")
		}
		_ = db.AutoMigrate(&Users{})
		_ = db.AutoMigrate(&Datas{})
		_ = db.AutoMigrate(&Usergroups{})
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
		})
		fmt.Println("用户名：admin@godcloud.com\n密码：123456")
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 3,
	}))
	e.Use(middleware.CORS())
	e.Static("/", "public")
	RegisterRoutes(e)
	root = ReadIni(
		"filepath",
		"path",
	)
	aria2enable := ReadIni("aria2", "enable")
	if aria2enable != "no" {
		aria2url := ReadIni("aria2", "url")
		aria2token := ReadIni("aria2", "token")
		aria2client = aria2begin(aria2url, aria2token)
	}
	e.HideBanner = true
	fmt.Print("\n   __________  ____     ________    ____  __  ______ \n  / ____/ __ \\/ __ \\   / ____/ /   / __ \\/ / / / __ \\\n / / __/ / / / / / /  / /   / /   / / / / / / / / / /\n/ /_/ / /_/ / /_/ /  / /___/ /___/ /_/ / /_/ / /_/ / \n\\____/\\____/_____/   \\____/_____/\\____/\\____/_____/  \n                                                     \n")
	e.Logger.Fatal(e.Start(":2020"))
}
