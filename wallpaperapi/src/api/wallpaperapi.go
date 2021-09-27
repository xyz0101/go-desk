package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}
type Category struct {
	Count        int           `json:"count"`
	Ename        string        `json:"ename"`
	Rname        string        `json:"rname"`
	CoverTemp    string        `json:"cover_temp"`
	Name         string        `json:"name"`
	Cover        string        `json:"cover"`
	Rank         int           `json:"rank"`
	Filter       []interface{} `json:"filter"`
	Sn           int           `json:"sn"`
	Icover       string        `json:"icover"`
	Atime        int           `json:"atime"`
	Type         int           `json:"type"`
	ID           string        `json:"id"`
	PicassoCover string        `json:"picasso_cover"`
}
type Wallpaper struct {
	Views   int           `json:"views"`
	Ncos    int           `json:"ncos"`
	Rank    int           `json:"rank"`
	Tag     []string      `json:"tag"`
	Wp      string        `json:"wp"`
	Xr      bool          `json:"xr"`
	Cr      bool          `json:"cr"`
	Favs    int           `json:"favs"`
	Atime   int           `json:"atime"`
	ID      string        `json:"id"`
	Desc    string        `json:"desc"`
	Thumb   string        `json:"thumb"`
	Img     string        `json:"img"`
	Cid     []string      `json:"cid"`
	URL     []interface{} `json:"url"`
	Preview string        `json:"preview"`
	Store   string        `json:"store"`
}
type Res struct {
	Wallpaper []Wallpaper `json:"wallpaper"`
	Category  []Category  `json:"category"`
}
type Data struct {
	Msg  string `json:"msg"`
	Res  Res    `json:"res"`
	Code int    `json:"code"`
}

const (
	token         = "kTqzFtkEPUufjkgKf9WBYBaWS2raNb7f9rTlPjCKTMtaLir4JFUmLsHk1Etw9thfCpmvFwhlwJKRxXO3DICtKmytzNddaznDIMHQB9z1Hrq2R5PkkOEvHlXgO8fi8j+FGjCYiSNrJIUtBvZ0kD2v5kNF4nJvRuX9/ip64QEgOJN1ijPqAjsBI/sVUnZQUNJRObHcVS8S+uViWyfqrNjbgiY8/HQXxFZF0dj+XwBi8ur8H981Zy3yNaJzpB403j7k"
	host          = "http://tencent.jenkin.tech:8000"
	category      = host + "/lsc/files/aibizhi/getAbzCategory"
	wallpaperList = host + "/lsc/files/aibizhi/getAbzWallpaperWin?category={category}&skip={skip}"
)

//获取分类
func GetCategory() Response {
	response := &Response{}
	GetHttp(category, response)
	return *response
}

//获取壁纸列表
func GetWallpaperList(categoryId string, skip int) Response {
	response := &Response{}
	wallpaperUrl := strings.Replace(wallpaperList, "{category}", categoryId, 1)
	wallpaperUrl = strings.Replace(wallpaperUrl, "{skip}", strconv.Itoa(skip), 1)
	GetHttp(wallpaperUrl, response)
	return *response
}

//发送http请求
func GetHttp(url string, result interface{}) interface{} {
	t := time.Now()
	fmt.Println("请求的路径：", url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("createRequest异常 ", err)
	}
	req.Header.Add("token", token)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body := resp.Body
	err = json.NewDecoder(body).Decode(result)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("耗时：", time.Since(t).Milliseconds(), "毫秒")
	return result
}

func main() {
	res := GetWallpaperList("4e4d610cdf714d2966000000", 1)
	fmt.Println(res)
}
