package robot

import (
	"ama/db"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

const (
	January = 1
	February = 2
	March = 3
	April = 4
	May = 5
	June = 6
	July = 7
	August = 8
	September = 9
	October = 10
	November = 11
	December = 12
)
	var dateMap = map[string]int{

	"January" :1,
	"February" :2,
	"March" : 3,
	"April" : 4,
	"May" : 5,
	"June" : 6,
	"July" : 7,
	"August" : 8,
	"September" : 9,
	"October" : 10,
	"November" : 11,
	"December" : 12,
}
func TransformDate(date string)string{
	datestr := strconv.Itoa(dateMap[date])
	if len(datestr) == 1{
		datestr = "0"+datestr
	}
	return datestr
}

func TranslateDay(day string)string{
	daystr := strings.Split(day , ",")[0]
	//fmt.Println("day:",daystr)
	if len(daystr) == 1{
		daystr = "0" + daystr
	}
	return daystr
}

//R2N6YPK5RIHB9 5.0 out of 5 stars Color Name: TealSize: 50x60 Reviewed in Canada on December 14, 2020
//func ParseData(page int , url string,asin string){
func ParseData(data [][]string,asin string){
	//data := GetPageContent(page , url)
	//fmt.Println(data)
	fmt.Println("开始格式化数据")
	for _ , v := range data{
		id := strings.TrimSpace(v[0])
		star := strings.TrimSpace(strings.Split(v[1],"out")[0])
		sizeAndcolor := v[2]
		size := ""
		color := ""
		if len(sizeAndcolor) > 0{
			fmt.Println(sizeAndcolor)
			if strings.Contains(sizeAndcolor,"Size") {
				color = strings.TrimSpace(strings.Split(strings.Split(sizeAndcolor, ":")[1], "Size")[0])
				size = strings.TrimSpace(strings.Split(sizeAndcolor, "Size:")[1])
			}else{
				color = strings.TrimSpace((strings.Split(sizeAndcolor, ":")[1]))
				size = ""
			}
		}else{
			size = ""
			color = ""
		}
		countryAndDate := v[3]
		country := strings.TrimSpace(strings.ReplaceAll(strings.TrimSpace(strings.Split(countryAndDate,"on")[0]),"Reviewed in",""))
		dateSpilt := strings.Split(strings.TrimSpace(strings.Split(countryAndDate,"on")[1])," ")
		day := TranslateDay(dateSpilt[1])
		date :=TransformDate(strings.TrimSpace(dateSpilt[0]))
		//fmt.Println(id ,star , color , size , country,dateSpilt[2]+"-"+date+"-"+day)
		time ,_ := time.Parse("2006-01-02 15:04:05",dateSpilt[2]+"-"+date+"-"+day+" 00:00:00")
		review := db.Review{id , star,color,size , country,time}
		fmt.Println("插入数据库的操作"	,review)
		db.InsertReview(asin , &review)
	}


}

func ParseHtml(resp string)[][]string{

	fmt.Println("开始解析浏览器数据")
	data := [][]string{}
	pattern := "div[data-hook=review]"
	dom , _ := goquery.NewDocumentFromReader(strings.NewReader(resp))
	dom.Find(pattern).Each(func(i int, selection *goquery.Selection) {
		star := selection.Find("span[class=a-icon-alt]").Text()
		dateAndCountry := selection.Find("span[data-hook=review-date]").Text()
		color := selection.Find("a[data-hook=format-strip]").Text()
		col := strings.ReplaceAll(color, "\n", "")
		_ = selection.Find("span[class=a-profile-name]").Text()
		id, _ := selection.Attr("id")
		data = append(data, []string{id, star, col, dateAndCountry})

	})
	return data
}