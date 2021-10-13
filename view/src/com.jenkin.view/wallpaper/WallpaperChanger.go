package wallpaper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"unsafe"
	"view/src/com.jenkin.view/downloader"
)

const (
	CurrentPathDir = "cache/"
	//最大允许缓存的图片数目
	MaxSize = 100
)

var dll *syscall.LazyDLL
var proc *syscall.LazyProc

func Init() {
	_ = os.Mkdir(CurrentPathDir, 0755)
	dll = syscall.NewLazyDLL("user32.dll")
	proc = dll.NewProc("SystemParametersInfoW")
}

// EncodeMD5 MD5编码
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

// SetWindowsWallpaper 设置windows壁纸
func setWindowsWallpaper(imagePath string) error {

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
	deleteLastWhenOverMaxSize()
	if !exist {
		downloader.Download(imageURL, path)
	} else {

		fmt.Println("壁纸：", fileName, "已存在,不用下载")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

//删除超过限制文件
func deleteLastWhenOverMaxSize() {
	//file, _ := os.OpenFile(CurrentPathDir,os.O_RDONLY,os.ModeDir)
	infos, _ := ioutil.ReadDir(CurrentPathDir)
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().Unix() > infos[j].ModTime().Unix()
	})
	if len(infos) > MaxSize {
		info := infos[len(infos)-1]
		name := CurrentPathDir + info.Name()
		e := os.Remove(name)
		fmt.Println("滚动删除文件：", name)
		if e != nil {
			fmt.Println("文件滚动删除失败", e)
		} else {
			fmt.Println("文件数量超过 ", MaxSize, " 文件滚动删除成功")
		}
	}

}

func SetWallpaper(imageURL string) {

	fmt.Println("下载图片...", imageURL)
	imagePath, err := DownloadImage(imageURL)
	if err != nil {
		fmt.Println("下载图片失败: " + err.Error())
		return
	}
	fmt.Println("设置桌面...")
	PreImageCh <- imagePath
	err = setWindowsWallpaper(imagePath)
	if err != nil {
		fmt.Println("设置桌面背景失败: " + err.Error())
		return
	}
	fmt.Println("桌面设置成功")
}
