package user

import (
	"crud-golang/internal/dto"
)

type Service interface {
	Create(user dto.Registration) (*dto.Account, error)
	FindAll() (u *[]dto.Account, err error)
	FindOne(id string) (*dto.Account, error)
	Update(user dto.Registration, id string) (*dto.Account, error)
	Delete(id string) error
	FindByEmail(email string) (*dto.Account, error)
	Login(email string, password string) (*dto.Account, error)
}
