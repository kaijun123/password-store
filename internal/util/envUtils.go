package util

import (
	"errors"
	"os"
)

func GetEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", errors.New("env variable not set")
	}
	return value, nil
}

func SetEnv(name string, value string) error {
	if err := os.Setenv(name, value); err != nil {
		return errors.New("env variable not set")
	}
	return nil
}
