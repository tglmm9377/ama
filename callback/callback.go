package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

type cookie http.Cookie


var wait sync.WaitGroup
func main() {

	engine := gin.Default()

	engine.GET("/user/:id", user)
	engine.GET("/user",userquery)

	engine.Run("0.0.0.0:8080")

}

func user(c *gin.Context){
	id := c.Param("id")
	fmt.Println(id)
	c.JSON(200 , gin.H{
		"message":"okokok",
		"id":id,
	})
}

func userquery(c *gin.Context)  {
	username := c.DefaultQuery("username","clf")
	m := c.QueryMap("username")
	fmt.Println(m)
	c.JSON(200 , gin.H{
		"username":username,
	})
}
