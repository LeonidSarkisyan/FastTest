package handlers

import (
	"App/internal/middlewares"
	"App/internal/models"
	"App/pkg/systems"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
)

type ResultService interface {
	SaveResult(
		studentID, accessID, passID int, questions, questionsFromUser []models.QuestionWithAnswers,
		access models.AccessOut, timePass int,
	) (models.ResultStudent, error)

	Reset(passID int, access models.AccessOut) error
}

type StudentService interface {
	CreateStudentFromExcel(userID, groupID int, file io.Reader) ([]models.Student, error)
	GetAll(userID, groupID int) ([]models.Student, error)
	Delete(userID, groupID, studentID int) error
	Create(userID, groupID int, student models.Student) (int, error)
	Get(studentID int) (models.Student, error)
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

	GetAllQuestionsWithAnswers(testID int) ([]models.QuestionWithAnswers, error)
}

type TestService interface {
	Create(title string, userID int) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(userID, testID int, title string) error

	CreateAccess(userID, testID, groupID int, accessIn models.Access) (int, error)
	GetAccess(userID, accessID int) (models.AccessOut, error)
	GetAllAccessess(userID int) ([]models.AccessOut, error)

	GetResult(resultID int) (models.AccessOut, error)

	CreatePasses(groupID, accessID int) error
	GetPasses(accessID int) ([]models.Passes, error)
	GetPassByCode(resultID int, code int64) (models.Passes, error)
	GetPassByStudentID(passID, studentID int) (models.Passes, error)
	ClosePass(passID int) error

	GetPassesAndStudents(resultID, userID int) (models.ForResultTable, error)
}

type UserService interface {
	Register(in models.UserIn) error
	Login(in models.UserIn) (string, error)
}

type Channels struct {
	Broadcast         map[int]*chan Message
	BroadcastStudents map[int]*chan Message
	TimesMap          map[int]*chan int
	ResetMap          map[int]*chan int
}

type Handler struct {
	*ClientManager
	*Channels
	UserService
	TestService
	QuestionService
	AnswerService
	GroupService
	StudentService
	ResultService
	Config *systems.AppConfig
}

func NewHandler(
	cm *ClientManager,
	u UserService, t TestService, q QuestionService, a AnswerService, g GroupService, s StudentService,
	r ResultService,
) *Handler {
	return &Handler{ClientManager: cm, UserService: u, TestService: t, QuestionService: q,
		AnswerService: a, GroupService: g, StudentService: s, ResultService: r}
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

	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})

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
		pages.GET("/results", h.ResultPage)
		pages.GET("/results/:result_id", h.OneResultPage)
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

	results := router.Group("/results", middlewares.AuthProtect)
	{
		results.GET("", h.GetResults)
		results.GET("/:result_id", h.GetPassesAndStudents)
		results.GET("/:result_id/ws", h.CreateStreamConnect)
		results.PATCH("/:result_id/reset/:pass_id", h.ResetResult)
	}

	studentsPage := router.Group("/passing")
	{
		studentsPage.GET("/:result_id", h.PassingPage)
		studentsPage.POST("/:result_id", h.GetStartedTest)

		studentsPage.GET("/:result_id/solving/:pass_id", h.IssueTestPage)
		studentsPage.GET("/:result_id/solving/:pass_id/questions", h.GetQuestionsForStudent)
		studentsPage.POST("/:result_id/solving/:pass_id/results", h.CreateResult)
		studentsPage.GET("/:result_id/ws/student/:pass_id", h.CreateWSStudentConnect)
		studentsPage.GET("/abort", h.AbortPage)
	}

	return router
}
