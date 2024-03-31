package main

import (
	"App/internal/email"
	"App/internal/handlers"
	"App/internal/models"
	"App/internal/repository"
	"App/internal/service"
	"App/pkg/server"
	"App/pkg/systems"
	"context"
	"github.com/rs/zerolog/log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	systems.SetupLogger()
	cfg := systems.MustConfig()
	conn := systems.MustConn(cfg)

	userRepo := repository.NewUserPostgres(conn)
	testRepo := repository.NewTestPostgres(conn)
	questionRepo := repository.NewQuestionPostgres(conn)
	answerRepo := repository.NewAnswerPostgres(conn)
	groupRepo := repository.NewGroupPostgres(conn)
	studentRepo := repository.NewStudentPostgres(conn)
	resultRepo := repository.NewResultPostgres(conn)

	userService := service.NewUserService(userRepo)
	answerService := service.NewAnswerService(answerRepo, testRepo, questionRepo)
	questionService := service.NewQuestionService(questionRepo, testRepo, answerService)
	groupService := service.NewGroupService(groupRepo)
	testService := service.NewTestService(testRepo, studentRepo, questionService, groupService, resultRepo)
	studentService := service.NewStudentService(studentRepo, groupRepo)
	resultService := service.NewResultService(resultRepo)
	aiService := service.NewAiService(questionRepo, testRepo)

	manager := &handlers.ClientManager{
		Clients:   make([]*handlers.Client, 0),
		Broadcast: make(chan handlers.Message),
		TimesMap:  make(map[int]chan int),
		ResetMap:  make(map[int]chan int),
	}

	auth := smtp.PlainAuth("", os.Getenv("SMTP_EMAIL"), os.Getenv("SMTP_PASSWORD"), cfg.SMTP.Host)
	emailClient := email.NewEmailClient(auth, cfg)

	handler := handlers.NewHandler(
		manager, userService, testService, questionService, answerService, groupService, studentService, resultService,
		aiService,
	)
	handler.Config = cfg
	handler.Channels = &handlers.Channels{
		Broadcast:         make(map[int]chan handlers.Message),
		BroadcastStudents: make(map[int]chan handlers.Message),
		TimesMap:          make(map[int]*chan int),
		ResetMap:          make(map[int]*chan int),
	}

	handler.EmailClient = emailClient
	handler.EmailCodeMap = make(map[int64]models.UserIn)
	handler.ResetPasswordMap = make(map[int64]models.User)

	server_ := new(server.Server)

	go func() {
		if err := server_.Run(cfg.Port, handler.InitRoutes()); err != nil {
			log.Fatal().Err(err).Msg("ошибка при запуске сервера")
		}
	}()

	log.Printf("MakeTest старует на порту: %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT)
	<-quit

	log.Print("Сервер остановил свою работу")

	if err := server_.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("ошибка при остановке сервера")
	}

	if err := conn.Close(); err != nil {
		log.Error().Err(err).Msg("ошибка при закрытии соединения с БД")
	}
}
