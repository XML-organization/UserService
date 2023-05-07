package service

import (
	"user_service/model"
	"user_service/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: repo,
	}
}

func (service *UserService) UpdateUser(user model.ChangeUserDTO) (model.RequestMessage, error) {
	err := service.UserRepo.UpdateUser(user)
	if err != nil {
		return model.RequestMessage{Message: "An error occurred, please try again!"}, err
	}
	return model.RequestMessage{Message: "Success!"}, nil
}

func (service *UserService) ChangePassword(userPasswords model.UserPassword) (model.RequestMessage, error) {
	user, err := service.UserRepo.FindById(userPasswords.ID)

	if err != nil {
		message := model.RequestMessage{
			Message: "An error occurred, please try again!",
		}
		return message, err
	} else if err := bcrypt.CompareHashAndPassword(user.Password, []byte(userPasswords.OldPassword)); err != nil {
		message := model.RequestMessage{
			Message: "The old password is not correct!",
		}
		return message, err
	}

	newPassword, _ := bcrypt.GenerateFromPassword([]byte(userPasswords.NewPassword), 14)
	userPasswords.NewPassword = string(newPassword)

	return service.UserRepo.ChangePassword(userPasswords)
}

func (service *UserService) CreateUser(user model.User) (model.RequestMessage, error) {

	println("usao sam u create user metodu na user service strani")

	response := model.RequestMessage{
		Message: service.UserRepo.CreateUser(user).Message,
	}

	return response, nil
}
