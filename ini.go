package main

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

const IniPath = "./config.ini"
/*func main(){
	_, err := os.Stat(IniPath)
	if err != nil || os.IsNotExist(err){
		CreateIni(IniPath)
		cfg := OpenIni(IniPath)
		if cfg != nil {
			_, _ = cfg.Section("test").NewKey("name", "张三")
			err = cfg.SaveTo(IniPath)
			fmt.Println("配置写入成功")
		}
	}else{
		cfg := OpenIni(IniPath)
		if cfg != nil {
			fmt.Println(cfg.Section("test").Key("name"))
		}
	}
}*/

func OpenIni(path string) *ini.File {
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Println("配置读取失败")
	}else {
		return cfg
	}
	return nil
}

func CreateIni(path string)  *os.File{
	cfg, err := os.Create(path)
	if err != nil {
		fmt.Println("创建失败")
	}else {
		return cfg
	}
	return nil
}

func ReadIni(Section string,Key string) string {
	cfg := OpenIni("./config.ini")
	text := cfg.Section(Section).Key(Key).String()
	return text
}
