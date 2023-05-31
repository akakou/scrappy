//go:build example
// +build example

package scrappy

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const LOOP_NUM = 1000000000

const HOST_NAME = "http://localhost:8080"

func somethingHeavy() int {
	num := 0

	for i := 0; i < LOOP_NUM; i++ {
		num *= i
		num %= LOOP_NUM
	}

	return num
}

func ServVerifier() {
	r := ginServ()

	r.GET("/", func(c *gin.Context) {
		period := Now()

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

	r.POST("/slow_with_attest", VerifyMiddleware(HOST_NAME), func(c *gin.Context) {
		fmt.Printf("%v", somethingHeavy())
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.Run()
}
