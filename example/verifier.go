package main

import (
	"fmt"
	"net/http"

	"github.com/akakou/scrappy"
	scrappy_gin "github.com/akakou/scrappy/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const LOOP_NUM = 1000000000

const HOST_NAME = "http://server:8081"

// const HOST_NAME = "http://192.168.10.105:8081"

func somethingHeavy() int {
	num := 0

	for i := 0; i < LOOP_NUM; i++ {
		num *= i
		num %= LOOP_NUM
	}

	return num
}

func runServer() {
	secret := []byte("secret")
	r := gin.Default()
	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob("../core/gin/templates/*.html")

	r.GET("/", func(c *gin.Context) {
		period := scrappy.Now()

		session := sessions.Default(c)

		session.Set("period", period)
		session.Save()

		c.HTML(http.StatusOK, "index.html", gin.H{
			"period": period,
		})
	})

	r.GET("/fast", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_without_attest", func(c *gin.Context) {
		fmt.Printf("%v", somethingHeavy())

		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_with_attest", scrappy_gin.VerifyMiddleware(HOST_NAME), func(c *gin.Context) {
		fmt.Printf("%v", somethingHeavy())
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.GET("/callback", func(c *gin.Context) {
		c.HTML(http.StatusOK, "callback.html", gin.H{})
	})

	r.Run(":8081")
}
