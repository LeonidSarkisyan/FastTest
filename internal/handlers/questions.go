package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func (h *Handler) CreateQuestion(c *gin.Context) {
	userID := c.GetInt("userID")

	testIDStr := c.Param("test_id")

	testID, err := strconv.Atoi(testIDStr)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	id, ids, err := h.QuestionService.Create(testID, userID)

	if err != nil {
		log.Err(err).Send()
		SendErrorResponse(c, 400, err.Error())
	}

	c.JSON(201, gin.H{
		"id":          id,
		"answers_ids": ids,
	})
}

func (h *Handler) GetAllQuestion(c *gin.Context) {
	userID := c.GetInt("userID")

	testIDStr := c.Param("test_id")

	testID, err := strconv.Atoi(testIDStr)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	log.Info().Int("testID", testID).Int("userID", userID).Send()

	questions, err := h.QuestionService.GetAllQuestionsWithAnswers(testID)

	log.Info().Any("questions", questions).Send()

	var questionsWithAnswers []models.QuestionWithAnswersWithOutIsCorrect

	for _, question := range questions {
		questionWithAnswer := models.QuestionWithAnswersWithOutIsCorrect{
			ID:   question.ID,
			Text: question.Text,
		}

		for _, answer := range question.Answers {
			questionWithAnswer.Answers = append(questionWithAnswer.Answers, models.AnswerWithCorrect{
				ID:        answer.ID,
				Text:      answer.Text,
				IsCorrect: answer.IsCorrect,
			})
		}

		questionsWithAnswers = append(questionsWithAnswers, questionWithAnswer)
	}

	if err != nil {
		log.Err(err).Send()
		SendErrorResponse(c, 400, err.Error())
	}

	c.JSON(200, responses.NewListResponse(questionsWithAnswers))
}

func (h *Handler) UpdateQuestion(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	var question models.QuestionUpdate

	if err := c.BindJSON(&question); err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err := h.QuestionService.Update(userID, testID, questionID, question)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (h *Handler) DeleteQuestion(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	err := h.QuestionService.Delete(userID, testID, questionID)

	if err != nil {
		SendErrorResponse(c, 400, err.Error())
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
