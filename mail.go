package main

import (
	"github.com/les-cours/user-service/utils"
	"log"
)

func main() {

	in :=
		struct {
			Firstname string
			Lastname  string
			Email     string
		}{
			Firstname: "chouaib",
			Lastname:  "amine",
			Email:     "chouaibe708@gmail.com",
		}
	var emailData = struct {
		CompanyName string
		Receiver    string
		Email       string
	}{
		CompanyName: in.Firstname + " " + in.Lastname,
		Receiver:    in.Email,
	}

	var emailSubject = "Account Registration confirmation"
	var emailTemplate = "registration-confirmation"

	err := utils.GenerateEmail(in.Email, emailSubject, emailTemplate, emailData)

	if err != nil {
		log.Println(err)

	}
}
