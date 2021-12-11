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

func getSelenium(url string) *seleniumStruct {
	// firefoxDriverPath = "/usr/local/bin/geckodriver"
	chromeDriverPath, err := exec.LookPath("chromedriver")
	panicerr(err)
	servicePort := Int(randint(30000, 65535))

	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	// selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(chromeDriverPath, servicePort, opts...)
	panicerr(err)

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", servicePort))
	panicerr(err)

	err = wd.Get(url)
	panicerr(err)

	return &seleniumStruct{
		service: service,
		driver:  wd,
	}
}

func (c *seleniumStruct) close() {
	c.service.Stop()
	c.driver.Close()
}

func (c *seleniumStruct) cookie() (co string) {
	cookies, err := c.driver.GetCookies()
	panicerr(err)
	var coo []string
	for _, cookie := range cookies {
		coo = append(coo, cookie.Name+"="+cookie.Value)
	}
	return String(";").Join(coo).Get()
}

func (c *seleniumStruct) url() string {
	u, err := c.driver.CurrentURL()
	panicerr(err)
	return u
}

func (c *seleniumStruct) scrollRight(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy("+Str(pixel)+",0);", []interface{}{})
	panicerr(err)
}

func (c *seleniumStruct) scrollLeft(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy("+Str(pixel*-1)+",0);", []interface{}{})
	panicerr(err)
}

func (c *seleniumStruct) scrollUp(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy(0, "+Str(pixel*-1)+");", []interface{}{})
	panicerr(err)
}

func (c *seleniumStruct) scrollDown(pixel int) {
	_, err := c.driver.ExecuteScript("window.scrollBy(0, "+Str(pixel)+");", []interface{}{})
	panicerr(err)
}

func (c *seleniumStruct) resizeWindow(width int, height int) *seleniumStruct {
	cur, err := c.driver.CurrentWindowHandle()
	panicerr(err)

	err = c.driver.ResizeWindow(cur, width, height)
	panicerr(err)

	return c
}

type seleniumElementStruct struct {
	element selenium.WebElement
}

func (c *seleniumStruct) find(xpath string, nowait ...bool) *seleniumElementStruct {
	if len(nowait) != 0 && nowait[0] {
		we, err := c.driver.FindElement(selenium.ByXPATH, xpath)
		panicerr(err)
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

func (c *seleniumElementStruct) clear() *seleniumElementStruct {
	err := c.element.Clear()
	panicerr(err)
	return c
}

func (c *seleniumElementStruct) click() *seleniumElementStruct {
	err := c.element.Click()
	panicerr(err)
	return c
}

func (c *seleniumElementStruct) text() string {
	s, err := c.element.Text()
	panicerr(err)
	return s
}

func (c *seleniumElementStruct) input(s string) *seleniumElementStruct {
	err := c.element.SendKeys(s)
	panicerr(err)
	return c
}

func (c *seleniumElementStruct) submit() *seleniumElementStruct {
	err := c.element.Submit()
	panicerr(err)
	return c
}

func (c *seleniumElementStruct) pressEnter() *seleniumElementStruct {
	c.input(selenium.EnterKey)
	return c
}
