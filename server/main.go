package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	secret := []byte("secret")
	r := gin.Default()

	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		period := time.Now()
		period = time.Date(period.Year(), period.Month(), period.Day(), period.Hour(), period.Minute(), 0, 0, time.UTC)

		session := sessions.Default(c)

		session.Set("period", period.Unix())
		session.Save()

		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/fast", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_without_attest", func(c *gin.Context) {
		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.POST("/slow_with_attest", func(c *gin.Context) {
		session := sessions.Default(c)
		period := session.Get("period")

		attestation := c.PostForm("attestation")

		fmt.Printf("attestation=%v\n", attestation)
		fmt.Printf("period=%v\n", period)

		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.Run()
}
