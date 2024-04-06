package utils

import (
	"github.com/les-cours/user-service/api/users"
	"regexp"
)

func ValidateUsername(username string) bool {
	if len(username) < 1 || len(username) > 64 {
		return false
	}
	return true
}

func ValidateFirstname(name string) bool {
	if len(name) < 1 || len(name) > 64 {
		return false
	}
	return true
}

func ValidateLastname(name string) bool {
	if len(name) < 1 || len(name) > 64 {
		return false
	}
	return true
}

func ValidateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}

func ValidatePassword(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func ValidateStudentDetails(user *users.StudentUpdateRequest) (bool, string) {

	var fields string
	var valid = true
	if !ValidateUsername(user.Username) {
		fields = "username, "
		valid = false
	}
	if !ValidateFirstname(user.Firstname) {
		fields += "firstname, "
		valid = false
	}
	if !ValidateLastname(user.Lastname) {
		fields += "lastname, "
		valid = false
	}
	if !ValidateEmail(user.Email) {
		fields += "email "
		valid = false
	}

	return valid, fields
}
