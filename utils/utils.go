package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetIntFromFile returns the first number in a 
// file
func GetIntFromFile(filepath string) (value int, err error) {
	output, err := ioutil.ReadFile(filepath)
	if err == nil {
		value, err = strconv.Atoi(strings.TrimSpace(string(output)))
	}
	return
}
