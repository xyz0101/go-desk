package wallpaper

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
	"view/src/com.jenkin.view/wallpaperstruct"
)

var rLock sync.RWMutex
var Conn net.Conn
var Opt wallpaperstruct.Option
var Strategy wallpaperstruct.WallStrategy

//func main() {
//	Start()
//}

func Start(option wallpaperstruct.Option) {

	Conn = getConnection()
	if Conn == nil {
		fmt.Println("客户端建立连接失败,5秒后退出程序")
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}
	login()
	//go HeartBeatHandler(conn)
	go WallpaperHandler()
	// 根据策略循环
	loopNext()
}

func getConnection() net.Conn {
	//conn, err := net.Dial("tcp", "127.0.0.1:9010")
	conn, err := net.Dial("tcp", "tencent.jenkin.tech:9010")
	if err != nil {
		fmt.Println("客户端建立连接失败")
		return nil
	}
	fmt.Println("连接获取成功")
	return conn
}

func login() {
	option := wallpaperstruct.Option{
		OperateType: "login",
		//UserCode:"jenkin",
		UserCode:    os.Args[1],
		OperateData: "password",
	}
	writeServer(option, Conn)
}

func loopNext() {

	for {

		rLock.RLock()
		second := getSleepSecondTime()
		fmt.Println("间隔时间为：", second)
		if &Opt != nil && second > 0 {
			getNextFromServer(Opt, Conn)
			time.Sleep(time.Duration(second) * time.Second)
		} else {
			fmt.Println(second, Strategy)
			fmt.Println("循环检测未登录")
			time.Sleep(2 * time.Second)
		}
		rLock.RUnlock()
	}
}

/**
 MINUTE("second",0,"秒"),
 MINUTE("minute",1,"分钟"),
UN_START("hour",2,"小时"),
 WAITING("day",3,"天")
*/
func getSleepSecondTime() int {

	switch Strategy.TimeUnit {
	case 0:
		return Strategy.TimeGap
	case 1:
		return Strategy.TimeGap * 60
	case 2:
		return Strategy.TimeGap * 60 * 60
	case 3:
		return Strategy.TimeGap * 60 * 60 * 24
	default:
		return -1
	}

}

func WallpaperHandler() {
	fmt.Println("监听壁纸返回")
	//缓存 conn 中的数据
	buf := make([]byte, 1024*10)
	for {
		rLock.RLock()
		if &Opt == nil {
			fmt.Println("监听壁纸未登录")
			time.Sleep(time.Second * 2)
			continue
		}
		rLock.RUnlock()
		fmt.Println("等待数据")
		info, err := readInfo(Conn, buf)
		if err != nil {
			tryReconnect()
		}
		if info != nil {
			data := *info
			opType := data.OperateType
			if opType != "" {
				fmt.Println("操作类型：", opType, " 操作人：", data.UserCode, " 数据：", data.OperateData)
				switch opType {
				case "changeWallpaper":
					changeWallpaper(data, Conn)
				case "changeStrategy":
					changeStrategy(data, Conn)
				case "loginSuccess":
					loginSuccess(data, Conn)
				case "loginFailed":
					loginFailed(data, Conn)

				}
			} else {
				fmt.Println("操作类型为空")
			}
		}
	}
}

func tryReconnect() {
	for {
		connection := getConnection()
		if connection == nil {
			fmt.Println("5秒后重试获取连接")
			time.Sleep(time.Second * 5)
		} else {
			Conn = connection
			break
		}
	}
}

func loginFailed(option wallpaperstruct.Option, conn net.Conn) {
	fmt.Println("登录失败，用户未注册，或未配置规则,5秒后退出")
	time.Sleep(time.Second * 5)
	os.Exit(0)
}

func loginSuccess(option wallpaperstruct.Option, conn net.Conn) {
	rLock.Lock()
	opdata := option.OperateData

	strategy := &wallpaperstruct.WallStrategy{}
	JsonToStruct(opdata, strategy)
	Strategy = *strategy
	Opt = option
	fmt.Println("登录获取到的配置：", Strategy)
	fmt.Println("登录成功：", option)
	defer rLock.Unlock()

}

func changeStrategy(option wallpaperstruct.Option, conn net.Conn) {
	//for {
	fmt.Println("策略变更，变更前：", Strategy, "变更后：", option.OperateData)
	strategy := &wallpaperstruct.WallStrategy{}
	JsonToStruct(option.OperateData, strategy)
	Strategy = *strategy
	Opt.OperateData = option.OperateData
	getNextFromServer(option, conn)
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

func readInfo(conn net.Conn, buf []byte) (*wallpaperstruct.Option, error) {
	cnt, err := conn.Read(buf)
	if err != nil {
		fmt.Println("客户端读取数据失败 %s\n", err)
		return nil, err
	}
	data := buf[0:cnt]

	dataStr := string(data)
	res := &wallpaperstruct.Option{}
	JsonToStruct(dataStr, res)
	//回显服务器端回传的信息
	fmt.Println("服务器端回复: " + dataStr)
	return res, nil
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
		fmt.Println("json转结构体异常：", text)
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
