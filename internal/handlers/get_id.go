package handlers

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func MustID(c *gin.Context, key string) int {
	keyIDStr := c.Param(key)

	keyID, err := strconv.Atoi(keyIDStr)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return 0
	}

	return keyID
}
