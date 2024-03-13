package handlers

import (
	"App/internal/middlewares"
	"App/internal/models"
	"github.com/gin-gonic/gin"
)

type AnswerService interface {
	Create(userID, testID, questionID int) (int, error)
	CreateThree(userID, testID, questionID int) error
	GetAllByQuestionID(userID, testID, questionID int) ([]models.Answer, error)
	Update(userID, testID, questionID, answerID int, answerUpdate models.AnswerUpdate) error
	Delete(userID, testID, questionID, answerID int) error
}

type QuestionService interface {
	Create(testID, userID int) (int, error)
	GetAll(testID, userID int) ([]models.Question, error)
	Update(userID, testID, questionID int, question models.QuestionUpdate) error
	Delete(userID, testID, questionID int) error
}

type TestService interface {
	Create(title string, userID int) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(userID, testID int, title string) error
}

type UserService interface {
	Register(in models.UserIn) error
	Login(in models.UserIn) (string, error)
}

type Handler struct {
	UserService
	TestService
	QuestionService
	AnswerService
}

func NewHandler(u UserService, t TestService, q QuestionService, a AnswerService) *Handler {
	return &Handler{u, t, q, a}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(middlewares.LoggerMiddleware())

	router.Use(func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Next()
	})

	router.LoadHTMLGlob("templates/*")

	router.Static("/static", "./static")

	router.GET("/", h.MainPage)
	router.GET("/auth", h.AuthPage)

	auth := router.Group("/auth", middlewares.IsNotEmptyRequestBody)
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	withAuth := router.Group("/p", middlewares.AuthProtect)
	{
		withAuth.GET("/header", h.Header)
		withAuth.GET("/tests", h.TestPage)
		withAuth.GET("/tests/:test_id", h.OneTestPage)
	}

	tests := router.Group("/tests", middlewares.AuthProtect)
	{
		tests.POST("", h.CreateTest)
		tests.GET("", h.GetAll)
		tests.GET("/:test_id", h.GetTest)
		tests.PATCH("/:test_id", h.UpdateTestTitle)

		questions := tests.Group("/:test_id/questions")
		{
			questions.POST("", h.CreateQuestion)
			questions.GET("", h.GetAllQuestion)
			questions.PATCH("/:question_id", h.UpdateQuestion)
			questions.DELETE("/:question_id", h.DeleteQuestion)

			answers := questions.Group("/:question_id/answers")
			{
				answers.POST("", h.CreateAnswer)
				answers.GET("", h.GetAnswers)
				answers.PATCH("/:answer_id", h.UpdateAnswer)
				answers.DELETE("/:answer_id", h.DeleteAnswer)
			}
		}
	}

	return router
}
