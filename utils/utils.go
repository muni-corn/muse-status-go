package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetIntFromFile returns a number in a file
func GetIntFromFile(filepath string) (value int, err error) {
	str, err := GetStringFromFile(filepath)
	if err == nil {
		value, err = strconv.Atoi(str)
	}
	return
}

// GetStringFromFile returns a file as a string
func GetStringFromFile(filepath string) (value string, err error) {
	output, err := ioutil.ReadFile(filepath)
	if err == nil {
		value = strings.TrimSpace(string(output))
	}
	return
}
