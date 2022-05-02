package golanglibs

import (
	"fmt"
	"os/exec"

	"github.com/tebeka/selenium"
)

type seleniumStruct struct {
	service *selenium.Service
	driver  selenium.WebDriver
}

func getSeleniumLocal() *seleniumStruct {
	driverPath, err := exec.LookPath("chromedriver")
	Panicerr(err)
	servicePort := Int(randint(30000, 65535))

	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	// selenium.SetDebug(true)
	var service *selenium.Service
	service, err = selenium.NewChromeDriverService(driverPath, servicePort, opts...)
	Panicerr(err)

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", servicePort))
	Panicerr(err)

	return &seleniumStruct{
		service: service,
		driver:  wd,
	}
}

func getSeleniumRemote(serverURL string) *seleniumStruct {
	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}

	if String(serverURL).EndsWith("/") {
		serverURL = String(serverURL).Strip("/").S
	}

	if !String(serverURL).EndsWith("/wd/hub") {
		serverURL = serverURL + "/wd/hub"
	}

	wd, err := selenium.NewRemote(caps, serverURL)
	Panicerr(err)

	return &seleniumStruct{
		service: nil,
		driver:  wd,
	}
}

func (c *seleniumStruct) Get(url string) *seleniumStruct {
	err := c.driver.Get(url)
	Panicerr(err)
	return c
}

func (c *seleniumStruct) Close() {
	if c.service != nil {
		c.service.Stop()
	}
	c.driver.Close()
}

func (c *seleniumStruct) Cookie() (co string) {
	cookies, err := c.driver.GetCookies()
	Panicerr(err)
	var coo []string
	for _, cookie := range cookies {
		coo = append(coo, cookie.Name+"="+cookie.Value)
	}
	return String(";").Join(coo).Get()
}

func (c *seleniumStruct) Url() string {
	u, err := c.driver.CurrentURL()
	Panicerr(err)
	return u
}

func (c *seleniumStruct) Title() string {
	u, err := c.driver.Title()
	Panicerr(err)
	return u
}

func (c *seleniumStruct) ScrollRight(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy("+Str(pixel)+",0);", []interface{}{})
	Panicerr(err)
}

func (c *seleniumStruct) ScrollLeft(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy("+Str(pixel*-1)+",0);", []interface{}{})
	Panicerr(err)
}

func (c *seleniumStruct) ScrollUp(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy(0, "+Str(pixel*-1)+");", []interface{}{})
	Panicerr(err)
}

func (c *seleniumStruct) ScrollDown(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy(0, "+Str(pixel)+");", []interface{}{})
	Panicerr(err)
}

func (c *seleniumStruct) ResizeWindow(width int, height int) *seleniumStruct {
	cur, err := c.driver.CurrentWindowHandle()
	Panicerr(err)

	err = c.driver.ResizeWindow(cur, width, height)
	Panicerr(err)

	return c
}

type seleniumElementStruct struct {
	element selenium.WebElement
}

func (c *seleniumStruct) Find(xpath string, nowait ...bool) *seleniumElementStruct {
	if len(nowait) != 0 && nowait[0] {
		we, err := c.driver.FindElement(selenium.ByXPATH, xpath)
		if String("no such element").In(Str(err)) {
			return nil
		}
		Panicerr(err)
		return &seleniumElementStruct{element: we}
	} else {
		for {
			we, err := c.driver.FindElement(selenium.ByXPATH, xpath)
			if String("no such element").In(Str(err)) {
				sleep(1)
				continue
			}
			return &seleniumElementStruct{element: we}
		}
	}
}

func (c *seleniumStruct) PageSource() string {
	source, err := c.driver.PageSource()
	Panicerr(err)
	return source
}

func (c *seleniumElementStruct) Clear() *seleniumElementStruct {
	err := c.element.Clear()
	Panicerr(err)
	return c
}

func (c *seleniumElementStruct) Click() *seleniumElementStruct {
	err := c.element.Click()
	Panicerr(err)
	return c
}

func (c *seleniumElementStruct) Text() *StringStruct {
	s, err := c.element.Text()
	Panicerr(err)
	return String(s)
}

func (c *seleniumElementStruct) Attribute(name string) *StringStruct {
	s, err := c.element.GetAttribute(name)
	Panicerr(err)
	return String(s)
}

func (c *seleniumElementStruct) Input(s string) *seleniumElementStruct {
	err := c.element.SendKeys(s)
	Panicerr(err)
	return c
}

func (c *seleniumElementStruct) Submit() *seleniumElementStruct {
	err := c.element.Submit()
	Panicerr(err)
	return c
}

func (c *seleniumElementStruct) PressEnter() *seleniumElementStruct {
	c.Input(selenium.EnterKey)
	return c
}
