package user

import (
	"crud-golang/internal/dto"
	"crud-golang/internal/models"
	"crud-golang/internal/repositories"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repository repositories.Repository
}

func (s service) Create(user dto.AccountDto) (models.Account, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Account{}, err
	}

	account := models.Account{
		Password: string(hash),
		Email:    user.Email,
		Username: user.Username,
	}
	account, err = s.repository.Create(&account)
	if err != nil {
		return account, err
	}
	return account, err
}

func (s service) FindAll() (u []models.Account, err error) {
	accounts, err := s.repository.FindAll()
	return accounts, err
}

func (s service) FindOne(id string) (models.Account, error) {
	account, err := s.repository.FindOne(id)
	return account, err
}

func (s service) Update(user dto.AccountDto, id string) (models.Account, error) {
	account := models.Account{
		Password: user.Password,
		Email:    user.Email,
		Username: user.Username,
	}
	account, err := s.repository.Update(id, account)
	if err != nil {
		return account, err
	}
	return account, err
}

func (s service) FindByEmail(email string) (models.Account, error) {
	return s.repository.FindByEmail(email)
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) Login(email string, password string) (models.Account, error) {
	account, err := s.FindByEmail(email)
	if err != nil {
		return models.Account{}, errors.New("invalid login credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return models.Account{}, errors.New("incorrect password")
	}

	return account, nil
}

func NewService(repository repositories.Repository) Service {
	return &service{
		repository: repository,
	}
}
