package user

import (
	"crud-golang/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

const jwtSecret = "secret"

func CreateJWTToken(user models.Account) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.Id
	claims["exp"] = exp
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}
