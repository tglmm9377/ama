package main

import (
	"ama/httpserver"
)

//如果数据库对应asin 数据小于页面实际数据超过10则开启爬取
const STANDDBDATA = 10

func main()  {

	httpserver.StartServer()
}
