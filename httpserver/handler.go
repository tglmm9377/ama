package httpserver

import (
	"ama/db"
	"ama/jobstate"
	"ama/middlevars"
	"ama/robot"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"sync"
)


var Lock sync.WaitGroup
var mutex sync.Mutex

func index(ctx *gin.Context){
	ctx.HTML(200 , "index.html",gin.H{
		"title":"index",
	})
}

func ProcessAsin(ctx *gin.Context){



	middlevars.Site = "www.amazon.ca"
	asin := ctx.PostForm("asin")
	//var url string
	var j jobstate.Job
	db.DB.Sync2(&j)

	mutex.Lock()
	middlevars.Asin = asin
	var reviewdb db.Review

	db.DB.Sync2(&reviewdb)
	mutex.Unlock()

	mutex.Lock()
	has , _ := db.DB.Where("name=?",asin).Get(&j)
	if has{
		status := j.Status
		fmt.Println(status)
		if status == "running"{
			log.Println("已经有进程开始运行")
			// 查询出任务总量 ， 和当前完成情况
			_ , url := GetUrlWithAsin(asin)
			review := GetReviewNum(asin,url)
			////查询当前 数据表中数据条目
			count , _ := db.DB.Table(asin).Count(&reviewdb)
			ctx.HTML(200 , "index.html",gin.H{
				"message":"任务"+asin+"已经开始运行,review总数:"+strconv.Itoa(review)+"当前数据库总数:"+strconv.Itoa(int(count)),
				//"message":"任务"+asin+"已经开始运行,review总数:",

			})
			mutex.Unlock()
			return
		}else{
			j.Name = asin
			j.Status = "running"
			db.DB.Where("name=?",asin).Update(&j)
		}

	}else{
		j.Name = asin
		fmt.Println("状态改为running")
		j.Status = "running"
		db.DB.Insert(&j)
	}
	mutex.Unlock()

	//middlevars.Asin = ctx.PostForm("asin")

	tableExist , err := db.DB.IsTableExist(asin)
	if err != nil{
		log.Println("judge asin or in db error: ",err)
		return
	}
	//表已经存在，判断数据是否完整
	//严格模式 ， 一条数据都不允许丢失
	//非严格模式 STANDDBDATA  = 10
	var result []map[string]string
	var message string
	//判断表是否存在
	if tableExist {
		empty , err := db.DB.IsTableEmpty(asin)
		if !empty {
			fmt.Print("存在查询操作")
			sql := "select color,size,country,count(*) as count from " + asin + " where color != \"\" group by country,color,size"
			result, err = db.DB.QueryString(sql)
			if err != nil {
				message = "数据表" + asin + "查询失败,error: " + fmt.Sprintf("%s", err)
				return
			}
			ctx.HTML(200 , "result.html",gin.H{
				"message":message,
				"result":result,
				"title":"result",
			})
			func() {
				j.Name = asin
				j.Status = "stoped"
				db.DB.Where("name=?",asin).Update(&j)
			}()
		}else{
			ok , url := GetUrlWithAsin(asin)
			if !ok{
				//通过asin查询list失败 ，返回建议，用户通过review地址拉取数据
				ctx.HTML(200 , "index.html",gin.H{
					"message":"无法通过asin查询到相关listing 请使用直接通过填写review地址的方式",

				})
				func() {
					j.Name = asin
					j.Status = "stoped"
					db.DB.Where("name=?",asin).Update(&j)
				}()
				return
			}
			reviewnum := GetReviewNum(asin , url+"1")
			Lock.Add(1)
			go job2(reviewnum , url,asin)
			ctx.HTML(200 ,"index.html",gin.H{
				"message": "查询的数据不存在正在录入， 可能需要一段时间",
			})
		}

	}else{
		//插入数据操作
		ok , url := GetUrlWithAsin(asin)
		if !ok{
			//通过asin查询list失败 ，返回建议，用户通过review地址拉取数据
			ctx.HTML(200 , "index.html",gin.H{
				"message":"无法通过asin查询到相关listing 请使用直接通过填写review的方式",
			})
			return
		}
		fmt.Println(url+"1")
		reviewnum := GetReviewNum(asin , url+"1")
		fmt.Printf("数据表%s不存在，执行job任务",asin)
		Lock.Add(1)
		go job2(reviewnum , url,asin)

		ctx.HTML(200 ,"index.html",gin.H{
			"message": "查询的数据不存在正在录入， 可能需要一段时间",
		})


	}

}


func openrobot(c *gin.Context){
	robot.Init()
	c.JSON(200,gin.H{
		"message":"open chrome ok",
	})
}






func SetUrl(c *gin.Context){
	url := c.Query("url")
	middlevars.Url = url
}


func jobProgress(asin string){

}

//func SearchUrl(c *gin.Context){
//	//获取前段提交的 url
//	//https://www.amazon.ca/store-Rambler-Tumbler-Thermik-Tumblers/product-reviews/B01J5Y1N8Y/ref=cm_cr_arp_d_viewopt_rvwer?ie=UTF8&reviewerType=avp_only_reviews&pageNumber=1
//	url := c.PostForm("searchurl")
//	asin := strings.Split(url,"/")[5]
//	mutex.Lock()
//	middlevars.Asin = asin
//	mutex.Unlock()
//
//	mutex.Lock()
//	CheckJobStatus(asin , c)
//	mutex.Unlock()
//
//	fmt.Println(url , asin)
//	reviewnum := GetReviewNum(asin , url)
//
//
//
//	go job2(reviewnum , asin , url)
//	//通过提交的url提取asin
//	c.HTML(200 ,"index.html",gin.H{
//		"message" :"",
//	})
//}


//func CheckJobStatus(asin string , ctx *gin.Context){
//	var j jobstate.Job
//	has , _ := db.DB.Where("name=?",asin).Get(&j)
//	if has{
//		status := j.Status
//		fmt.Println(status)
//		if status == "running"{
//			log.Println("已经有进程开始运行")
//			// 查询出任务总量 ， 和当前完成情况
//			ctx.HTML(200 , "index.html",gin.H{
//				"message":"任务"+asin+"已经开始运行",
//			})
//			mutex.Unlock()
//			return
//		}else{
//			j.Name = asin
//			j.Status = "running"
//			db.DB.Where("name=?",asin).Update(&j)
//		}
//
//	}else{
//		j.Name = asin
//		fmt.Println("状态改为running")
//		j.Status = "running"
//		db.DB.Insert(&j)
//	}
//}