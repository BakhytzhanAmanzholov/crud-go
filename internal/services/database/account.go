package database

import (
	"context"
	"crud-golang/internal/dto"
	"crud-golang/internal/mappers"
	"crud-golang/internal/models"
	"crud-golang/internal/repositories"
	"crud-golang/internal/services"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository repositories.Repository
}

func (s service) Create(user dto.Registration) (*dto.Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	account := models.Account{
		Id:       primitive.NewObjectID(),
		Password: string(hash),
		Email:    user.Email,
		Username: user.Username,
	}

	result, err := s.repository.Create(&account, context.Background())
	if err != nil {
		return nil, err
	}

	data, err := mappers.MapperOneToDto(result)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s service) FindAll() (u *[]dto.Account, err error) {
	accounts, err := s.repository.FindAll(context.Background())
	data, err := mappers.MapperManyToDto(*accounts)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s service) FindOne(id string) (*dto.Account, error) {
	account, err := s.repository.FindOne(id, context.Background())
	data, err := mappers.MapperOneToDto(account)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s service) Update(user dto.Registration, id string) (*dto.Account, error) {
	account := models.Account{
		Password: user.Password,
		Email:    user.Email,
		Username: user.Username,
	}
	result := &account
	result, err := s.repository.Update(id, account, context.Background())
	if err != nil {
		return nil, err
	}
	data, err := mappers.MapperOneToDto(result)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s service) FindByEmail(email string) (*dto.Account, error) {
	account, err := s.repository.FindByEmail(email, context.Background())
	if err != nil {
		return nil, errors.New("invalid email")
	}
	data, err := mappers.MapperOneToDto(account)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id, context.Background())
}

func (s service) Login(email string, password string) (*dto.Account, error) {
	account, err := s.repository.FindByEmail(email, context.Background())
	if err != nil {
		return nil, errors.New("invalid login credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, errors.New("incorrect password")
	}
	data, err := mappers.MapperOneToDto(account)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewService(repository repositories.Repository) user.Service {
	return &service{
		repository: repository,
	}
}
