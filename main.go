package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var root string
func main() {
	e := echo.New()
	_, err := os.Stat("data.db")
	if err != nil {
		db, err_ := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
			SkipDefaultTransaction: true,
		})
		if err_ != nil {
			panic("创建Sqlite数据库失败")
		}
		db.Create(Users{
			UserID:   1,
			Email:    "admin@godcloud.com",
			Password: "123456",
			GroupID:  1,
			Volume:   1048576,
		})
		db.Create(Usergroups{
			GroupID: 1,
			name:    "管理员",
			volume:  1048576,
		})
		fmt.Println("用户名：admin@godcloud.com\n密码：123456")
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 3,
	}))
	e.Use(middleware.CORS())
	e.Static("/","public")
	RegisterRoutes(e)
	root = ReadIni(
		"filepath",
		"path",
	)
	e.HideBanner = true
	fmt.Print("\n   __________  ____     ________    ____  __  ______ \n  / ____/ __ \\/ __ \\   / ____/ /   / __ \\/ / / / __ \\\n / / __/ / / / / / /  / /   / /   / / / / / / / / / /\n/ /_/ / /_/ / /_/ /  / /___/ /___/ /_/ / /_/ / /_/ / \n\\____/\\____/_____/   \\____/_____/\\____/\\____/_____/  \n                                                     \n")
	e.Logger.Fatal(e.Start(":2020"))
}
