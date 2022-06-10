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
	// sn := Tools.SeleniumRemote("http://192.168.168.22:4444").Get("https://projectkorra.com/forum/forums/general-discussion.15").ResizeWindow(1300, 1000)
	// defer sn.Close()
	// el := sn.Find("/html/body/div[1]/div[4]/div/div[2]/div/div[2]/div/div[2]/div[2]/div/div[1]/div[1]/div[2]/div[1]/a")
	// Lg.Trace(el.Attribute("href"))

	// sn := Tools.SeleniumLocal()
	// defer sn.Close()
	// go sn.Get("http://www.gshadow.com/Aban/index.php?name=Forums&file=viewtopic&t=5110").ResizeWindow(1300, 1000)
	// Time.Sleep(15)
	// Print(1)
	// sn.Get("http://google.com")
	// Print(sn.Title())

	sn := Tools.SeleniumRemote("http://localhost:22222/wd/hub")
	// Print(sn.GetSession())
	sn.SetSession("fa8ad2c197992ea57fa0eea2d0185405")
	// sn.Get("https://google.com")
	sn.Find("/html/body/div[1]/div[3]/form/div[1]/div[1]/div[1]/div/div[2]/input").Input("abc")
}
