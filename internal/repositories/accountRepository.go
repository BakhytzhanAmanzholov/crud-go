package repositories

import (
	"crud-golang/internal/models"
)

type Repository interface {
	Create(user *models.Account) (models.Account, error)
	FindAll() (u []models.Account, err error)
	FindOne(id string) (models.Account, error)
	Update(id string, user models.Account) (models.Account, error)
	Delete(id string) error
	FindByEmail(email string) (models.Account, error)
}
