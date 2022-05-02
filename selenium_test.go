package golanglibs

import (
	"testing"
)

func TestSelenium(t *testing.T) {
	// sn := Tools.SeleniumLocal().Get("https://google.com").ResizeWindow(1300, 1000)
	// defer sn.Close()

	// // 登录
	// //Lg.Trace("选语言")
	// sn.Find(`/html/body/div/div[1]/div[1]/div[2]/form/div[3]/div/select/option[2]`).Click()
	// //Lg.Trace("输入用户名")
	// sn.Find(`//*[@id="login"]`).Clear().Input("Starview")
	// //Lg.Trace("输入密码")
	// sn.Find(`//*[@id="password"]`).Clear().Input("SqKhvbRbA4jYBPCN")
	// Lg.Trace("登录")
	// sn.Find(`/html/body/div[1]/div[1]/div[1]/div[2]/form/center/div/input`).Click()

	sn := Tools.SeleniumRemote("http://192.168.168.22:4444").Get("https://projectkorra.com/forum/forums/general-discussion.15").ResizeWindow(1300, 1000)
	defer sn.Close()

	el := sn.Find("/html/body/div[1]/div[4]/div/div[2]/div/div[2]/div/div[2]/div[2]/div/div[1]/div[1]/div[2]/div[1]/a")

	Lg.Trace(el.Attribute("href"))
}
