package main

import (
	"core"
	"fmt"
	"net/http"

	"github.com/akakou/ecdaa"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const LOOP_NUM = 1000000000

type Config struct {
	Cred   *ecdaa.MiddleEncodedCredential
	Isk    *ecdaa.MiddleEncodedISK
	Ipk    *ecdaa.MiddleEncodedIPK
	Handle []byte
}

const HOST_NAME = "http://localhost:8080"

func somethingHeavy() int {
	num := 0

	for i := 0; i < LOOP_NUM; i++ {
		num *= i
		num %= LOOP_NUM
	}

	return num
}

func main() {
	secret := []byte("secret")
	r := gin.Default()

	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*.html")

	r.GET("/", func(c *gin.Context) {
		period := core.Now()

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

	r.POST("/slow_with_attest", func(c *gin.Context) {
		session := sessions.Default(c)
		period := session.Get("period").(int)

		attestation := c.PostForm("attestation")

		err := core.Verify(attestation, HOST_NAME, int(period))

		if err != nil {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
		} else {
			fmt.Printf("%v", somethingHeavy())
			c.HTML(http.StatusOK, "hello.html", gin.H{})
		}
	})

	r.Run()
}
