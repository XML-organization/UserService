package startup

import (
	"fmt"
	"log"
	"net"
	"user_service/handler"
	"user_service/repository"
	"user_service/service"
	"user_service/startup/config"

	user "github.com/XML-organization/common/proto/user_service"
	saga "github.com/XML-organization/common/saga/messaging"
	"github.com/XML-organization/common/saga/messaging/nats"
	"github.com/neo4j/neo4j-go-driver/neo4j"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

const (
	QueueGroup = "user_service"
)

func (server *Server) Start() {
	postgresClient := server.initPostgresClient()
	neo4jDriver := server.initNeo4jDriver()
	userRepo := server.initUserRepository(postgresClient)
	notRepo := server.initNotificationRepository(postgresClient)
	userNotRepo := server.initUserNotificationRepository(postgresClient)
	userNeo4jRepo := server.initUserNeo4jRepository(neo4jDriver)
	ratingRepo := server.initRatingRepository(postgresClient)

	//change password orchestrator
	commandPublisher := server.initPublisher(server.config.ChangePasswordCommandSubject)
	replySubscriber := server.initSubscriber(server.config.ChangePasswordReplySubject, QueueGroup)
	changePasswordOrchestrator := server.initChangePasswordOrchestrator(commandPublisher, replySubscriber)

	//update user orchestrator
	commandPublisher1 := server.initPublisher(server.config.UpdateUserCommandSubject)
	replySubscriber1 := server.initSubscriber(server.config.UpdateUserReplySubject, QueueGroup)
	updateUserOrchestrator := server.initUpdateUserOrchestrator(commandPublisher1, replySubscriber1)

	userService := server.initUserService(userRepo, notRepo,userNotRepo, userNeo4jRepo, changePasswordOrchestrator, updateUserOrchestrator)
	ratingService := server.initRatingService(ratingRepo)

	//update user
	commandSubscriber2 := server.initSubscriber(server.config.UpdateUserCommandSubject, QueueGroup)
	replyPublisher2 := server.initPublisher(server.config.UpdateUserReplySubject)
	server.initUpdateUserHandler(userService, replyPublisher2, commandSubscriber2)

	//change password
	commandSubscriber1 := server.initSubscriber(server.config.ChangePasswordCommandSubject, QueueGroup)
	replyPublisher1 := server.initPublisher(server.config.ChangePasswordReplySubject)
	server.initCreateUserHandler(userService, replyPublisher1, commandSubscriber1)

	//create user
	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(userService, replyPublisher, commandSubscriber)

	userHandler := server.initUserHandler(userService, ratingService)

	server.startGrpcServer(userHandler)
}

func (server *Server) initUpdateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *service.UpdateUserOrchestrator {
	orchestrator, err := service.NewUpdateUserOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initChangePasswordOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) *service.ChangePasswordOrchestrator {
	orchestrator, err := service.NewChangePasswordOrchestrator(publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
	return orchestrator
}

func (server *Server) initPublisher(subject string) saga.Publisher {
	publisher, err := nats.NewNATSPublisher(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject)
	if err != nil {
		log.Fatal(err)
	}
	return publisher
}

func (server *Server) initSubscriber(subject, queueGroup string) saga.Subscriber {
	subscriber, err := nats.NewNATSSubscriber(
		server.config.NatsHost, server.config.NatsPort,
		server.config.NatsUser, server.config.NatsPass, subject, queueGroup)
	if err != nil {
		log.Fatal(err)
	}
	return subscriber
}

func (server *Server) initUpdateUserHandler(service *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handler.NewUpdateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initCreateUserHandler(service *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) {
	_, err := handler.NewCreateUserCommandHandler(service, publisher, subscriber)
	if err != nil {
		log.Fatal(err)
	}
}

func (server *Server) initPostgresClient() *gorm.DB {
	client, err := repository.GetClient(
		server.config.UserDBHost, server.config.UserDBUser,
		server.config.UserDBPass, server.config.UserDBName,
		server.config.UserDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initNeo4jDriver() *neo4j.Driver {
	driver, err := repository.GetNeo4jClient("bolt://accommodation_recommendation_db:7687", "neo4j", "password")

	if err != nil {
		return nil
	}

	return &driver
}

func (server *Server) initUserRepository(client *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(client)
}

func (server *Server) initNotificationRepository(client *gorm.DB) *repository.NotificationRepository {
	return repository.NewNotificationRepository(client)
}

func (server *Server) initUserNotificationRepository(client *gorm.DB) *repository.UserNotificationRepository {
	return repository.NewUserNotificationRepository(client)
}

func (server *Server) initUserNeo4jRepository(driver *neo4j.Driver) *repository.Neo4jUserRepository {
	return repository.NewNeo4jUserRepository(*driver)
}

func (server *Server) initRatingRepository(client *gorm.DB) *repository.RatingRepository {
	return repository.NewRatingRepository(client)
}

func (server *Server) initUserService(repo *repository.UserRepository, notRepo *repository.NotificationRepository, userNotRepo *repository.UserNotificationRepository, neo4jRepo *repository.Neo4jUserRepository, changePassOrchestrator *service.ChangePasswordOrchestrator, updateUserOrchestrator *service.UpdateUserOrchestrator) *service.UserService {
	return service.NewUserService(repo, notRepo, userNotRepo, neo4jRepo, changePassOrchestrator, updateUserOrchestrator)
}

func (server *Server) initRatingService(repo *repository.RatingRepository) *service.RatingService {
	return service.NewRatingService(repo)
}

func (server *Server) initUserHandler(service *service.UserService, ratingService *service.RatingService) *handler.UserHandler {
	return handler.NewUserHandler(service, ratingService)
}

func (server *Server) startGrpcServer(userHandler *handler.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	user.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)
	println("GRPC SERVER USPJESNO NAPRAVLJEN")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
		println("GRPC SERVER NIJE USPJESNO NAPRAVLJEN")
	}
}
