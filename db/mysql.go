package db

import (
	"ama/middlevars"
	"ama/utils/file"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
	"xorm.io/xorm"
)
type MysqlInfo struct{
	Address string `yaml:"address"`
	Port int `yaml:"port"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	DB string `yaml:"db"`
}

type Review struct{
	Id string `xorm:"id pk varchar(20)" json:"id"`
	Star string `xorm:"star" json:"star"`
	Color string `xorm:"color" json:"color"`
	Size string `xorm:"size" json:"size"`
	Country string `xorm:"country" json:"country"`
	Date time.Time `xorm:"date" json:"date"`
}
func (r *Review)TableName() string{
	return middlevars.Asin
}

var info map[string]MysqlInfo

var DB *xorm.Engine

//初始化数据库连接
func init(){
	log.Println("start init mysql ...")
	getInfo("etc/mysql.yml")
	var err error
	DB , err = xorm.NewEngine("mysql",info["mysql"].User+":"+info["mysql"].Password+"@/"+info["mysql"].DB+"?"+"charset=utf8")
	fmt.Println("info:",info["mysql"])
	if err != nil{
		log.Println("init mysql engine faild: ",err)
		os.Exit(3)

	}
	log.Println("init mysql success!🎉 ")

}

func getInfo(path string)  {
	exist , err := file.IsExist(path)
	if !exist{
		fmt.Println("mysql conf file not exist: ",err)
		return
	}
	var mysqlinfo  = new(map[string]MysqlInfo)
	file.ReadYaml(path , mysqlinfo)
	if err != nil{
		log.Println("read mysql config file error: ",err)
	}
	info = *mysqlinfo
}



func InsertReview(asin string , review *Review){
	DB.Table(asin)
	fmt.Println("调用sync2")

	DB.TableName(asin)
	//DB.Sync2(review)
	DB.Table(asin).Insert(review)
}