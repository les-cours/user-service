package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
	"log"
)

func (s *Server) StudentPasswordReset(ctx context.Context, in *users.StudentPasswordResetRequest) (*users.OperationStatus, error) {

	var studentInfo = struct {
		Firstname string
		Lastname  string
		Email     string
	}{
		Firstname: "",
		Lastname:  "",
		Email:     "",
	}
	var err = s.DB.QueryRow(
		`	
        SELECT 
            email, firstname, lastname
		FROM 
		    students
		INNER JOIN accounts ON students.account_id = accounts.account_id
		WHERE 
		    students.id = $1
		  AND accounts.password = crypt($2, accounts.password)
		`,
		in.AccountId, in.OldPassword).Scan(&studentInfo.Email, &studentInfo.Firstname, &studentInfo.Lastname)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidInput("oldPassword", " wrong password")
		}
		return nil, ErrInternal
	}

	var _ = s.DB.QueryRow(
		`UPDATE accounts SET password = crypt($1,gen_salt('bf', 8)) WHERE id = $2`,
		in.NewPassword, in.AccountId)
	if err != nil {
		return nil, ErrInternal
	}

	var emailData = struct {
		Username string
		Receiver string
	}{
		Username: studentInfo.Firstname + " " + studentInfo.Lastname,
		Receiver: studentInfo.Email,
	}

	var emailSubject = "Password reset success"
	var emailTemplate = "password-reset-success"

	message, err := utils.GenerateEmail(studentInfo.Email, emailSubject, emailTemplate, emailData)

	go func() {
		_, err = s.SendGridClient.Send(message)
		if err != nil {
			log.Println(err)
		}
	}()

	return &users.OperationStatus{
		Completed: true,
	}, nil
}
