package scrappy

import (
	"github.com/akakou/ecdaa"
	"github.com/akakou/scrappy/crypto"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func ServIssuer() {
	r := ginServ()
	rng := ecdaa.InitRandom()
	issuer, err := crypto.IssuerFromConfigFile()

	if err != nil {
		panic(err)
	}

	r.GET("/gen_seed", func(c *gin.Context) {
		session := sessions.Default(c)

		seed, issuerB, err := crypto.GenSeed(rng)
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

		cred, err := crypto.MakeCred([]byte(joinReq), []byte(issuerB), issuer, rng)
		if err != nil {
			panic(err)
		}

		c.String(200, "%v", string(cred))
	})

	r.Run()
}
