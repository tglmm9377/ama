package httpserver

import "github.com/gin-gonic/gin"

func Config(engine *gin.Engine){
	engine.GET("/index",index)
	engine.POST("/postdata",ProcessAsin)
	engine.GET("/openrobot",openrobot)
	//engine.POST("/searchurl",SearchUrl)
}
