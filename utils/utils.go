package utils

import (
	"errors"
	"os"
)

// GetGOPATH returns environment variable GOPATH
func GetGOPATH() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("Please set the $GOPATH environment variable")
	}

	return gopath, nil
}
