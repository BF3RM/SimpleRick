package env

import (
	"github.com/joho/godotenv"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
}

func GetString(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
