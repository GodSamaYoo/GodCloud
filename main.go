package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var root string
func main() {
	e := echo.New()
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
