package user

import (
	"crud-golang/internal/dto"
	"crud-golang/internal/models"
)

type Service interface {
	Create(user dto.AccountDto) (models.Account, error)
	FindAll() (u []models.Account, err error)
	FindOne(id string) (models.Account, error)
	Update(user dto.AccountDto, id string) (models.Account, error)
	Delete(id string) error
	FindByEmail(email string) (models.Account, error)
	Login(email string, password string) (models.Account, error)
}
