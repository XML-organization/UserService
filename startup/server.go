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

	"google.golang.org/grpc"
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
	userRepo := server.initUserRepository(postgresClient)
	userService := server.initUserService(userRepo)

	commandSubscriber := server.initSubscriber(server.config.CreateUserCommandSubject, QueueGroup)
	replyPublisher := server.initPublisher(server.config.CreateUserReplySubject)
	server.initCreateUserHandler(userService, replyPublisher, commandSubscriber)

	userHandler := server.initUserHandler(userService)

	server.startGrpcServer(userHandler)
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

func (server *Server) initUserRepository(client *gorm.DB) *repository.UserRepository {
	return repository.NewUserRepository(client)
}

func (server *Server) initUserService(repo *repository.UserRepository) *service.UserService {
	return service.NewUserService(repo)
}

func (server *Server) initUserHandler(service *service.UserService) *handler.UserHandler {
	return handler.NewUserHandler(service)
}

func (server *Server) startGrpcServer(userHandler *handler.UserHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	user.RegisterUserServiceServer(grpcServer, userHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
