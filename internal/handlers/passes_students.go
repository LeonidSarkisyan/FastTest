package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetResults(c *gin.Context) {
	userID := c.GetInt("userID")

	results, err := h.TestService.GetAllAccessess(userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.JSON(200, responses.NewListResponse(results))
}

func (h *Handler) GetPassesAndStudents(c *gin.Context) {
	userID := c.GetInt("userID")
	resultID := MustID(c, "result_id")

	passes, students, err := h.TestService.GetPassesAndStudents(resultID, userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.JSON(200, gin.H{
		"passes":   passes,
		"students": students,
	})
}

func (h *Handler) GetStartedTest(c *gin.Context) {
	resultID := MustID(c, "result_id")

	var code models.Code

	if err := c.Bind(&code); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	pass, err := h.TestService.GetPass(resultID, code.Code)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	log.Info().Any("pass", pass).Send()
}
