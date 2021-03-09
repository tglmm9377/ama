package httpserver

import (
	"ama/db"
	"ama/jobstate"
	"ama/middlevars"
	"ama/robot"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//通过asin 获取所需信息 ，
//1. 获取对应listing review 地址
func GetUrlWithAsin(asin string)(bool , string){
	//defer Lock.Done()
	wbd := robot.GetWebdriver()
	defer wbd.Quit()
	Getlisturl  := "https://"+middlevars.Site+"/s/ref=nb_sb_noss_2?url=search-alias%3Daps&field-keywords="+asin
	fmt.Println(Getlisturl)
	html := robot.GetOnePageSource(wbd , Getlisturl)
	dom , _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	selection := dom.Find("span[data-component-type=s-product-image]").First().Find("a")
	val , ok := selection.Attr("href")
	if !ok{
		log.Println("url='' ，获取review地址失败，检查浏览器返回的网页是否正常,出错的asin: ",asin,"开始重试")
		return false , ""
	}
	fmt.Println("val:",val)
	urlSplit := strings.Split(val ,"/")[1]
	url := "https://"+middlevars.Site+"/"+urlSplit+"/"+"product-reviews/"+asin+"/ref=cm_cr_arp_d_viewopt_rvwer?ie=UTF8&reviewerType=avp_only_reviews&pageNumber="
	return true ,url

}


//2 . 获取总review数量
func GetReviewNum(asin string , url string)int{
	//defer Lock.Done()
	wdb := robot.GetWebdriver()
	defer wdb.Quit()
	pagesource := robot.GetOnePageSource(wdb , url)

	pattern := ".*?[|](.*?)global reviews.*?"
	re := regexp.MustCompile(pattern)
	numstr := strings.TrimSpace(re.FindStringSubmatch(pagesource)[1])
	//resolve 2,345 format
	num := strings.Join(strings.Split(numstr , ","),"")
	total , err := strconv.Atoi(num)
	if err != nil{
		//log.Println("")
		return 0
	}
	return total
}



func job2(reviewnum int, url , asin string){
	reviewNum := reviewnum
	baseurl := url
	var pageNum int
	var groupNum int
	pageNum = RightNum(reviewNum , 10)
	groupNum = RightNum(pageNum , 40)
	//记录每组 多少页
	groupPage := make(map[int]int)
	allPage := pageNum
	for i:=0;i<groupNum;i++{
		allPage = allPage - 40
		if allPage > 0{
			groupPage[i] = 40
		}else{
			groupPage[i] = allPage+40
		}

	}
	fmt.Println("分组后每组的页数",groupPage)
	for i, v := range groupPage{
		webdriver := robot.GetWebdriver()
		Lock.Add(1)
		go func(group int , page int,webdriver selenium.WebDriver) {
			for j:=1;j<=page;j++{
				fmt.Println( group , page , j)
				url := 	baseurl + strconv.Itoa(group*40+j)
				robot.OpenNewTab(webdriver , url)
				//robot.ParseData(robot.ParseHtml(robot.SwitchTab(webdriver , j,url)))
			}
			handler ,_ := webdriver.WindowHandles()
			for _ , hand := range handler{
				webdriver.SwitchWindow(hand)
				pageSource , _ := webdriver.PageSource()
				robot.ParseData(robot.ParseHtml(pageSource),asin)
			}
			defer webdriver.Quit()
			Lock.Done()
		}(i,v,webdriver)
	}

	func() {
		var j jobstate.Job
		j.Name = asin
		j.Status = "stoped"
		db.DB.Where("name=?",asin).Update(&j)
	}()

}

func RightNum(num int , remainder int)(pageNum int){
	if num % remainder > 0 {
		pageNum = num / remainder + 1
	}else{
		pageNum = num / remainder
	}
	return pageNum
}