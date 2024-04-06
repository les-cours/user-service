package resolvers

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
	"log"
)

func (s *Server) DoesUserNameExist(ctx context.Context, in *users.DoesUserNameExistRequest) (*users.DoesUserNameExistResponse, error) {
	exists := false

	err := s.DB.QueryRow(
		`SELECT exists (SELECT 1 FROM students WHERE username = $1 LIMIT 1);
	`, in.Username).Scan(&exists)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return &users.DoesUserNameExistResponse{
		Exists: exists,
	}, nil
}

func (s *Server) DoesEmailExist(ctx context.Context, in *users.DoesEmailExistRequest) (*users.DoesEmailExistResponse, error) {
	exists := false

	err := s.DB.QueryRow(
		`SELECT exists (SELECT 1 FROM accounts WHERE email = $1 LIMIT 1);
	`, in.Email).Scan(&exists)

	if err != nil {
		log.Println(err)
		return nil, ErrInternal
	}

	return &users.DoesEmailExistResponse{
		Exists: exists,
	}, nil
}

func (s *Server) DoesSignupLinkExist(ctx context.Context, email string) (bool, error) {
	exists := false

	err := s.DB.QueryRow(
		`SELECT exists (SELECT 1 FROM signup_links WHERE email = $1 LIMIT 1);
	`, email).Scan(&exists)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return exists, nil
}

func (s *Server) Signup(ctx context.Context, in *users.SignupRequest) (*users.SignupResponse, error) {

	var err error

	/*
		VALIDATION
	*/
	switch {
	case err != nil:
		return nil, ErrInternal

	case !utils.ValidateFirstname(in.Firstname):
		return nil, ErrInvalidInput("first", "must be > 1 and < 64")

	case !utils.ValidateLastname(in.Lastname):
		return nil, ErrInvalidInput("lastName", "must be > 1 and < 64")
	}

	/*
		SQL
	*/

	var tx *sql.Tx
	tx, err = s.DB.BeginTx(context.Background(), nil)
	defer tx.Rollback()
	if err != nil {
		log.Printf("Err when BEGINTX err %v", err)
		return nil, ErrInternal
	}

	var accountID string
	err = tx.QueryRowContext(context.Background(),
		`INSERT INTO accounts 
		(email, password, status, user_type) VALUES($1, crypt($2, gen_salt('bf')), $3, $4)
		 RETURNING account_id;`,
		in.GetEmail(), in.GetPassword(), "active", "student").Scan(&accountID)

	if err != nil {
		log.Println("107 : " + err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	_, err = tx.ExecContext(context.Background(),
		`
		INSERT INTO
		permissions (account_id,live,write_comment,settings)
		VALUES
			($1,FALSE, FALSE, FALSE)
		
		`, &accountID)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, ErrInternal
	}

	if len(in.Password) > 64 {
		tx.Rollback()
		return nil, ErrInvalidInput("password", "password length should not exceed 64 characters")
	}

	// generating avatar
	var defaultAvatar string
	defaultAvatar, err = utils.GenerateAvatar(in.Firstname, in.Lastname)

	var studentID uuid.UUID
	studentID, err = uuid.NewRandom()
	_, err = tx.ExecContext(context.Background(),
		`
		INSERT INTO 
		students
		(id, account_id, username, firstname, lastname,default_avatar)
		VALUES($1, $2, $3, $4, $5, $6) 
		RETURNING id;
		`,
		studentID,
		accountID,
		in.Firstname+in.Lastname,
		in.Firstname,
		in.Lastname,
		defaultAvatar,
	)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, ErrInternal
	}

	var emailData = struct {
		CompanyName string
		Receiver    string
	}{
		CompanyName: "",
		Receiver:    in.Email,
	}

	var emailSubject = "Account Registration confirmation"
	var emailTemplate = "registration-confirmation"

	message, err := utils.GenerateEmail(in.Email, emailSubject, emailTemplate, emailData)
	if err != nil {
		log.Println("Error generating email")
		log.Println(err)
		return nil, ErrInternal
	}

	go func() {
		var _, err = s.SendGridClient.Send(message)
		if err != nil {
			log.Println("Error sending email")
			log.Println(err)
		}
	}()

	err = tx.Commit()

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return nil, ErrInternal
	}
	//_, err = s.StripeService.CreateFreeCustomer(ctx, &stripe.CreateFreeCustomerRequest{
	//	AccountID: accountID,
	//	Email:     in.Email,
	//	Name:      in.Firstname + " " + in.Lastname,
	//},
	//)
	//if err != nil {
	//	log.Println(err)
	//	tx.Rollback()
	//	return nil, ErrInternal
	//}

	// Generate token and send it to the user

	res, err := s.AuthService.Signup(ctx, &auth.SignUpRequest{
		UserName:  in.Firstname + in.Lastname,
		Password:  in.Password,
		AccountID: accountID,
		StudentID: studentID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &users.SignupResponse{
		Succeeded: true,
		AccessToken: &users.AccessToken{
			Token:     res.AccessToken.Token,
			TokenType: res.AccessToken.TokenType,
			ExpiresAt: res.AccessToken.ExpiresAt,
		},
		RefreshToken: &users.RefreshToken{
			Token:     res.RefreshToken.Token,
			ExpiresAt: res.RefreshToken.ExpiresAt,
		},
		SignupToken: &users.SignupToken{
			Token:     res.SignupToken.Token,
			ExpiresAt: res.SignupToken.ExpiresAt,
		},
	}, nil
}
