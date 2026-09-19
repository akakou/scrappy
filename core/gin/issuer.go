package scrappy_gin

import (
	"github.com/akakou/mcl_utils"
	"github.com/akakou/scrappy/ecdaa_helper"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func RunIssuer(template string) {
	secret := []byte("secret")
	r := gin.Default()
	store := cookie.NewStore(secret)
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob(template)

	SetupIssuerEndpoints(r)

	r.Run(":8080")
}

func SetupIssuerEndpoints(r *gin.Engine) {
	rng := mcl_utils.InitRandom()
	issuer, err := ecdaa_helper.IssuerFromConfigFile()

	if err != nil {
		panic(err)
	}

	r.GET("/gen_seed", func(c *gin.Context) {
		session := sessions.Default(c)

		seed, issuerB, err := ecdaa_helper.GenSeed(rng)
		if err != nil {
			panic(err)
		}

		session.Set("issuerB", string(issuerB))
		session.Save()

		c.String(200, "%s", string(seed))

	})

	r.POST("/make_cred", func(c *gin.Context) {
		session := sessions.Default(c)

		joinReq := c.PostForm("join_req")
		issuerB := session.Get("issuerB").(string)

		cred, err := ecdaa_helper.MakeCred([]byte(joinReq), []byte(issuerB), issuer, rng)
		if err != nil {
			panic(err)
		}

		c.String(200, "%v", string(cred))
	})

}
