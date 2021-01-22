package main

import (
	"fmt"
	"github.com/mholt/archiver"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
)

//更新文件名（本机）

func UpdateFiles(tmp *UpdateType){

	if tmp.Path == "/" {
		_ = os.Rename(root+tmp.Email+"/"+tmp.OldName, root+tmp.Email+"/"+tmp.NewName)
	} else {
		_ = os.Rename(root+tmp.Email+tmp.Path+"/"+tmp.OldName, root+tmp.Email+tmp.Path+"/"+tmp.NewName)
	}
	return
}


//删除文件（本机）
func DeleteFiles(files *DeleteFile) string{
	email := "123456@qq.com"
	var DeletePath string
	path := files.Path
	if path == "/"{
		DeletePath = root+email+"/"
	}else{
		DeletePath = root+email+path+"/"
	}
	if files.Type != "dir" {
		err := os.Remove(DeletePath + files.FileName)
		if err != nil {
			fmt.Println(DeletePath + files.FileName+"\n删除失败")
			return "failed"
		}
	} else {
		err := os.RemoveAll(DeletePath + files.FileName)
		if err != nil {
			fmt.Println(DeletePath + files.FileName+"\n删除失败")
			return "failed"
		}
	}
	return "succeed"
}

//移动文件 文件夹
func MoveFiles(tmp *MoveStruct) string {
	OldPath := root+tmp.Email+tmp.OldPath
	NewPath := root+tmp.Email+tmp.NewPath
	if tmp.OldPath == "/"{
		OldPath += tmp.FileName
	}else {
		OldPath += "/"+tmp.FileName
	}
	if tmp.NewPath == "/"{
		NewPath += tmp.FileName
	}else {
		NewPath += "/"+tmp.FileName
	}
	if _, err := os.Stat(NewPath); os.IsNotExist(err) {
		// 当文件不存在, 才执行创建
		err = os.Rename(OldPath, NewPath)
		if err != nil {
			fmt.Println("文件移动失败，路径："+OldPath)
			fmt.Println(err)
			return "failed"
		}
	}
	return "succeed"
}

//解压缩
func UnArchiver(tmp *UnArchiveFile) bool {
	TmpPath := md5_(string(time.Now().Unix()))
	path := root+tmp.Email+tmp.Path
	if tmp.Path == "/"{
		path += tmp.FileName
	}else {
		path += "/"+tmp.FileName
	}
	os.Mkdir("./"+TmpPath,777)
	if tmp.PassWord == "" {
		err := archiver.Unarchive(path,"./"+TmpPath)
		if err != nil {
			fmt.Println("解压失败:")
			fmt.Println(err)
			return false
		}
	} else {
		a := archiver.NewRar()
		a.Password = tmp.PassWord
		a.Unarchive(path,"./"+TmpPath)

	}
	i := 0
	email := "123456@qq.com"
	_ = filepath.Walk("./"+TmpPath, func(path string, info os.FileInfo, err error) error {
		if i != 0 {
			lens := len("./"+TmpPath)
			NowPath := string([]rune(path)[lens-1:])
			PathSlice := strings.Split(NowPath,`\`)
			SavePath := AddChange(tmp.NewPath,email)
			var dir PathData
			if info.IsDir() {
				dir.DataPath = tmp.NewPath
				if len(PathSlice) != 1 {
					dir.DataPath = dir.DataPath[:len(dir.DataPath)-1]
					for q,v := range PathSlice {
						if q == len(PathSlice)-1 {
							continue
						}
						SavePath += "/"+v
						dir.DataPath += "/"+v
					}
				}
				_ = os.Mkdir(SavePath+"/"+info.Name(), 777)
				dir.DataEmail = email
				dir.DataType = "dir"
				dir.DataFileId = md5_(info.Name()+dir.DataPath+string(time.Now().Unix()))
				dir.DataTime = time.Now().Format("2006-01-02 15:04:05")
				dir.DataName = info.Name()
				AddData(dir)
			} else {
				dir.DataPath = tmp.NewPath
				if len(PathSlice) != 1 {
					dir.DataPath = dir.DataPath[:len(dir.DataPath)-1]
					SavePath=SavePath[:len(SavePath)-1]
					for q,v := range PathSlice {
						if q == len(PathSlice)-1 {
							continue
						}
						SavePath += "/"+v
						dir.DataPath += "/"+v
					}
				}
				err = os.Rename(path, SavePath+"/"+info.Name())
				if err != nil {
					fmt.Println("移动失败："+path)
					fmt.Println(err)
					return err
				}
				dir.DataEmail = email
				dir.DataType = "file"
				dir.DataFileId = md5_(info.Name()+dir.DataPath+string(time.Now().Unix()))
				dir.DataTime = time.Now().Format("2006-01-02 15:04:05")
				dir.DataName = info.Name()
				dir.DataSize = int(info.Size()/1024)
				AddData(dir)
			}
		}
		i++
		return nil
	})
	err := os.RemoveAll("./"+TmpPath)
	if err != nil {
		fmt.Println("移除临时文件失败")
		fmt.Println(err)
	}
	return true
}

//网盘地址转本地存储地址

func AddChange(NetPath string,email string) string {
	var SavePath string
	if NetPath == "/"{
		SavePath = root+email+"/"
	}else{
		SavePath = root+email+NetPath+"/"
	}
	return SavePath
}

//视频缩略图生成 100张合成1张

func VideoThumbnails(path string) {
	newImg := image.NewNRGBA(image.Rect(0, 0, 16000, 90))
	i := 0
	filepath.Walk(path, func (path string, info os.FileInfo, err error) error {
		if path != path {
			m := openAndDecode(path)
			draw.Draw(newImg, newImg.Bounds(), m, image.Pt(-160*i,0), draw.Over)
			fmt.Println(i)
		}
		i++
		return nil
	})
	imgfile, _ := os.Create("123.jpg")
	defer imgfile.Close()
	jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
}

func openAndDecode(imgPath string) image.Image {
	img, err := os.Open(imgPath)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to open %v", err))
	}

	decoded, err := jpeg.Decode(img)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to open %v", err))
	}
	defer img.Close()

	m := resize.Resize(0, 90, decoded, resize.Lanczos3)

	return m
}