package wallpaper

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"view/src/com.jenkin.view/wallpaperstruct"
)

//登录状态
var LoginCh chan bool

//预览图
var PreImageCh chan string

var Conn net.Conn
var Opt wallpaperstruct.Option
var Strategy wallpaperstruct.WallStrategy
var start = true

// 从服务端接受的指令
const (
	//更换壁纸
	CmdChangeWallpaper = "changeWallpaper"
	//更换策略
	CmdChangeStrategy = "changeStrategy"
	//登录成功
	CmdLoginSuccess = "loginSuccess"
	//登录失败
	CmdLoginFailed = "loginFailed"
)

// 发送给服务端的指令
const (
	//登录
	CmdLogin = "login"
	//下一张壁纸
	CmdNext = "next"
)

//func main() {
//	Start()
//}
//建立连接之前的初始化工作
func PreConn() bool {
	//无缓冲channel，阻塞
	LoginCh = make(chan bool)
	PreImageCh = make(chan string, 10)
	Conn = getConnection()
	if Conn == nil {
		return false
	}
	return true
}

//处理服务端响应数据
func handleEvent() {
	go WallpaperHandler()
}

//启动循环获取壁纸
func Start() {
	start = true
	if Conn != nil {
		// 根据策略循环
		loopNext()
	}
	start = false
}

//获取连接
func getConnection() net.Conn {
	//conn, err := net.Dial("tcp", "127.0.0.1:9010")
	conn, err := net.Dial("tcp", "tencent.jenkin.tech:9010")
	if err != nil {
		log.Println("客户端建立连接失败")
		return nil
	}
	log.Println("连接获取成功")
	return conn
}

//登录
func Login(code string, password string) {
	option := wallpaperstruct.Option{
		OperateType: CmdLogin,
		UserCode:    code,
		OperateData: password,
	}
	writeServer(option, Conn)
	handleEvent()
}

// 循环获取壁纸
func loopNext() {

	for Strategy.OnFlag {
		second := getSleepSecondTime()
		log.Println("间隔时间为：", second)
		if &Opt != nil && second > 0 {
			getNextFromServer(Opt, Conn)
			time.Sleep(time.Duration(second) * time.Second)
		} else {
			log.Println(second, Strategy)
			log.Println("循环检测未登录")
			time.Sleep(2 * time.Second)
		}
	}

}

/**
获取循环时间
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

//处理服务端返回的指令
func WallpaperHandler() {
	log.Println("监听壁纸返回")
	//缓存 conn 中的数据
	buf := make([]byte, 1024*10)
	for {

		if &Opt == nil {
			log.Println("监听壁纸未登录")
			time.Sleep(time.Second * 2)
			continue
		}

		log.Println("等待数据")
		info, err := readInfo(Conn, buf)
		if err != nil {
			tryReconnect()
		}
		if info != nil {
			data := *info
			opType := data.OperateType
			if opType != "" {
				log.Println("操作类型：", opType, " 操作人：", data.UserCode, " 数据：", data.OperateData)
				switch opType {
				case CmdChangeWallpaper:
					changeWallpaper(data, Conn)
				case CmdChangeStrategy:
					changeStrategy(data, Conn)
				case CmdLoginSuccess:
					loginSuccess(data, Conn)
				case CmdLoginFailed:
					loginFailed(data, Conn)
				default:
					log.Print("非法请求")
				}
			} else {
				log.Println("操作类型为空")
			}
		}
	}
}

// 断线重连
func tryReconnect() {
	for {
		connection := getConnection()
		if connection == nil {
			log.Println("5秒后重试获取连接")
			time.Sleep(time.Second * 5)
		} else {
			Conn = connection
			break
		}
	}
}

// 登录失败
func loginFailed(option wallpaperstruct.Option, conn net.Conn) {
	LoginCh <- false
	//log.Println("登录失败，用户未注册，或未配置规则,5秒后退出")
	//time.Sleep(time.Second * 5)
	//os.Exit(0)
}

// 登陆成功操作
func loginSuccess(option wallpaperstruct.Option, conn net.Conn) {
	opdata := option.OperateData

	strategy := &wallpaperstruct.WallStrategy{}
	JsonToStruct(opdata, strategy)
	Strategy = *strategy
	Opt = option
	log.Println("登录获取到的配置：", Strategy)
	log.Println("登录成功：", option)
	// 登录状态管道写入成功
	LoginCh <- true

}

// 更换策略
func changeStrategy(option wallpaperstruct.Option, conn net.Conn) {

	log.Println("策略变更，变更前：", Strategy, "变更后：", option.OperateData)
	strategy := &wallpaperstruct.WallStrategy{}
	JsonToStruct(option.OperateData, strategy)
	Strategy = *strategy
	Opt.OperateData = option.OperateData
	getNextFromServer(option, conn)
	// 禁用后重启
	if !start {
		go Start()
	}

}

// 下一张壁纸
func Next() {
	getNextFromServer(Opt, Conn)
}

// 从服务器获取下一张壁纸
func getNextFromServer(option wallpaperstruct.Option, conn net.Conn) {
	option.OperateType = CmdNext
	writeServer(option, conn)
}

//更换壁纸
func changeWallpaper(option wallpaperstruct.Option, conn net.Conn) {
	wall := &wallpaperstruct.Wallpaper{}
	JsonToStruct(option.OperateData, wall)
	log.Println("准备设置桌面")
	go SetWallpaper(wall.Img)
}

//读取服务端数据
func readInfo(conn net.Conn, buf []byte) (*wallpaperstruct.Option, error) {
	cnt, err := conn.Read(buf)
	if err != nil {
		log.Println("客户端读取数据失败 %s\n", err)
		return nil, err
	}
	data := buf[0:cnt]

	dataStr := string(data)
	res := &wallpaperstruct.Option{}
	JsonToStruct(dataStr, res)
	//回显服务器端回传的信息
	log.Println("服务器端回复: " + dataStr)
	return res, nil
}

//结构体转json
func StructToJson(data interface{}) string {
	res, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	return string(res)
}

//json转结构体
func JsonToStruct(text string, data interface{}) {
	err := json.Unmarshal([]byte(text), data)
	if err != nil {
		log.Println("json转结构体异常：", text)
	}
}

//心跳检测,暂时不用
func HeartBeatHandler(c net.Conn) {
	log.Println("监听心跳")
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
			log.Println("客户端读取数据失败 %s\n", err)
			time.Sleep(time.Second * 5)
			continue
		}
		//回显服务器端回传的信息
		log.Println("服务器端回复" + string(buf[0:cnt]))
		time.Sleep(time.Second * 5)
	}
}

//向服务端写数据
func writeServer(option wallpaperstruct.Option, conn net.Conn) {
	toJson := StructToJson(option)
	log.Println("发送 数据：", toJson)
	_, _ = conn.Write([]byte(toJson + "````"))
}
