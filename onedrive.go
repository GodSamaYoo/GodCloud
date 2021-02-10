package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func FirstGetRefreshToken(tmp *storemodel) error{
	url_ := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	method := "POST"
	payload := strings.NewReader("client_id="+tmp.ClientID+"&code="+tmp.Code+
		"&redirect_uri="+tmp.RedirectUrl+"&grant_type=authorization_code&client_secret="+
		tmp.ClientSecret)
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url_, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var data onedrivererefresh
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	tmp.RefreshToken = data.RefreshToken
	if tmp.RefreshToken != "" {
		if StoreAdd(tmp) {
			return nil
		}
	}
	return nil
}

//刷新token
func RefreshToken(id int) error {
	tmp := StoreInfoQuery(&Store{
		ID: id,
	})
	url := "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	method := "POST"

	payload := strings.NewReader("client_id="+tmp.ClientID+"&grant_type=refresh_token&client_secret="+tmp.ClientSecret+
		"&refresh_token="+tmp.RefreshToken)

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var data onedrivererefresh
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	onedrivetokens[id] = data.Token
	tmp_ := new(Store)
	tmp_.ID = id
	tmp_.RefreshToken = data.RefreshToken
	OneDriveRefreshUpdate(tmp_)

	return nil
}

//配合刷新
func RefreshToken2() {
	var tmp []Store
	db.Where(&Store{
		Type: "onedrive",
	}).Find(&tmp)
	for _,v := range tmp {
		i :=0
		for {
			err := RefreshToken(v.ID)
			if err == nil || i > 4 {
				fmt.Println("OneDrive刷新成功")
				break
			}
			i++
		}
	}
}
//获得上传地址

func UploadAddress (token, path1,path2,email, filename string) string {
	var url string
	if path2 != "/" {
		url = "https://graph.microsoft.com/v1.0/me/drive/root:"+path1+"/"+email+path2+"/"+filename+":/createUploadSession"
	} else {
		url = "https://graph.microsoft.com/v1.0/me/drive/root:"+path1+"/"+email+"/"+filename+":/createUploadSession"
	}
	method := "POST"
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}
	var data onedriveAddress
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ""
	}
	return data.UploadAddress
}
func GetOneDriveAdd(email, path, filename string,Need int) *OnedriveAddReturn {
	ttt := new(OnedriveAddReturn)
	Need = Need/1024 + 1
	a,_ := GetUserInfo(email)
	if Need > (a.Volume - a.Used - 10240) {
		return ttt
	}
	b := UserInfoQuery(&Users{
		Email: email,
	})
	c := UserGroupQuery(&Usergroups{
		GroupID: b.GroupID,
	})
	d := strings.Split(c.StoreID,",")
	f := new(Store)
	var e int
	for i,v := range d {
		e,_ =strconv.Atoi(v)
		f = StoreInfoQuery(&Store{
			ID: e,
		})
		if Need > (f.Volume - f.Used - 10240) {
			if i == len(d) {
				return ttt
			}
			continue
		}
		break
	}
	token := onedrivetokens[e]
	g := UploadAddress(token,f.Path,path,email,filename)
	ttt.StoreID = e
	ttt.Address = g
	return ttt
}

//获取缩略图
func GetThumbnail(itemid string,onedriveid int) string {
	url := "https://graph.microsoft.com/v1.0/me/drive/items/"+itemid+"/thumbnails?select=large"
	method := "GET"
	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+onedrivetokens[onedriveid])
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	url_ := gjson.Get(string(body),"value.#(large).large.url").Str
	return url_
}

func Aria2OneDriveUp(filepath string,size int,url string,storeid int) string {
	f,_ := os.Open(filepath)
	defer f.Close()
	buf := make([]byte,1024 * 320 * 80)
	n := 0
	var status string
	for {
		num,err := f.Read(buf)
		n++
		if err==io.EOF{
			break
		}
		if num < 1024 * 320 * 80 {
			buf_ := make([]byte,num)
			buf_ = buf [:num]
			reader := bytes.NewReader(buf_)
			status = OneDriveUp(url,reader,(n-1)*1024 * 320 * 80,(n-1)*1024 * 320 * 80+num-1,size,storeid)
			if status == "-1" {
				i := 0
				for {
					status = OneDriveUp(url,reader,(n-1)*1024 * 320 * 80,(n-1)*1024 * 320 * 80+num-1,size,storeid)
					if status != "-1" || i==5 {
						break
					}
					i++
				}
			}
		} else {
			reader := bytes.NewReader(buf)
			status = OneDriveUp(url,reader,(n-1)*1024 * 320 * 80,n*1024 * 320 * 80-1,size,storeid)
			if status == "-1" {
				i := 0
				for {
					status = OneDriveUp(url,reader,(n-1)*1024 * 320 * 80,n*1024 * 320 * 80-1,size,storeid)
					if status != "-1" || i==5 {
						break
					}
					i++
				}
			}
		}
	}
	return gjson.Get(status,"id").Str
}

func DeleteOneDriveFile(itemid string,onedriveid int) string {
	url := "https://graph.microsoft.com/v1.0/me/drive/items/"+itemid
	method := "DELETE"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+onedrivetokens[onedriveid])
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	return "succeed"
}

func OneDriveUp(url string,chunk io.Reader,begin,end,filesize int,storeid int) string {
	method := "PUT"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, chunk)

	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	begin_ := strconv.Itoa(begin)
	end_ := strconv.Itoa(end)
	filesize_ := strconv.Itoa(filesize)
	req.Header.Add("Content-Range", "bytes "+begin_+"-"+end_+"/"+filesize_)
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Authorization", "Bearer "+onedrivetokens[storeid])

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	defer res.Body.Close()

	if res.StatusCode == 202 {
		return "1"
	} else if res.StatusCode == 201 {
		body, _ := ioutil.ReadAll(res.Body)
		return string(body)
	} else {
		return "-1"
	}

}

