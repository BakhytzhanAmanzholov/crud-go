package user

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

func CreateJWTToken(id primitive.ObjectID) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 30).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = id
	claims["exp"] = exp
	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", 0, err
	}

	return t, exp, nil
}

func VerifyJWT(tokenHeader string) (bool, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("StatusUnauthorized")
		}
		return []byte("secret"), nil
	}

	_, err := jwt.Parse(tokenHeader, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errors.New("StatusUnauthorized")) {
			return false, errors.New("StatusUnauthorized")
		}
		return false, errors.New("StatusUnauthorized")
	}

	return true, nil
}
