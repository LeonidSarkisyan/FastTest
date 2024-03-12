package middlewares

import (
	"App/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
)

func AuthProtect(c *gin.Context) {
	token, err := c.Cookie("Authorization")

	if err != nil {
		log.Err(err).Send()
		c.Redirect(http.StatusPermanentRedirect, "/auth")
		c.Abort()
		return
	}

	userID, err := utils.GetUserIDFromToken(token)

	if err != nil {
		log.Err(err).Send()
		c.Redirect(http.StatusPermanentRedirect, "/auth")
		c.Abort()
		return
	}

	c.Set("userID", int(userID))
	c.Next()
}
