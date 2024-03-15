package handlers

import (
	"App/internal/middlewares"
	"App/internal/models"
	"github.com/gin-gonic/gin"
	"io"
)

type StudentService interface {
	CreateStudentFromExcel(userID, groupID int, file io.Reader) ([]models.Student, error)
	GetAll(userID, groupID int) ([]models.Student, error)
	Delete(userID, groupID, studentID int) error
	Create(userID, groupID int, student models.Student) (int, error)
}

type GroupService interface {
	Create(in models.GroupIn, userID int) (int, error)
	GetAll(userID int) ([]models.GroupOut, error)
	Get(groupID, userID int) (models.GroupOut, error)
	UpdateTitle(userID, groupID int, title string) error
}

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

	CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error)
	CreatePasses(groupID, accessID int) error
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
	GroupService
	StudentService
}

func NewHandler(
	u UserService, t TestService, q QuestionService, a AnswerService, g GroupService, s StudentService,
) *Handler {
	return &Handler{u, t, q, a, g, s}
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

	pages := router.Group("/p", middlewares.AuthProtect)
	{
		pages.GET("/header", h.Header)
		pages.GET("/tests", h.TestPage)
		pages.GET("/tests/:test_id/access", h.OneTestAccessPage)
		pages.GET("/tests/:test_id", h.OneTestPage)
		pages.GET("/groups", h.GroupPage)
		pages.GET("/groups/:group_id", h.OneGroupPage)
	}

	tests := router.Group("/tests", middlewares.AuthProtect)
	{
		tests.POST("", h.CreateTest)
		tests.GET("", h.GetAll)
		tests.POST("/:test_id/access/:group_id", h.CreateAccess)
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

	groups := router.Group("/groups", middlewares.AuthProtect)
	{
		groups.POST("", h.CreateGroup)
		groups.GET("", h.GetGroups)
		groups.GET("/:group_id", h.GetGroup)
		groups.PATCH("/:group_id", h.UpdateGroupTitle)

		students := groups.Group("/:group_id/students")
		{
			students.POST("", h.CreateStudent)
			students.POST("/excel", h.AddStudentsFromExcel)
			students.GET("", h.GetAllStudents)
			students.DELETE("/:student_id", h.DeleteStudent)
		}
	}

	return router
}
