package service

import (
	"user_service/model"
	"user_service/repository"

	events "github.com/XML-organization/common/saga/update_user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	UserRepo               *repository.UserRepository
	UserNeo4jRepo          *repository.Neo4jUserRepository
	orchestrator           *ChangePasswordOrchestrator
	updateUserOrchestrator *UpdateUserOrchestrator
}

func NewUserService(repo *repository.UserRepository, neo4jRepo *repository.Neo4jUserRepository, orchestrator *ChangePasswordOrchestrator, updateUserOrchestrator *UpdateUserOrchestrator) *UserService {
	return &UserService{
		UserRepo:               repo,
		UserNeo4jRepo:          neo4jRepo,
		orchestrator:           orchestrator,
		updateUserOrchestrator: updateUserOrchestrator,
	}
}

func (service *UserService) UpdateUser(user model.ChangeUserDTO) (model.RequestMessage, error) {
	loggedUser, err1 := service.UserRepo.FindById(user.ID)
	if err1 != nil {
		return model.RequestMessage{Message: "An error occurred, please try again!"}, err1
	}

	err := service.UserRepo.UpdateUser(user)
	if err != nil {
		return model.RequestMessage{Message: "An error occurred, please try again!"}, err
	}

	if loggedUser.Email != user.Email {
		err2 := service.updateUserOrchestrator.Start(&events.UpdateUserDTO{OldEmail: loggedUser.Email, NewEmail: user.Email})
		if err2 != nil {
			service.UserRepo.ChangeEmail(model.UpdateEmailDTO{OldEmail: user.Email, NewEmail: loggedUser.Email})
			return model.RequestMessage{Message: "An error occurred, please try again!"}, err2
		}
	}

	return model.RequestMessage{Message: "Success!"}, nil
}

func (service *UserService) ChangePassword(userPasswords model.UserPassword) (model.RequestMessage, error) {
	user, err := service.UserRepo.FindByEmail(userPasswords.Email)

	passwords := model.UserPassword{
		Email:       userPasswords.Email,
		OldPassword: userPasswords.OldPassword,
		NewPassword: userPasswords.NewPassword,
	}

	println("MJENJANJE SIFREE")
	println("Email: " + passwords.Email)
	println("Novi password: " + userPasswords.NewPassword)
	println("Stari password: " + userPasswords.OldPassword)

	if err != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message, err
	} else if err := bcrypt.CompareHashAndPassword(user.Password, []byte(userPasswords.OldPassword)); err != nil {
		message := model.RequestMessage{
			Message: "The old password is not correct!",
		}
		return message, status.Error(codes.OK, "The old password is not correct!")
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(userPasswords.NewPassword), 14)
	userPasswords.NewPassword = string(newPassword)

	message, err := service.UserRepo.ChangePassword(userPasswords)
	if err != nil {
		message := model.RequestMessage{
			Message: message.Message,
		}
		return message, err
	}

	err1 := service.orchestrator.Start(&passwords)
	if err1 != nil {
		service.UserRepo.ChangePassword(*mapRollbackChangePasswordDTO(&userPasswords))
		message := model.RequestMessage{
			Message: "An error occured, please try again!",
		}
		return message, err
	}

	message1 := model.RequestMessage{
		Message: "Successful",
	}
	return message1, err
}

func (service *UserService) CreateUser(user model.User) (model.RequestMessage, error) {

	println("USAO U METODU CREATE USER U SERVICE")

	message, err := service.UserRepo.CreateUser(user)
	err1 := service.UserNeo4jRepo.SaveUser(user.ID.String())

	if err1 != nil {
		println("nisam uspio upisati usera u neo4j bazu")
	}

	if err != nil {
		response := model.RequestMessage{
			Message: "An error occured, please try again!",
		}

		return response, err
	}

	response := model.RequestMessage{
		Message: message.Message,
	}

	return response, nil
}

func (service *UserService) ChangeEmail(emails *model.UpdateEmailDTO) error {
	_, err := service.UserRepo.FindByEmail(emails.OldEmail)
	if err != nil {
		return err
	}

	err1 := service.UserRepo.ChangeEmail(*emails)
	if err1 != nil {
		return err1
	}

	return nil
}

func (service *UserService) DeleteUser(email string) (model.RequestMessage, error) {
	user, err := service.UserRepo.FindByEmail(email)
	println("front email---" + email)
	println("user email---" + user.Email)
	message, err := service.UserRepo.DeleteUserById(user.ID)

	if err != nil {
		response := model.RequestMessage{
			Message: "An error occured, please try again!",
		}

		return response, err
	}

	response := model.RequestMessage{
		Message: message.Message,
	}

	return response, nil
}
