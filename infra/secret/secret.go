package secret

import (
	"errors"
	"os"
)

var verifyToken = ""

func Init(env string) error {
	switch env {
	case "local", "replit", "qcloud":
		verifyToken = os.Getenv("BASE_VERIFY_TOKEN")
		if verifyToken == "" {
			return errors.New("verify token not found")
		}
		return nil
	default:
		return errors.New("verify token not found")
	}
}

func GetVerifyToken() string {
	return verifyToken
}
