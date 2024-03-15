package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
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

	log.Info().Any("tests", tests).Send()

	c.JSON(200, responses.NewListResponse(tests))
}

func (h *Handler) GetTest(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")

	test, err := h.TestService.Get(testID, userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.JSON(200, test)
}

func (h *Handler) UpdateTestTitle(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")

	var testUpdate models.TestUpdate

	if err := c.BindJSON(&testUpdate); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err := h.TestService.UpdateTitle(userID, testID, testUpdate.Title)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, gin.H{})
}
