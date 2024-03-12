package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		log.Info().Msgf(
			"Response - Path: %s, Status: %d, Method: %s", c.Request.URL.Path, c.Writer.Status(), c.Request.Method)
	}
}
