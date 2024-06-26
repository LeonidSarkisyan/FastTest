package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"App/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"path/filepath"
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

func (h *Handler) CreateQuestionWithType(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	type_ := c.Param("type")

	id, data, err := h.QuestionService.CreateWithType(testID, userID, type_)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	text := models.GetTextFromType(type_)

	c.JSON(201, gin.H{
		"id":   id,
		"data": data,
		"type": type_,
		"text": text,
	})
}

func (h *Handler) SaveQuestionWithType(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	type_ := c.Param("type")
	questionID := MustID(c, "question_id")

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	err = h.QuestionService.Save(testID, userID, questionID, type_, data)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		c.Abort()
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
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
			ID:       question.ID,
			Text:     question.Text,
			Data:     question.Data,
			Type:     question.Type,
			ImageURL: question.ImageURL,
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

func (h *Handler) UploadImageQuestion(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	filename := utils.GenerateUniqueFilename(file.Filename)

	err = c.SaveUploadedFile(file, filepath.Join("static/media/questions/", filename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s, err := h.QuestionService.UploadImage(userID, testID, questionID, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": s})
}

func (h *Handler) DeleteImageQuestion(c *gin.Context) {
	userID := c.GetInt("userID")
	testID := MustID(c, "test_id")
	questionID := MustID(c, "question_id")

	err := h.QuestionService.DeleteImage(userID, testID, questionID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
