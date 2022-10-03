package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/fast", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_without_attest", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_with_attest", func(c *gin.Context) {
		attestation := c.PostForm("attestation")

		fmt.Printf("attestation=%v\n", attestation)

		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.Run()
}
