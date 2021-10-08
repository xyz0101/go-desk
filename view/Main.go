package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"view/src/com.jenkin.view/wallpaper"
)

func main() {
	//option := initOption()
	//wallpaper.Start(option)
	initUI()

}

func initUI() {
	var userTE, pwTE *walk.TextEdit

	userNameLabel := Label{Text: "用户名"}
	passwordLabel := Label{Text: "密码"}
	passwordEdit := TextEdit{Name: "密码", AssignTo: &pwTE, MaxSize: Size{100, 30}}
	userCodeEdit := TextEdit{Name: "用户名", AssignTo: &userTE, MaxSize: Size{100, 30}}
	//var inTE, outTE *walk.TextEdit
	var windowMain *walk.MainWindow
	_ = MainWindow{
		AssignTo: &windowMain,
		Children: []Widget{
			Composite{

				Layout: Grid{Columns: 2},
				Children: []Widget{
					userNameLabel,
					userCodeEdit,
					passwordLabel,
					passwordEdit,

					PushButton{
						ColumnSpan: 2,
						MaxSize:    Size{100, 30},
						Text:       "登录",
						OnClicked: func() {
							go execLogin(windowMain, userTE, pwTE)
						},
					},
					PushButton{
						ColumnSpan: 2,
						MaxSize:    Size{100, 30},
						Text:       "下一张",
						OnClicked: func() {
							go wallpaper.Next()
						},
					},
				},
			},
		},
		Size:   Size{500, 400},
		Title:  "GoDesk!",
		Layout: HBox{},
	}.Create()

	windowMain.Run()
}

func execLogin(mainWindow *walk.MainWindow, userTE *walk.TextEdit, pwTE *walk.TextEdit) {
	fmt.Println("登录")
	wallpaper.PreConn()
	wallpaper.Login(userTE.Text(), pwTE.Text())
	res := <-wallpaper.LoginCh
	if res {
		wallpaper.Start()
	} else {
		walk.MsgBox(
			mainWindow,
			"Title",
			"Message",
			walk.MsgBoxYesNoCancel)

	}
}
