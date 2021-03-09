package httpserver

import (
	"github.com/gin-gonic/gin"
	"net"
)

func StartServer(){
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.LoadHTMLGlob("templates/*")
	Config(r)
	listener,_ := net.Listen("tcp","localhost:8080")
	r.RunListener(listener)
	r.Run()
}
