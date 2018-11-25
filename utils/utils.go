package utils

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// GetIntFromFile parses the first line of a file into an
// integer.
func GetIntFromFile(filepath string) (value int, err error) {
	output, err := ioutil.ReadFile(filepath)
	if err == nil {
		value, err = strconv.Atoi(strings.TrimSpace(string(output)))
	}
	return
}
