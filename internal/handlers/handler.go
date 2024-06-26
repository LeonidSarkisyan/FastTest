package handlers

import (
	"App/internal/ai"
	"App/internal/email"
	"App/internal/middlewares"
	"App/internal/models"
	"App/pkg/systems"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
)

type AiService interface {
	CreateQuestionsFromGPT(userID, testID int, params ai.PromptParams) ([]models.QuestionWithAnswersWithOutIsCorrect, error)
}

type ResultService interface {
	SaveResult(
		studentID, accessID, passID int, questionsFromUser []models.QuestionWithAnswers,
		access models.AccessOut, timePass int,
	) (models.ResultStudent, error)

	Reset(passID int, access models.AccessOut) error
	GetResultByPassID(passID int) (models.ResultStudent, error)
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
	GetIfNotDelete(groupID, userID int) (models.GroupOut, error)
	UpdateTitle(userID, groupID int, title string) error
	Delete(userID, groupID int) error
}

type AnswerService interface {
	Create(userID, testID, questionID int) (int, error)
	CreateThree(userID, testID, questionID int) ([]int, error)
	GetAllByQuestionID(userID, testID, questionID int) ([]models.Answer, error)
	Update(userID, testID, questionID, answerID int, answerUpdate models.AnswerUpdate) error
	Delete(userID, testID, questionID, answerID int) error
}

type QuestionService interface {
	Create(testID, userID int) (int, []int, error)
	CreateWithType(testID, userID int, type_ string) (int, any, error)
	Save(testID, userID, questionID int, type_ string, data []byte) error
	GetAll(testID, userID int) ([]models.Question, error)
	Update(userID, testID, questionID int, question models.QuestionUpdate) error
	Delete(userID, testID, questionID int) error

	UploadImage(userID, testID, questionID int, filename string) (string, error)
	DeleteImage(userID, testID, questionID int) error

	GetAllQuestionsWithAnswers(testID int) ([]models.QuestionWithAnswers, error)
}

type TestService interface {
	Create(title string, userID int) (int, error)
	Get(testID, userID int) (models.TestOut, error)
	GetIfNotDelete(testID, userID int) (models.TestOut, error)
	GetAll(userID int) ([]models.TestOut, error)
	UpdateTitle(userID, testID int, title string) error
	Delete(userID, testID int) error

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
	GetByEmail(email string) (models.User, error)
	GetByID(userID int) (models.User, error)
	ChangePassword(userID int, newPassword models.NewPassword) error

	Exist(email string) bool
}

type Channels struct {
	Broadcast         map[int]chan Message
	BroadcastStudents map[int]chan Message
	TimesMap          map[int]*chan int
	ResetMap          map[int]*chan int
}

type Handler struct {
	*email.EmailClient
	EmailCodeMap     map[int64]models.UserIn
	ResetPasswordMap map[int64]models.User
	*ClientManager
	*Channels
	AiService
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
	r ResultService, ai AiService,
) *Handler {
	return &Handler{ClientManager: cm, UserService: u, TestService: t, QuestionService: q,
		AnswerService: a, GroupService: g, StudentService: s, ResultService: r, AiService: ai}
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
		auth.POST("/register", h.CreateAccount)
		auth.POST("/login", h.Login)

		// router - специально
		router.GET("/auth/confirm/:code", h.ConfirmAccount)
		router.POST("/auth/reset", middlewares.AuthProtect, h.CreateResetPasswordCode)
		router.POST("/auth/change/:code", middlewares.AuthProtect, h.ResetPassword)
		router.GET("/auth/logout", h.LogOut)
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
		pages.GET("/account", h.Account)
	}

	tests := router.Group("/tests", middlewares.AuthProtect)
	{
		tests.POST("", h.CreateTest)
		tests.GET("", h.GetAll)
		tests.POST("/:test_id/access/:group_id", h.CreateAccess)
		tests.GET("/:test_id", h.GetTest)
		tests.PATCH("/:test_id", h.UpdateTestTitle)
		tests.DELETE("/:test_id", h.DeleteTest)

		questions := tests.Group("/:test_id/questions")
		{
			questions.POST("", h.CreateQuestion)
			questions.POST("/chat-gpt", h.CreateQuestionsFromGPT)
			questions.POST("/type/:type", h.CreateQuestionWithType)
			questions.PATCH("/type/:type/:question_id", h.SaveQuestionWithType)
			questions.GET("", h.GetAllQuestion)
			questions.PATCH("/:question_id", h.UpdateQuestion)
			questions.DELETE("/:question_id", h.DeleteQuestion)
			questions.POST("/:question_id/image", h.UploadImageQuestion)
			questions.DELETE("/:question_id/image", h.DeleteImageQuestion)

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
		groups.DELETE("/:group_id", h.DeleteGroup)

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

	router.GET("/news", h.NewsPage)

	return router
}
