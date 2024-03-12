package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) CreateTest(c *gin.Context) {
	userID := c.GetInt("userID")

	log.Info().Int("userID", userID).Send()

	var testIn models.TestIn

	if err := c.BindJSON(&testIn); err != nil {
		SendErrorResponse(c, 422, err.Error())
	}

	id, err := h.TestService.Create(testIn.Title, userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}

func (h *Handler) GetAll(c *gin.Context) {
	userID := c.GetInt("userID")

	tests, err := h.TestService.GetAll(userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
	}

	c.JSON(200, responses.NewListResponse(tests))
}
