package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/mholt/archiver"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//更新文件名（本机）

func UpdateFiles(tmp *UpdateType) {

	if tmp.Path == "/" {
		_ = os.Rename(root+tmp.Email+"/"+tmp.OldName, root+tmp.Email+"/"+tmp.NewName)
	} else {
		_ = os.Rename(root+tmp.Email+tmp.Path+"/"+tmp.OldName, root+tmp.Email+tmp.Path+"/"+tmp.NewName)
	}
	return
}

//删除文件（本机）
func DeleteFiles(files *DeleteFile,email string) string {
	var DeletePath string
	path := files.Path
	if path == "/" {
		DeletePath = root + email + "/"
	} else {
		DeletePath = root + email + path + "/"
	}
	if files.Type != "dir" {
		err := os.Remove(DeletePath + files.FileName)
		if err != nil {
			fmt.Println(DeletePath + files.FileName + "\n删除失败")
			return "failed"
		}
	} else {
		err := os.RemoveAll(DeletePath + files.FileName)
		if err != nil {
			fmt.Println(DeletePath + files.FileName + "\n删除失败")
			return "failed"
		}
	}
	return "succeed"
}

//移动文件 文件夹
func MoveFiles(tmp *MoveStruct) string {
	OldPath := root + tmp.Email + tmp.OldPath
	NewPath := root + tmp.Email + tmp.NewPath
	if tmp.OldPath == "/" {
		OldPath += tmp.FileName
	} else {
		OldPath += "/" + tmp.FileName
	}
	if tmp.NewPath == "/" {
		NewPath += tmp.FileName
	} else {
		NewPath += "/" + tmp.FileName
	}
	if _, err := os.Stat(NewPath); os.IsNotExist(err) {
		// 当文件不存在, 才执行创建
		err = os.Rename(OldPath, NewPath)
		if err != nil {
			fmt.Println("文件移动失败，路径：" + OldPath)
			fmt.Println(err)
			return "failed"
		}
	}
	return "succeed"
}

//解压缩
func UnArchiver(tmp *UnArchiveFile,email string) bool {
	TmpPath := md5_(time.Now().String())
	path := root + tmp.Email + tmp.Path
	if tmp.Path == "/" {
		path += tmp.FileName
	} else {
		path += "/" + tmp.FileName
	}
	_ = os.Mkdir("./"+TmpPath, 777)
	if tmp.PassWord == "" {
		err := archiver.Unarchive(path, "./"+TmpPath)
		if err != nil {
			fmt.Println("解压失败:")
			fmt.Println(err)
			return false
		}
	} else {
		a := archiver.NewRar()
		a.Password = tmp.PassWord
		_ = a.Unarchive(path, "./"+TmpPath)

	}
	i := 0
	_ = filepath.Walk("./"+TmpPath, func(path string, info os.FileInfo, err error) error {
		if i != 0 {
			lens := len("./" + TmpPath)
			NowPath := string([]rune(path)[lens-1:])
			PathSlice := strings.Split(NowPath, `\`)
			SavePath := AddChange(tmp.NewPath, email)
			var dir PathData
			if info.IsDir() {
				dir.DataPath = tmp.NewPath
				if len(PathSlice) != 1 {
					dir.DataPath = dir.DataPath[:len(dir.DataPath)-1]
					for q, v := range PathSlice {
						if q == len(PathSlice)-1 {
							continue
						}
						SavePath += "/" + v
						dir.DataPath += "/" + v
					}
				}
				_ = os.Mkdir(SavePath+"/"+info.Name(), 777)
				dir.DataType = "dir"
				dir.DataFileId = md5_(info.Name() + dir.DataPath + time.Now().String())
				dir.DataName = info.Name()
				AddData(dir,email)
			} else {
				dir.DataPath = tmp.NewPath
				if len(PathSlice) != 1 {
					dir.DataPath = dir.DataPath[:len(dir.DataPath)-1]
					SavePath = SavePath[:len(SavePath)-1]
					for q, v := range PathSlice {
						if q == len(PathSlice)-1 {
							continue
						}
						SavePath += "/" + v
						dir.DataPath += "/" + v
					}
				}
				err = os.Rename(path, SavePath+"/"+info.Name())
				if err != nil {
					fmt.Println("移动失败：" + path)
					fmt.Println(err)
					return err
				}
				dir.DataType = "file"
				dir.DataFileId = md5_(info.Name() + dir.DataPath + time.Now().String())
				dir.DataName = info.Name()
				dir.DataSize = int(info.Size() / 1024)
				AddData(dir,email)
			}
		}
		i++
		return nil
	})
	err := os.RemoveAll("./" + TmpPath)
	if err != nil {
		fmt.Println("移除临时文件失败")
		fmt.Println(err)
	}
	return true
}

//网盘地址转本地存储地址

func AddChange(NetPath string, email string) string {
	var SavePath string
	if NetPath == "/" {
		SavePath = root + email + "/"
	} else {
		SavePath = root + email + NetPath + "/"
	}
	return SavePath
}

//文件上传  文件夹创建处理

func FileUploadAndCreate(ctx echo.Context) bool {
	form, _ := ctx.MultipartForm()
	path := ctx.FormValue("path")
	path, _ = url.QueryUnescape(path)
	files := form.File["file"]
	enkey,err := ctx.Cookie("GODKEY")
	if err != nil {
		return false
	}
	email,_ := DesDecrypt(enkey.Value,deskey)
	var SavePath string

	if path == "/" {
		SavePath = root + email + "/"
	} else {
		SavePath = root + email + path + "/"
	}
	_ = os.Mkdir(SavePath, 0777)

	if form.File["file"] == nil {
		var tmp PathData
		NameList := strings.Split(ctx.FormValue("name"), "/")
		for i, name_ := range NameList {
			name := strings.TrimSpace(name_)
			if name == "." || name == "" {
				continue
			}
			tmp.DataName = name_
			tmp.DataPath = path
			tmp.DataType = "dir"
			tmp.DataFileId = md5_(name_ + path + time.Now().String())
			for _, name__ := range NameList[0:i] {
				if tmp.DataPath == "/" {
					if name__ == "." {
						continue
					}
					tmp.DataPath = ""
				}
				if name__ == "." {
					continue
				}
				tmp.DataPath += "/" + name__
			}
			if tmp.DataPath == "/" {
				SavePath = root + email + "/" + name_
			} else {
				SavePath = root + email + tmp.DataPath + "/" + name_
			}
			if AddData(tmp,email) {
				_ = os.Mkdir(SavePath, 0777)
			} else {
				return false
			}

		}
	} else {
		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				return false
			}

			// Destination
			dst, err := os.Create(SavePath + file.Filename)
			if err != nil {
				return false
			}

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return false
			}

			DataFileId := md5_(file.Filename + path + time.Now().String())

			var tmp PathData
			tmp.DataSize = int(file.Size / 1024)
			tmp.DataName = file.Filename
			tmp.DataFileId = DataFileId
			tmp.DataType = "file"
			tmp.DataPath = path

			/*
				video := []string{".wav", ".avi", ".mov", ".mp4", ".flv", ".rmvb", ".mpeg", ".mpg"}
				for _,v := range video {
					if v == path2.Ext(file.Filename) {
						go videothumbnail(SavePath+file.Filename,DataFileId)
					}
				}*/

			err = src.Close()
			err = dst.Close()
			if err != nil {
				return false
			}
			AddData(tmp,email)
		}
	}
	return true
}

//用户容量是否足够

func UserVolumeIsOK(email string,Need int) bool {
	a,_ := GetUserInfo(email)
	if Need < (a.Volume - a.Used - 10240) {
		return true
	}
	return false
}

//储存策略容量是否足够
func StoreVolumeIsOK(id int,Need int) bool {
	tmp := StoreInfoQuery(&Store{
		ID: id,
	})
	if Need < (tmp.Volume - tmp.Used - 10240) {
		return true
	}
	return false
}



/* cgo编译存在问题

//视频缩略图生成 100张合成1张

func videothumbnail(videopath string, imgname string) {
	a, _ := screengen.NewGenerator(videopath)
	cuttime := int(a.Duration / 101)
	newImg := image.NewNRGBA(image.Rect(0, 0, 16000, 90))
	for i := 1; i < 101; i++ {
		tmpimg, _ := a.ImageWxH(int64(cuttime*i), 160, 90)
		draw.Draw(newImg, newImg.Bounds(), tmpimg, image.Pt(-160*i, 0), draw.Over)
	}
	_ = os.Mkdir("./thumbnail", 0777)
	imgfile, _ := os.Create("./thumbnail/"+imgname + ".jpg")
	defer imgfile.Close()
	jpeg.Encode(imgfile, newImg, &jpeg.Options{100})
}

*/
