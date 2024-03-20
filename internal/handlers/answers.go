package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) CreateAnswer(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	id, err := h.AnswerService.Create(userID, testID, questionID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.JSON(201, gin.H{
		"id": id,
	})
}

func (h *Handler) GetAnswers(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	answers, err := h.AnswerService.GetAllByQuestionID(userID, testID, questionID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	var answers_ []models.AnswerWithCorrect

	for _, a := range answers {
		answers_ = append(answers_, models.AnswerWithCorrect{
			ID:        a.ID,
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		})
	}

	c.JSON(200, responses.NewListResponse(answers_))
}

func (h *Handler) UpdateAnswer(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")
	answerID := MustID(c, "answer_id")

	var answerUpdate models.AnswerUpdate

	if err := c.BindJSON(&answerUpdate); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err := h.AnswerService.Update(userID, testID, questionID, answerID, answerUpdate)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) DeleteAnswer(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")
	answerID := MustID(c, "answer_id")

	err := h.AnswerService.Delete(userID, testID, questionID, answerID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.AbortWithStatusJSON(http.StatusNoContent, gin.H{})
}
