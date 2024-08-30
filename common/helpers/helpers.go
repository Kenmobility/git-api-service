package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"math/rand"

	"gopkg.in/go-playground/validator.v9"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

// IsLocal returns true or false depending on APP_ENV environmental variable's value
func IsLocal() bool {
	return os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local"
}

// Getenv gets the env variable value or set a default if empty
func Getenv(variable string, defaultValue ...string) string {
	env := os.Getenv(variable)
	if env == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return ""
	}
	return env
}

// IsRepositoryNameValid validates the repository name
func IsRepositoryNameValid(repoName string) bool {
	return strings.Contains(repoName, "/")
}

// ValidateInput validates structs fields with tags
func ValidateInput(input interface{}) []string {
	var errors []string
	v := validator.New()

	err := v.Struct(input)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.ActualTag() {
			case "required":
				errors = append(errors, fmt.Sprintf("%s field is required", e.Field()))
			default:
				errors = append(errors, "an error occurred")
			}
		}
	}

	return errors
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomRepositoryName generates a random repository name
func RandomRepositoryName() string {
	return fmt.Sprintf("%s/%s", RandomString(6), RandomString(6))
}

// RandomRepositoryUrl generates a random repository URL
func RandomRepositoryUrl() string {
	return fmt.Sprintf("https://github.com/%s/%s", RandomString(6), RandomString(6))
}

// RandomFetchStartDate generates a random start date for fetching commits
func RandomFetchStartDate() time.Time {
	return time.Now().AddDate(0, -8, 0)
}

// RandomFetchEndDate generates a random end date for fetching commits
func RandomFetchEndDate() time.Time {
	return time.Now()
}

// RandomWords generates some random words of specified length
func RandomWords(words int) string {
	var sb strings.Builder

	for i := 0; i < words; i++ {
		m := RandomString(5)
		sb.WriteString(m)
		sb.WriteString(" ")
	}

	return sb.String()
}
