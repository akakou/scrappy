package scrappy_gin

import (
	"net/http"

	"github.com/akakou/scrappy"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func VerifyMiddleware(HOST_NAME string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		period := session.Get("period").(int)

		attestation := c.PostForm("scrappy")

		err := scrappy.Verify(attestation, HOST_NAME, int(period))

		if err != nil {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})

			c.Abort()
		}

		c.Next()
	}
}
