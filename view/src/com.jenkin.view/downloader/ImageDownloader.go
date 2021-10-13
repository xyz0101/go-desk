package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Downloader struct {
	io.Reader
	Total   int64
	Current int64
	Name    string
}

//var percentMap *lru.LRUCache
var PercentCh = make(chan float64, 10)

func (d *Downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.Current += int64(n)
	percent := float64(d.Current*10000/d.Total) / 100
	fmt.Printf("\r正在下载，下载进度：%.2f%%", percent)
	if d.Current == d.Total {
		fmt.Printf("\r下载完成，下载进度：%.2f%%", 100.00)
	}
	//if percentMap == nil {
	//	percentMap = lru.GetLRUCache(10)
	//}
	//percentMap.Put(d.Name,percent)
	PercentCh <- percent
	return
}

func downloadFile(url, filePath string) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	file, err := os.Create(filePath)
	defer func() {
		_ = file.Close()
	}()
	downloader := &Downloader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
		Name:   filePath,
	}
	if _, err := io.Copy(file, downloader); err != nil {
		log.Fatalln(err)
	}
}

var wg sync.WaitGroup

func Download(url string, path string) {
	wg.Add(1)
	downloadFile(url, path)

	wg.Wait()
}

//func main() {
//	task := make(map[string]string)
//	task["http://img.aibizhi.adesk.com/616171aee7bce73098e3ac62?sign=30746747ac29d123bac18e19928dc845&t=6166118d"] = "D:\\新建文件夹\\test.jpg"
//	for k, v := range task {
//		wg.Add(1)
//		downloadFile(k, v)
//	}
//	wg.Wait()
//}
