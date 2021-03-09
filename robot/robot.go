package robot

import (
	"ama/middlevars"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"os"
)

func init(){
	initPre()

}
func Init(){
	initPre()
}

func GetWebdriver()selenium.WebDriver{
	const (
		seleniumPath = `./chromedriver`
		port            = 9515
	)

	opts := []selenium.ServiceOption{

	}


	//selenium.SetDebug(true)

	_, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		os.Exit(3)
	}

	////server关闭之后，chrome窗口也会关闭
	//defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		//"browserName": "/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome Dev",
		"browserName": "Google Chrome Dev",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome Dev",
		Args: []string{
			//静默执行请求
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	wbd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	return wbd


}


func initPre(){
	const (
		seleniumPath = `./chromedriver`
		port            = 9515
	)

	opts := []selenium.ServiceOption{

	}


	//selenium.SetDebug(true)

	_, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		os.Exit(3)
	}

	////server关闭之后，chrome窗口也会关闭
	//defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		//"browserName": "/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome Dev",
		"browserName": "Google Chrome Dev",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "/Applications/Google Chrome Dev.app/Contents/MacOS/Google Chrome Dev",
		Args: []string{
			//静默执行请求
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	wbd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	middlevars.Chrome = wbd


}


//开启新的tab
func OpenNewTab(wbd selenium.WebDriver , url string){
	wbd.ExecuteScript("window.open('"+url+"');",nil)
	//获取当前浏览器窗口的所有tab列表
	//windowhandler , _ := wbd.WindowHandles()
	//wbd.SwitchWindow(windowhandler[0])
	//return getSource(wbd , url)

}

func GetOnePageSource(wbd selenium.WebDriver , url string)string{
	wbd.Get(url)
	return getSource(wbd , url)

}


func getSource(wbd selenium.WebDriver,url string)string{
	pagesource , err := wbd.PageSource()
	if err != nil{
		//这里添加错误后的重试操作
		//wbd.Refresh()
		//time.Sleep(time.Second * 2)
		log.Printf("get pagesource %s error:%v",url , err)

	}
	//defer wbd.Quit()
	return pagesource
}

func SwitchTab(wbd selenium.WebDriver , tabindex int,url string)string{
	fmt.Println("切换浏览器标签到,",tabindex)
	windowshandler , err := wbd.WindowHandles()
	if err != nil{
		log.Println("获取浏览器tab对象失败🙅‍♂️,err ",err)
		return ""
	}
	wbd.SwitchWindow(windowshandler[tabindex])
	wbd.ExecuteScript("window.open('"+url+"');",nil)
	return getSource(wbd , url)
}

func CloseTab(wbd selenium.WebDriver){
	//获去窗口句柄
	handler , err := wbd.WindowHandles()
	if err != nil{
		log.Println("closeTab func get windowhandler err: ",err)
		return
	}
	if len(handler) > 30{
		for i:=0;i<25;i++{
			wbd.SwitchWindow(handler[i])
			wbd.Close()
		}
	}


}