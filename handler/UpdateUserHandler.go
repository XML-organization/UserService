package handler

import (
	"log"
	"user_service/model"
	"user_service/service"

	saga "github.com/XML-organization/common/saga/messaging"
	events "github.com/XML-organization/common/saga/update_user"
)

type UpdateUserCommandHandler struct {
	userService       *service.UserService
	replyPublisher    saga.Publisher
	commandSubscriber saga.Subscriber
}

func NewUpdateUserCommandHandler(userService *service.UserService, publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateUserCommandHandler, error) {
	o := &UpdateUserCommandHandler{
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

func (handler *UpdateUserCommandHandler) handle(command *events.UpdateUserCommand) {

	reply := events.UpdateUserReply{UpdateUserDTO: command.UpdateUserDTO}

	switch command.Type {
	case events.PrintSuccessful:
		log.Println("Saga (User servise side): User password changed successfuly!")
		reply.Type = events.SuccessfulyFinished
	case events.RollbackEmail:
		err := handler.userService.ChangeEmail(&model.UpdateEmailDTO{NewEmail: command.UpdateUserDTO.NewEmail, OldEmail: command.UpdateUserDTO.OldEmail})
		if err != nil {
			log.Println(err)
			return
		}
		println("Saga (User servise side): User old user password returned successfuly!")
		reply.Type = events.OldEmailReturned
	default:
		reply.Type = events.UnknownReply
	}

	if reply.Type != events.UnknownReply {
		_ = handler.replyPublisher.Publish(reply)
	}
}
