package handler

import (
	"log"
	"user_service/model"
	"user_service/service"

	events "github.com/XML-organization/common/saga/change_password"
	saga "github.com/XML-organization/common/saga/messaging"
)

type ChangePasswordCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewChangePasswordCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*ChangePasswordCommandHandler, error) {
	o := &ChangePasswordCommandHandler{
		userService:       userService,
		replyPublisher:    publisher,
		commandSubscriber: subscriber,
	}
	err := o.commandSubscriber.Subscribe(o.handle)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return o, nil
}

func (handler *ChangePasswordCommandHandler) handle(command *events.ChangePasswordCommand) {
	reply := events.ChangePasswordReply{ChangePasswordDTO: command.ChagePasswordDTO}

	log.Println("ChangePassword: Usao sam u handle metodu na User service strani")
	log.Println("Ovo je tip comande koju sam dobio:", command.Type)

	switch command.Type {
	case events.PrintSuccessful:
		log.Println("Saga (User servise side): User password changed successfuly!")
		reply.Type = events.SuccessfulyFinished
	case events.RollbackPassword:
		_, err := handler.userService.ChangePassword(model.UserPassword{Email: command.ChagePasswordDTO.Email,
			OldPassword: command.ChagePasswordDTO.OldPassword,
			NewPassword: command.ChagePasswordDTO.NewPassword})
		if err != nil {
			log.Println("nisam uspjesno obrisao usera")
			return
		}
		log.Println("Saga (User servise side): User old user password returned successfuly!")
		reply.Type = events.OldPasswordReturned
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
