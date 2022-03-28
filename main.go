package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func tmp(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	server := gin.Default()
	server.LoadHTMLGlob("template/html/*")
	//設定靜態資源的讀取
	server.Static("/assets", "./template/assets")

	server.GET("/", tmp)
	server.Run(":8888")
}
