package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
	"log"
)

func (s *Server) UserPasswordReset(ctx context.Context, in *users.UserPasswordResetRequest) (*users.OperationStatus, error) {

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
            email, s.firstname, s.lastname
		FROM 
		    accounts
		Inner Join public.students s on accounts.account_id = s.id
		WHERE 
		    accounts.account_id = $1
		  AND accounts.password = crypt($2, accounts.password)
		`,
		in.UserID, in.OldPassword).Scan(&studentInfo.Email, &studentInfo.Firstname, &studentInfo.Lastname)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidInput("oldPassword", " wrong password")
		}
		return nil, ErrInternal
	}

	var _ = s.DB.QueryRow(
		`UPDATE accounts SET password = crypt($1,gen_salt('bf', 8)) WHERE account_id = $2`,
		in.NewPassword, in.UserID)
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
