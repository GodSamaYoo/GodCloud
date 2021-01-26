package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

const IniPath = "./GodCloud.ini"

func CheckIni() {
	_, err := os.Stat(IniPath)
	if err != nil || os.IsNotExist(err) {
		CreateIni(IniPath)
		cfg := OpenIni(IniPath)
		if cfg != nil {
			_, _ = cfg.Section("filepath").NewKey("path", "./")
			_, _ = cfg.Section("aria2").NewKey("enable", "no")
			_, _ = cfg.Section("service").NewKey("port", "2020")
			err = cfg.SaveTo(IniPath)
		}
	}
}

func OpenIni(path string) *ini.File {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Println("配置读取失败")
	} else {
		return cfg
	}
	return nil
}

func CreateIni(path string) {
	_, err := os.Create(path)
	if err != nil {
		fmt.Println("创建失败")
	}
}

func ReadIni(Section string, Key string) string {
	cfg := OpenIni("./GodCloud.ini")
	text := cfg.Section(Section).Key(Key).String()
	return text
}
