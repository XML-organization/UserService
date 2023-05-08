package service

import (
	"user_service/model"

	events "github.com/XML-organization/common/saga/change_password"
)

func MapChangePasswordToSagaChangePasswordDTO(changePassword *model.UserPassword) *events.ChangePasswordDTO {
	return &events.ChangePasswordDTO{
		Email:       changePassword.Email,
		NewPassword: changePassword.NewPassword,
		OldPassword: changePassword.OldPassword,
	}
}

func mapRollbackChangePasswordDTO(p *model.UserPassword) *model.UserPassword {
	newPassword := p.OldPassword
	oldPassword := p.NewPassword

	return &model.UserPassword{
		Email:       p.Email,
		NewPassword: newPassword,
		OldPassword: oldPassword,
	}
}
