package utils

import "os"

func GetJWTKey() []byte {
	return []byte(os.Getenv("SECRET_KEY"))
}