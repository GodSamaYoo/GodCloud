package main

import (
	"context"
	"fmt"
	"github.com/zyxar/argo/rpc"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)
//下载服务启动
func aria2begin() rpc.Client {
	aria2enable := ReadIni("aria2", "enable")
	if aria2enable == "no" {
		return nil
	}
	var err error
	if ReadIni("aria2","tmpdownpath") != "" {
		tmppath = ReadIni("aria2","tmpdownpath")
	} else {
		tmppath,err = os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		tmppath += "\\tmp"
	}
	aria2token := ReadIni("aria2","token")
	port := ReadIni("aria2","port")
	aria2url := "http://127.0.0.1:"+port+"/jsonrpc"
	ctx := context.Background()
	var notifier rpc.Notifier = DummyNotifier{}
	t, _ := time.ParseDuration("9999h")
	client, err := rpc.New(ctx, aria2url, aria2token, t, notifier)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return client
}

type DummyNotifier struct{}
func (DummyNotifier) OnDownloadStart(events []rpc.Event)      { fmt.Printf("%s started.", events) }
func (DummyNotifier) OnDownloadPause(events []rpc.Event)      { fmt.Printf("%s paused.", events) }
func (DummyNotifier) OnDownloadStop(events []rpc.Event)       { fmt.Printf("%s stopped.", events) }
func (DummyNotifier) OnDownloadComplete(events []rpc.Event) {
	for _,v := range events {
		infos,_ := aria2client.TellStatus(v.Gid)
		if len(infos.FollowedBy) == 0 {
			a := strings.Split(infos.Files[0].Path,"/") [len(strings.Split(infos.Dir,"\\"))-1:len(strings.Split(infos.Dir,"\\"))+1]
			tmp_ := Aria2DataByAdd(infos.Dir)
			if IsOneDriveStore(tmp_.UserEmail) {
				for _,vv := range infos.Files {
					s,_ := os.Stat(vv.Path)
					ss := vv.Path[len(infos.Dir):]
					ppath := path.Dir(ss)
					if ppath != "/" && ppath != "."{
						var pppath string
						if tmp_.Path == "/" {
							pppath = path.Dir(ppath)
						} else {
							if path.Dir(ppath) == "/"{
								pppath = tmp_.Path
							} else {
								pppath = tmp_.Path + path.Dir(ppath)
							}
						}
						t1 := PathData{
							DataFileId: md5_(path.Base(ppath) + pppath + time.Now().String()),
							DataName:   path.Base(ppath),
							DataType:   "dir",
							DataPath:   pppath,
						}
						AddData(t1,tmp_.UserEmail)
					}
					q := GetOneDriveAdd(tmp_.UserEmail,tmp_.Path,s.Name(),int(s.Size()))
					itemid := Aria2OneDriveUp(infos.Files[0].Path,int(s.Size()),q.Address,q.StoreID)
					picture := []string{"jpg", "jpeg", "bmp", "gif", "png", "tif"}
					tp := path.Ext(s.Name())
					url_ := ""
					if tp != "" {
						for _,v_ := range picture{
							if tp[1:] == v_{
								url_ = GetThumbnail(itemid,q.StoreID)
								break
							}
						}
					}
					fileid := md5_(s.Name() + tmp_.Path + time.Now().String())
					if url_ != "" {
						os.Mkdir("Thumbnail",0777)
						file, _ := http.Get(url_)
						defer file.Body.Close()
						files,_ := ioutil.ReadAll(file.Body)
						_ = ioutil.WriteFile("Thumbnail/"+fileid+".jpg", files, 0644)
					}
					var jjj string
					if tmp_.Path == "/" {
						jjj = ppath
					}else {
						if ppath != "/" {
							jjj = tmp_.Path +ppath
						}else {
							jjj = tmp_.Path
						}
					}
					DatasAdd(&Datas{
						FileID:    fileid,
						Name:      s.Name(),
						Type:      "file",
						Path:      jjj,
						UserEmail: tmp_.UserEmail,
						Size:      int(s.Size()/1024),
						StoreID:   q.StoreID,
						ItemID:    itemid,
					})
					_ = os.RemoveAll(infos.Dir)
				}
			} else {
				add := AddChange(tmp_.Path,tmp_.UserEmail)
				_ = os.Mkdir(add,0777)
				tmp := filepath.Join(a...)
				_ = os.Rename(tmppath+"/"+tmp,add+a[1])
				_ = os.RemoveAll(infos.Dir)
				s,_ :=os.Stat(add+a[1])
				if  s.IsDir() {
					i := 0
					_ = filepath.Walk(add+a[1], func(path string, info os.FileInfo, err error) error {
						var dir PathData
						if i==0 {
							dir.DataPath = tmp_.Path
						} else {
							if tmp_.Path != "/" {
								dir.DataPath = tmp_.Path +"/"+ filepath.ToSlash(filepath.Dir(path[len(add)-2:]))
							} else {
								dir.DataPath = tmp_.Path + filepath.ToSlash(filepath.Dir(path[len(add)-2:]))
							}
						}
						if info.IsDir() {
							dir.DataType = "dir"
						}else {
							dir.DataType = "file"
							dir.DataSize = int(info.Size()/1024)
						}
						dir.DataName = info.Name()
						dir.DataFileId = md5_(s.Name() + dir.DataPath + time.Now().String())
						AddData(dir,tmp_.UserEmail)
						i++
						return err
					})
				}else {
					p := PathData{
						DataFileId: md5_(s.Name() + tmp_.Path + time.Now().String()),
						DataName:   s.Name(),
						DataType:   "file",
						DataPath:   tmp_.Path,
						DataSize:   int(s.Size()/1024),
					}
					AddData(p,tmp_.UserEmail)
				}
			}


		}
	}
	fmt.Printf("%s completed.", events)
}
func (DummyNotifier) OnDownloadError(events []rpc.Event)      { fmt.Printf("%s error.", events) }
func (DummyNotifier) OnBtDownloadComplete(events []rpc.Event) { fmt.Printf("bt %s completed.", events) }

//下载开始

func aria2download(url []string,path string,email string) []string {
	var gids []string
	for _, v := range url {
		var url_ string
		var gid string
		var err error
		tmp := tmppath + "\\" + md5_(time.Now().String())
		if len(v) == 40 && !strings.Contains(v, ".") {
			url_ = "magnet:?xt=urn:btih:" + v
			gid, err = aria2client.AddURI([]string{url_},rpc.Option{"dir":tmp})
		} else if strings.Contains(strings.ToLower(v), "magnet:?xt=urn:btih:") {
			url_ = v
			gid, err = aria2client.AddURI([]string{url_},rpc.Option{"dir":tmp})
		} else {
			gid, err = aria2client.AddURI([]string{v},rpc.Option{"dir":tmp})
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		Aria2DataBegin(gid,path,email,tmp)
		gids = append(gids, gid)
	}
	return gids
}

//下载进度查询

func aria2status(email string) []aria2downloadinfo {
	var infos []aria2downloadinfo
	task := Aria2DataInfo(email)
	for _,v := range task {
		totalinfo,_ := aria2client.TellStatus(v.Gid)
		if len(totalinfo.FollowedBy)>0 && totalinfo.Status == "complete" {
			totalinfo,_ = aria2client.TellStatus(totalinfo.FollowedBy[0])
		}
		var info aria2downloadinfo
		info.Gid = totalinfo.Gid
		info.TotalLength = totalinfo.TotalLength
		info.FileNums = len(totalinfo.Files)
		info.DownloadSpeed = totalinfo.DownloadSpeed
		info.Status = totalinfo.Status
		info.CompletedLength = totalinfo.CompletedLength
		info.BeginTime = v.Time
		info.Path = v.Path
		info.Files = totalinfo.Files
		for i,_ := range info.Files {
			tmp := strings.Split(info.Files[i].Path,"/")
			info.Files[i].Path = tmp[len(tmp)-1]
		}

		infos = append(infos, info)
	}
	return infos
}

//下载任务 暂停 取消 继续 

func aria2taskchange(tmp *aria2change) error {
	var err error
	if tmp.Type == 1 {
		_, err = aria2client.ForcePause(tmp.Gid)
		if err != nil {
			return err
		}
	}else if tmp.Type == 2 {
		_, err = aria2client.ForceRemove(tmp.Gid)
	} else if tmp.Type == 3 {
		_, err = aria2client.Unpause(tmp.Gid)
	} else if tmp.Type == 4 {
		_, err = aria2client.ForceRemove(tmp.Gid)
		totalinfo,_ := aria2client.TellStatus(tmp.Gid)
		err = os.RemoveAll(totalinfo.Dir)
		Aria2DeleteTask(totalinfo.Dir)
	} else if tmp.Type == 5 {
		totalinfo,_ := aria2client.TellStatus(tmp.Gid)
		err = os.RemoveAll(totalinfo.Dir)
		Aria2DeleteTask(totalinfo.Dir)
	}
	if err != nil {
		return err
	}
	return nil
}



