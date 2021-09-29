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
	"unsafe"
)

const (
	CurrentPathDir = "cache/"
)

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

func setSleepWallpaper(imagePath string) error {

	dll := syscall.NewLazyDLL("user32.dll")

	proc := dll.NewProc("SystemParametersInfoW")
	_t, _ := syscall.UTF16PtrFromString(imagePath)
	ret, _, _ := proc.Call(20, 1, uintptr(unsafe.Pointer(_t)), 0x1|0x2)
	if ret != 1 {
		return errors.New("系统调用失败")
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// DownloadImage 下载图片,保存并返回保存的文件名的绝对路径
func DownloadImage(imageURL string) (string, error) {
	fileName := EncodeMD5(imageURL)
	path := CurrentPathDir + fmt.Sprintf("%s", fileName) + ".jpg"
	fmt.Println("校验图片是否已存在", path)
	exist := Exists(path)
	if !exist {

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
		err = ioutil.WriteFile(path, body, 0755)
		if err != nil {
			return "", err
		}
	} else {
		fmt.Println("壁纸：", fileName, "已存在,不用下载")
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
