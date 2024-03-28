package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strings"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !strings.HasSuffix(c.Request.URL.Path, ".css") &&
			!strings.HasSuffix(c.Request.URL.Path, ".js") &&
			!strings.HasSuffix(c.Request.URL.Path, ".png") &&
			!strings.HasSuffix(c.Request.URL.Path, ".PNG") {
			log.Info().Msgf(
				"Response - Path: %s, Status: %d, Method: %s", c.Request.URL.Path, c.Writer.Status(), c.Request.Method)
		}
	}
}
