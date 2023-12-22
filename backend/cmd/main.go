package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ibook/backend/internal/conf"
	"net/http"
)

var IBookConfig *conf.Config

func main() {
	fmt.Println(conf.GetConf())
	server := gin.Default()
	server.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"data1": "this is my data",
		})
	})
	err := server.Run(conf.GetConf().ServerConf.Port)
	if err != nil {
		return
	}
}
