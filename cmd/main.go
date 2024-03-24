package main

import (
	"App/internal/handlers"
	"App/internal/repository"
	"App/internal/service"
	"App/pkg/server"
	"App/pkg/systems"
	"context"
	"github.com/rs/zerolog/log"
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

	manager := &handlers.ClientManager{
		Clients:    make(map[*handlers.Client]bool),
		Broadcast:  make(chan handlers.Message, 100),
		Register:   make(chan *handlers.Client),
		Unregister: make(chan *handlers.Client),
		TimesMap:   make(map[int]*chan int),
		ResetMap:   make(map[int]*chan int),
	}

	go manager.Start()

	handler := handlers.NewHandler(
		manager, userService, testService, questionService, answerService, groupService, studentService, resultService,
	)
	handler.Config = cfg
	handler.Channels = &handlers.Channels{
		Broadcast:         make(map[int]chan handlers.Message),
		BroadcastStudents: make(map[int]chan handlers.Message),
		TimesMap:          make(map[int]*chan int),
		ResetMap:          make(map[int]*chan int),
	}

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
