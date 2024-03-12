package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Detail string `json:"detail"`
}

func SendErrorResponse(c *gin.Context, statusCode int, detail string) {
	c.AbortWithStatusJSON(statusCode, detail)
}
