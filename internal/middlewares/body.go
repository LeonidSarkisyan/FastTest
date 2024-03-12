package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

func IsNotEmptyRequestBody(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.JSON(422, gin.H{"detail": "Ошибка при прочитывании реквест боди"})
		c.Abort()
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if len(body) == 0 {
		c.JSON(422, gin.H{"detail": "Пустой реквест боди"})
		c.Abort()
		return
	}

	c.Next()
}
