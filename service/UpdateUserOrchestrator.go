package service

import (
	"log"

	saga "github.com/XML-organization/common/saga/messaging"
	events "github.com/XML-organization/common/saga/update_user"
)

type UpdateUserOrchestrator struct {
	commandPublisher saga.Publisher
	replySubscriber  saga.Subscriber
}

func NewUpdateUserOrchestrator(publisher saga.Publisher, subscriber saga.Subscriber) (*UpdateUserOrchestrator, error) {
	o := &UpdateUserOrchestrator{
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

func (o *UpdateUserOrchestrator) Start(emails *events.UpdateUserDTO) error {
	event := &events.UpdateUserCommand{
		Type:          events.UpdateUser,
		UpdateUserDTO: *emails,
	}

	return o.commandPublisher.Publish(event)
}

func (o *UpdateUserOrchestrator) handle(reply *events.UpdateUserReply) {
	command := events.UpdateUserCommand{UpdateUserDTO: reply.UpdateUserDTO}
	command.Type = o.nextCommandType(reply.Type)
	if command.Type != events.UnknownCommand {
		_ = o.commandPublisher.Publish(command)
	}
}

func (o *UpdateUserOrchestrator) nextCommandType(reply events.UpdateUserReplyType) events.UpdateUserCommandType {

	switch reply {
	case events.UserUpdated:
		log.Println("Event: UserUpdated")
		return events.PrintSuccessful
	case events.UserNotUpdated:
		log.Println("Event: UserNotUpdated")
		return events.RollbackEmail
	default:
		return events.UnknownCommand
	}
}
