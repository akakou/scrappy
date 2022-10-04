package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/akakou/ecdaa"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Cred   *ecdaa.MiddleEncodedCredential
	Isk    *ecdaa.MiddleEncodedISK
	Ipk    *ecdaa.MiddleEncodedIPK
	Handle []byte
}

func main() {
	var config Config
	configBuf, err := ioutil.ReadFile("../init/config.json ")

	if err != nil {
		log.Fatal("error:", err)
	}

	err = json.Unmarshal(configBuf, &config)

	if err != nil {
		log.Fatal("error:", err)
	}

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
		var signature ecdaa.MiddleEncodedSignature
		session := sessions.Default(c)
		period := session.Get("period").(string)

		attestation := c.PostForm("attestation")

		attestBuf, err := base64.StdEncoding.DecodeString(attestation)
		if err != nil {
			log.Fatal("error:", err)
			c.HTML(http.StatusOK, "err.html", gin.H{})
		}

		err = json.Unmarshal(attestBuf, &signature)

		if err != nil {
			log.Fatal("error:", err)
		}

		fmt.Printf("attestation=%v\n", attestation)
		fmt.Printf("period=%v\n", period)

		err = ecdaa.Verify(
			[]byte{},
			[]byte(period),
			signature.Decode(),
			config.Ipk.Decode(),
		)

		if err != nil {
			c.HTML(http.StatusOK, "err.html", gin.H{})
		}

		c.HTML(http.StatusOK, "hello.html", gin.H{})
	})

	r.Run()
}
