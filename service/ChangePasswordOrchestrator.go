package service

import (
	"log"
	"user_service/model"

	events "github.com/XML-organization/common/saga/change_password"
	saga "github.com/XML-organization/common/saga/messaging"
)

type ChangePasswordOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewChangePasswordOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*ChangePasswordOrchestrator, error) {
	o := &ChangePasswordOrchestrator{
		commandPublisher: publisher,
		replySubscriber:  subscriber,
	}
	err := o.replySubscriber.Subscribe(o.handle)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return o, nil
}

func (o *ChangePasswordOrchestrator) Start(changePassword *model.UserPassword) error {
	event := &events.ChangePasswordCommand{
		Type:             events.ChangePassword,
		ChagePasswordDTO: *MapChangePasswordToSagaChangePasswordDTO(changePassword),
	}

	return o.commandPublisher.Publish(event)
}

func (o *ChangePasswordOrchestrator) handle(reply *events.ChangePasswordReply) {
	command := events.ChangePasswordCommand{ChagePasswordDTO: reply.ChangePasswordDTO}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *ChangePasswordOrchestrator) nextCommandType(reply events.ChangePasswordReplyType) events.ChangePasswordCommandType {

	log.Println("ChangePasswordOrchestrator: Usao sam u nextCommandType metodu")

	switch reply {
	case events.PasswordChanged:
		log.Println("Event: PasswordChanged")
		return events.PrintSuccessful
	case events.PasswordNotChanged:
		log.Println("Event: PasswordNotChanged")
		return events.RollbackPassword
	default:
		return events.UnknownCommand
	}
}
