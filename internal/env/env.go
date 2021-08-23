package env

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
	"reflect"
	"strings"
)

type envStructField struct {
	key          string
	defaultValue string
	value        reflect.Value
}

var ErrInvalidStruct = errors.New("given parameter is not a pointer to a struct")

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

func GetBytes(key string, defaultVal []byte) []byte {
	if value, exists := os.LookupEnv(key); exists {
		return []byte(value)
	}

	return defaultVal
}

func FillStruct(conf interface{}) {
	v := reflect.ValueOf(conf)
	if v.Kind() != reflect.Ptr {
		panic(ErrInvalidStruct)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(ErrInvalidStruct)
	}

	for i := 0; i < v.NumField(); i++ {
		tag, ok := v.Type().Field(i).Tag.Lookup("env")
		if !ok {
			continue
		}

		args := strings.Split(tag, ",")

		f := envStructField{key: args[0], value: v.Field(i)}
		if len(args) == 2 {
			f.defaultValue = args[1]
		}

		switch f.value.Kind() {
		case reflect.String:
			f.value.SetString(GetString(f.key, f.defaultValue))
		}
	}
}
