package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/win"
	"view/src/com.jenkin.view/wallpaperstruct"
)

func main() {
	//option := initOption()
	//wallpaper.Start(option)
	initUI()

}

func initUI() {
	window, _ := walk.NewMainWindow()

	// 设置窗体标题
	window.SetTitle(`你好世界！`)

	// 设置窗体的宽高
	window.SetWidth(400)
	window.SetHeight(400)

	// 设置窗体生成在屏幕的正中间
	// 窗体横坐标 = ( 屏幕宽度 - 窗体宽度 ) / 2
	// 窗体纵坐标 = ( 屏幕高度 - 窗体高度 ) / 2
	window.SetX((int(win.GetSystemMetrics(0)) - window.Width()) / 2)
	window.SetY((int(win.GetSystemMetrics(1)) - window.Height()) / 2)

	// 设置窗体为显示状态（默认：隐藏状态）
	window.Show()

	// 运行窗体
	window.Run()

}

func initOption() wallpaperstruct.Option {
	str := "{ " +
		"      \"timeGap\": 1, " +
		"      \"strategyCode\": \"RandomStrategy\"," +
		"      \"categories\": [" +
		"        \"4e4d610cdf714d2966000000\"," +
		"        \"4e4d610cdf714d2966000002\"," +
		"        \"4e4d610cdf714d2966000001\"" +
		"      ]," +
		"      \"userCode\": \"jenkin\"," +
		"      \"timeUnit\": 1" +
		"    }"
	op := wallpaperstruct.Option{
		"next",
		"jenkin",
		str,
	}
	return op
}

//func buildLabelText(font string, size int, blod bool, text string) {
//	label := ui.NewLabel(text)
//
//}
