package service

import (
	"github.com/dgrijalva/jwt-go"
	"os"
)

func GenerateJWT(studentId string) string {
	result, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": studentId,
	}).SignedString([]byte(os.Getenv("JWT_SECRET")))
	return result
}
