package handlers

import (
	"App/internal/handlers/responses"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
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

	forResultTable, err := h.TestService.GetPassesAndStudents(resultID, userID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	c.JSON(200, forResultTable)
}

func (h *Handler) GetStartedTest(c *gin.Context) {
	resultID := MustID(c, "result_id")

	var code models.Code

	if err := c.Bind(&code); err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	result, err := h.TestService.GetResult(resultID)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	pass, err := h.TestService.GetPassByCode(result.ID, code.Code)

	if err != nil {
		SendErrorResponse(c, 422, err.Error())
		return
	}

	if pass.IsActivated {
		SendErrorResponse(c, 403, "все попытки закончились")
		c.Abort()
		return
	}

	log.Info().Any("pass", pass).Send()

	c.SetCookie(
		"StudentID", strconv.Itoa(pass.StudentID), 2592000, "/",
		Domain, false, false,
	)

	c.AbortWithStatusJSON(200, gin.H{
		"id": pass.ID,
	})
}

func (h *Handler) IssueTestPage(c *gin.Context) {
	studentIDCookie, err := c.Cookie("StudentID")

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 401, доступ запрещён",
		})
		c.Abort()
		return
	}

	studentID, err := strconv.Atoi(studentIDCookie)

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 401, доступ запрещён",
		})
		c.Abort()
		return
	}

	passID := MustID(c, "pass_id")

	pass, err := h.TestService.GetPassByStudentID(passID, studentID)

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 404, неизвестный тест",
		})
		c.Abort()
		return
	}

	student, err := h.StudentService.Get(studentID)

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 404, неизвестный студент",
		})
		c.Abort()
		return
	}

	if pass.IsActivated {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 409, попытки закончились",
		})
		c.Abort()
		return
	}

	resultID := MustID(c, "result_id")

	result, err := h.TestService.GetResult(resultID)

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 404, неизвестный тест",
		})
		c.Abort()
		return
	}

	test, err := h.TestService.Get(result.TestID, result.UserID)

	if err != nil {
		c.HTML(401, "error.html", gin.H{
			"error": "Ошибка 404, неизвестный тест",
		})
		c.Abort()
		return
	}

	c.HTML(200, "solve.html", gin.H{
		"title":   "Прохождение теста",
		"access":  result,
		"test":    test,
		"student": student,
	})
}
