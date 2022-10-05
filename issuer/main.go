package main

import (
	"github.com/akakou/ecdaa"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	secret := []byte("secret")
	r := gin.Default()

	rng := ecdaa.InitRandom()
	issuer := ecdaa.RandomIssuer(rng)

	err := ecdaa.VerifyIPK(&issuer.Ipk)

	if err != nil {
		panic(err)
	}

	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))

	r.Run()
}
