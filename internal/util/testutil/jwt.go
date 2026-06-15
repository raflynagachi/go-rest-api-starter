package testutil

import (
	"time"

	appjwt "github.com/raflynagachi/go-rest-api-starter/pkg/jwt"
)

const TestJWTSecret = "test-secret"

func GenerateTestToken(email string) string {
	token, _ := appjwt.GenerateToken(email, TestJWTSecret, time.Hour)
	return "Bearer " + token
}
