package handlers

import (
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const Domain = "https://фаст-тест.рф/"
const DomainWithPort = "https://фаст-тест.рф/"

func (h *Handler) Register(c *gin.Context) {
	var in models.UserIn

	if err := c.BindJSON(&in); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	err := h.UserService.Register(in)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, "")
}

func (h *Handler) Login(c *gin.Context) {
	var in models.UserIn

	if err := c.BindJSON(&in); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	token, err := h.UserService.Login(in)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.SetCookie("Authorization", token, 2592000, "/", Domain, false, true)
	c.AbortWithStatusJSON(http.StatusNoContent, "")
}
