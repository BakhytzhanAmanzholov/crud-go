package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type Registration struct {
	Username string `json:"username" validate:"min=3"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6"`
}

type Login struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=6"`
}

type Account struct {
	Id       primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Username string             `json:"username" validate:"min=3"`
	Email    string             `json:"email" validate:"email"`
}
