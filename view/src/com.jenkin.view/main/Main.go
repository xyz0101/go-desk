package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"view/src/com.jenkin.view/wallpaper"
	"view/src/com.jenkin.view/wallpaperstruct"
)

func main() {
	option := initOption()
	wallpaper.Start(option)
	//err := ui.Main(setupUI)
	//if err != nil {
	//	panic(err)
	//}
}

func setupUI() {

	// 生成：水平容器
	box := ui.NewHorizontalBox()

	// 往 垂直容器 中添加 控件
	box.Append(ui.NewLabel("壁纸分类"), false)

	// 生成：窗口（标题，宽度，高度，是否有 菜单 控件）
	window := ui.NewWindow(`GoDesk`, 1366, 720, true)

	// 窗口容器绑定
	window.SetChild(box)

	// 设置：窗口关闭时
	window.OnClosing(func(*ui.Window) bool {
		// 窗体关闭
		ui.Quit()
		return true
	})

	// 窗体显示
	window.Show()
	option := initOption()
	wallpaper.Start(option)
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
