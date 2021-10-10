package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"view/src/com.jenkin.view/wallpaper"
)

var userNameLabelPoint *walk.Label
var passwordLabelPoint *walk.Label
var passwordEditPoint *walk.TextEdit
var userCodeEditPoint *walk.TextEdit
var loginButtonPoint *walk.PushButton
var nextButtonPoint *walk.PushButton
var windowMain *walk.MainWindow
var imgViewPoint *walk.ImageView
var imgPath string
var children []Widget

func main() {

	initUI()

}

func initUI() {

	userNameLabel := Label{Text: "用户名", AssignTo: &userNameLabelPoint}
	passwordLabel := Label{Text: "密码", AssignTo: &passwordLabelPoint}
	passwordEdit := TextEdit{Name: "密码", AssignTo: &passwordEditPoint, MaxSize: Size{100, 30}}
	userCodeEdit := TextEdit{Name: "用户名", AssignTo: &userCodeEditPoint, MaxSize: Size{100, 30}}
	loginButton := PushButton{
		ColumnSpan: 2,
		MaxSize:    Size{100, 30},
		Text:       "登录",
		OnClicked: func() {
			go execLogin(windowMain, userCodeEditPoint, passwordEditPoint)
		},
		AssignTo: &loginButtonPoint,
	}
	nextButton := PushButton{
		ColumnSpan: 2,
		MaxSize:    Size{100, 30},
		Text:       "下一张",
		OnClicked: func() {
			go wallpaper.Next()
		},
		Visible:  false,
		AssignTo: &nextButtonPoint,
	}
	imgPath = "app.ico"
	imgView := ImageView{
		Image:    imgPath,
		Margin:   0,
		Mode:     ImageViewModeShrink,
		Visible:  false,
		AssignTo: &imgViewPoint,
	}
	children = []Widget{
		userNameLabel,
		userCodeEdit,
		passwordLabel,
		passwordEdit,
		loginButton,
		nextButton,
		imgView,
	}

	_ = MainWindow{
		AssignTo: &windowMain,
		Children: []Widget{
			Composite{

				Layout:   Grid{Columns: 2},
				Children: children,
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
		go wallpaper.Start()
		showPreImage()
	} else {
		walk.MsgBox(
			mainWindow,
			"登录失败",
			"请检查用户名或密码是否错误",
			walk.MsgBoxOK)

	}
}

func showPreImage() {
	fmt.Println("预览图片")
	loginButtonPoint.SetVisible(false)
	userNameLabelPoint.SetVisible(false)
	passwordEditPoint.SetVisible(false)
	userCodeEditPoint.SetVisible(false)
	passwordLabelPoint.SetVisible(false)
	nextButtonPoint.SetVisible(true)
	imgViewPoint.SetVisible(true)
	for s := range wallpaper.PreImageCh {
		fmt.Println("收到预览图片：", s)
		image, _ := walk.NewImageFromFileForDPI(s, 96)
		setImage := imgViewPoint.SetImage(image)
		if setImage != nil {
			fmt.Println("设置预览失败：", setImage)
		}
	}
}
