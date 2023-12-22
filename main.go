package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	myMap := map[string]string{}
	for range myMap {
	}
	server := gin.Default()
	server.GET("/test", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"data1": "this is my data",
		})
	})
	err := server.Run(":6000")
	if err != nil {
		return
	}
}
