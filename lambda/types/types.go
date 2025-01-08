package types

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerdUser RegisterUser) (*User, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte (registerdUser.Password), 10)
	if err != nil {
		return nil, err
	}
	return &User{
		Username: registerdUser.Username,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidateFunction (hashedPassword, plainTextPassword string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword))
	return err == nil
}

func CreateToken(user User) string{
	now := time.Now()
	validUntil := now.Add(time.Hour * 1).Unix()
	claims := jwt.MapClaims{
		"user": user.Username,
		"expires": validUntil,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims, nil)
	secret := "secret"
	tokenString,err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return tokenString
}