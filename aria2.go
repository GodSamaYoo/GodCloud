package main

import (
	"context"
	"fmt"
	"github.com/zyxar/argo/rpc"
	"strings"
	"time"
)
//下载服务启动
func aria2begin(aria2url string, aria2token string) rpc.Client {
	ctx := context.Background()
	var notifier rpc.Notifier
	t, _ := time.ParseDuration("9999h")
	client, err := rpc.New(ctx, aria2url, aria2token, t, notifier)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return client
}

//下载开始

func aria2download(url []string,path string,email string) []string {
	var gids []string
	for _, v := range url {
		fmt.Println(v)
		var url_ string
		var gid string
		var err error
		if len(v) == 40 && !strings.Contains(v, ".") {
			url_ = "magnet:?xt=urn:btih:" + v
			gid, err = aria2client.AddURI([]string{url_})
		} else if strings.Contains(strings.ToLower(v), "magnet:?xt=urn:btih:") {
			url_ = v
			gid, err = aria2client.AddURI([]string{url_})
		} else {
			gid, err = aria2client.AddURI([]string{v})
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		Aria2DataBegin(gid,path,email)
		gids = append(gids, gid)
	}
	return gids
}

//下载进度查询



