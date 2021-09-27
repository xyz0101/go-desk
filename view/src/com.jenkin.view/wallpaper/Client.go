package wallpaper

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
	"view/src/com.jenkin.view/wallpaperstruct"
)

var Conn net.Conn
var Opt wallpaperstruct.Option

//func main() {
//	Start()
//}

func Start(option wallpaperstruct.Option) {

	conn, err := net.Dial("tcp", "127.0.0.1:23456")
	if err != nil {
		fmt.Println("客户端建立连接失败")
		return
	}
	Conn = conn
	Opt = option
	//go HeartBeatHandler(conn)
	go WallpaperHandler(conn)
	// 根据策略循环
	loopNext(option, conn)
}

func loopNext(option wallpaperstruct.Option, conn net.Conn) {
	for {
		getNextFromServer(option, conn)
		time.Sleep(time.Second * 600)
	}
}

func WallpaperHandler(conn net.Conn) {
	fmt.Println("监听壁纸返回")
	//缓存 conn 中的数据
	buf := make([]byte, 1024*10)
	for {
		fmt.Println("等待数据")
		info := readInfo(conn, buf)
		if info != nil {
			data := *info
			opType := data.OperateType
			if opType != "" {
				fmt.Println("操作类型：", opType, " 操作人：", data.UserCode, " 数据：", data.OperateData)
				switch opType {
				case "changeWallpaper":
					changeWallpaper(data, conn)
				case "changeStrategy":
					changeStrategy(data, conn)
				}
			} else {
				fmt.Println("操作类型为空")
			}
		}
	}
}

func changeStrategy(option wallpaperstruct.Option, conn net.Conn) {
	//for {
	getNextFromServer(option, conn)
	time.Sleep(time.Second * 10)
	//}
}

func getNextFromServer(option wallpaperstruct.Option, conn net.Conn) {
	option.OperateType = "next"
	writeServer(option, conn)
}

func changeWallpaper(option wallpaperstruct.Option, conn net.Conn) {
	wall := &wallpaperstruct.Wallpaper{}
	JsonToStruct(option.OperateData, wall)
	fmt.Println("准备设置桌面")
	SetWallpaper(wall.Img)
}

func readInfo(conn net.Conn, buf []byte) *wallpaperstruct.Option {
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println("客户端读取数据失败 %s\n", err)
		return nil
	}
	data := buf[0:cnt]

	dataStr := string(data)
	res := &wallpaperstruct.Option{}
	JsonToStruct(dataStr, res)
	//回显服务器端回传的信息
	fmt.Println("服务器端回复: " + dataStr)
	return res
}

func StructToJson(data interface{}) string {
	res, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	return string(res)
}

func JsonToStruct(text string, data interface{}) {
	err := json.Unmarshal([]byte(text), data)
	if err != nil {
		fmt.Println("json转结构体异常：", err)
	}
}

func HeartBeatHandler(c net.Conn) {
	fmt.Println("监听心跳")
	//缓存 conn 中的数据
	buf := make([]byte, 1024*10)
	for {
		req := wallpaperstruct.Option{
			UserCode:    "jenkin",
			OperateType: "heart",
		}
		//客户端请求数据写入 conn，并传输
		writeServer(req, c)

		cnt, err := c.Read(buf)
		if err != nil {
			fmt.Println("客户端读取数据失败 %s\n", err)
			time.Sleep(time.Second * 5)
			continue
		}
		//回显服务器端回传的信息
		fmt.Println("服务器端回复" + string(buf[0:cnt]))
		time.Sleep(time.Second * 5)
	}
}

func writeServer(option wallpaperstruct.Option, conn net.Conn) {
	toJson := StructToJson(option)
	fmt.Println("发送 数据：", toJson)
	_, _ = conn.Write([]byte(toJson + "````"))
}
