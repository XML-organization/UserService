package handler

import (
	"user_service/service"

	events "github.com/XML-organization/common/saga/create_user"
	saga "github.com/XML-organization/common/saga/messaging"
)

type CreateUserCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewCreateUserCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*CreateUserCommandHandler, error) {
	o := &CreateUserCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (handler *CreateUserCommandHandler) handle(command *events.CreateUserCommand) {
	reply := events.CreateUserReply{User: command.User}

	println("Usao sam u handle metodu na User service strani")
	//println("Tip komande " + command.Type)

	switch command.Type {
	case events.SaveUser:
		println("ID:")
		println(command.User.ID)
		user := mapSagaUserToUser(&command.User)
		_, err := handler.userService.CreateUser(*user)
		println("Saga (User servise side): User created successfuly!")
		if err != nil {
			reply.Type = events.UserNotSaved
			println("Saga (User servise side): User dont created successfuly!")
			break
		}
		reply.Type = events.UserSaved
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
