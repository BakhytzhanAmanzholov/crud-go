package mappers

import (
	"crud-golang/internal/dto"
	"crud-golang/internal/models"
)

func MapperOneToDto(account *models.Account) (*dto.Account, error) {
	data := dto.Account{
		Id:       account.Id,
		Username: account.Username,
		Email:    account.Email,
	}

	return &data, nil
}

func MapperManyToDto(accounts []models.Account) (*[]dto.Account, error) {
	var dtos []dto.Account

	for _, account := range accounts {
		data, err := MapperOneToDto(&account)
		if err != nil {
			return nil, err
		}
		dtos = append(dtos, *data)
	}

	return &dtos, nil
}
