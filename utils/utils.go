package utils

import (
	"log"
	"os"
	"regexp"
	"strings"
)

// ValidatePath helper to validate if a path exists
func ValidatePath(p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return err
	}

	return nil
}

// Sanitize is used to remove special chars and trim a string
func Sanitize(s string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	s = strings.Trim(s, " ")
	s = reg.ReplaceAllString(s, "")
	s = strings.Replace(s, " ", "-", -1)
	s = strings.ToLower(s)

	return s
}
