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

// browser can be chrome or firefox
// chrome ==> chromedriver
// firefox ==> geckodriver
func getSelenium(url string, browser ...string) *seleniumStruct {

	var browserv string

	drivermap := map[string]string{
		"chrome":  "chromedriver",
		"firefox": "geckodriver",
	}

	if len(browser) != 0 {
		browserv = browser[0]
	}
	// firefoxDriverPath = "/usr/local/bin/geckodriver"
	chromeDriverPath, err := exec.LookPath(drivermap[browserv])
	Panicerr(err)
	servicePort := Int(randint(30000, 65535))

	opts := []selenium.ServiceOption{
		// selenium.Output(os.Stderr), // Output debug information to STDERR.
	}
	// selenium.SetDebug(true)
	var service *selenium.Service
	if browserv == "firefox" {
		service, err = selenium.NewGeckoDriverService(chromeDriverPath, servicePort, opts...)
		Panicerr(err)
	} else if browserv == "chrome" {
		service, err = selenium.NewChromeDriverService(chromeDriverPath, servicePort, opts...)
		Panicerr(err)
	}

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", servicePort))
	Panicerr(err)

	err = wd.Get(url)
	Panicerr(err)

	return &seleniumStruct{
		service: service,
		driver:  wd,
	}
}

func (c *seleniumStruct) Close() {
	c.service.Stop()
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

func (c *seleniumElementStruct) Text() *stringStruct {
	s, err := c.element.Text()
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
