package handlers

import (
	"App/internal/ai"
	"App/internal/handlers/responses"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateQuestionsFromGPT(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")

	var promptParams ai.PromptParams

	if err := c.BindJSON(&promptParams); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	questions, err := h.AiService.CreateQuestionsFromGPT(userID, testID, promptParams)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(201, responses.NewListResponse(questions))
}
