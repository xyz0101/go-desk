package wallpaper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

const (
	UserAgent      = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36"
	BingHomeURL    = "https://cn.bing.com"
	CurrentPathDir = "cache/"
)

const (
	Size1k string = "1920,1080"
	Size2k string = "2560,1440"
	Size4k string = "3840,2160"
)

// ImageSize 图片大小
type ImageSize struct {
	w string
	h string
}

func init() {
	_ = os.Mkdir(CurrentPathDir, 0755)
}

// EncodeMD5 MD5编码
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// SetWindowsWallpaper 设置windows壁纸
func setWindowsWallpaper(imagePath string) error {
	dll := syscall.NewLazyDLL("user32.dll")
	proc := dll.NewProc("SystemParametersInfoW")
	_t, _ := syscall.UTF16PtrFromString(imagePath)
	ret, _, _ := proc.Call(20, 1, uintptr(unsafe.Pointer(_t)), 0x1|0x2)
	if ret != 1 {
		return errors.New("系统调用失败")
	}
	return nil
}

// DownloadImage 下载图片,保存并返回保存的文件名的绝对路径
func DownloadImage(imageURL string) (string, error) {

	client := http.Client{}

	request, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	day := time.Now().Format("2006-01-02")

	fileName := EncodeMD5(imageURL)
	path := CurrentPathDir + fmt.Sprintf("[%s]%s", day, fileName) + ".jpg"

	err = ioutil.WriteFile(path, body, 0755)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func SetWallpaper(imageURL string) {
	//imageURL :="http://img.aibizhi.adesk.com/614acf3be7bce72b931d3d2f?sign=497eb83aa7c6c6783cafdc9d10ed65f3&t=61507c5e"

	fmt.Println("下载图片...", imageURL)
	imagePath, err := DownloadImage(imageURL)
	if err != nil {
		fmt.Println("下载图片失败: " + err.Error())
		return
	}
	fmt.Println("设置桌面...")
	err = setWindowsWallpaper(imagePath)
	if err != nil {
		fmt.Println("设置桌面背景失败: " + err.Error())
		return
	}
	fmt.Println("桌面设置成功")
}
