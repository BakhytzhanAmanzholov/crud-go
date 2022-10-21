package repositories

import (
	"context"
	"crud-golang/internal/models"
)

type Repository interface {
	Create(user *models.Account, ctx context.Context) (*models.Account, error)
	FindAll(ctx context.Context) (u *[]models.Account, err error)
	FindOne(id string, ctx context.Context) (*models.Account, error)
	Update(id string, user models.Account, ctx context.Context) (*models.Account, error)
	Delete(id string, ctx context.Context) error
	FindByEmail(email string, ctx context.Context) (*models.Account, error)
}
