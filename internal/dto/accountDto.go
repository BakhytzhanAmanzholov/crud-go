package dto

type AccountDto struct {
	Username string `json:"username" validate:"min=3"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6"`
}

type LoginDto struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6"`
}
