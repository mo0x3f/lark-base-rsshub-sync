package secret

import (
	"errors"
	"os"
)

var verifyToken = ""

func Init(env string) error {
	switch env {
	case "local":
	case "replit":
		verifyToken = os.Getenv("BASE_VERIFY_TOKEN")
		if verifyToken == "" {
			return errors.New("verify token not found")
		}
	}
	return nil
}

func GetVerifyToken() string {
	return verifyToken
}
