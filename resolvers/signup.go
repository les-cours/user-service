package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
)

func (s *Server) DoesUserNameExist(ctx context.Context, in *users.DoesUserNameExistRequest) (*users.DoesUserNameExistResponse, error) {
	exists := false

	err := s.DB.QueryRow(
		`SELECT exists (SELECT 1 FROM accounts WHERE username = $1 LIMIT 1);
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

func (s *Server) StudentSignup(ctx context.Context, in *users.StudentSignupRequest) (*users.StudentSignupResponse, error) {

	var accountID uuid.UUID
	accountID, _ = uuid.NewRandom()

	var err error

	/*
		VALIDATION
	*/
	checkEmail, err := s.DoesEmailExist(ctx, &users.DoesEmailExistRequest{
		Type:  "student",
		Email: in.Email,
	})

	switch {
	case err != nil:
		return nil, ErrInternal
	case !utils.ValidateFirstname(in.Firstname):
		return nil, ErrInvalidInput("first", "must be > 1 and < 64")
	case !utils.ValidateLastname(in.Lastname):
		return nil, ErrInvalidInput("lastName", "must be > 1 and < 64")
	case checkEmail.Exists:
		return nil, ErrInvalidInput("email", "wrong")
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

	var userName = in.GetFirstname() + "_" + in.GetLastname()
	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO accounts 
		(account_id,email,password,username, status, user_type,plan_id) VALUES($1, $2,crypt($3,gen_salt('bf')),$4,$5,$6,$7);`,
		accountID, in.GetEmail(), in.GetPassword(), userName, "active", "student", "PLAN_free")

	if err != nil {
		log.Printf("Err when set accounts err %v", err)
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
		log.Println("Err When INSERT permissions err:", err)
		tx.Rollback()
		return nil, ErrInternal
	}

	// generating avatar
	var defaultAvatar string
	defaultAvatar, err = utils.GenerateAvatar(in.Firstname, in.Lastname)

	_, err = tx.ExecContext(context.Background(),
		`
		INSERT INTO 
		students
		(student_id, firstname, lastname,default_avatar,grade_id,date_of_birth,gender,city_id)
		VALUES($1, $2, $3, $4, $5,$6,$7,$8);
		`,
		accountID,
		in.Firstname,
		in.Lastname,
		defaultAvatar,
		in.GradID,
		in.Dob,
		in.Gender,
		in.CityID,
	)
	if err != nil {
		log.Println("Err When INSERT students err:", err)
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

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err := s.AuthService.Signup(ctx, &auth.SignUpRequest{
		AccountID: accountID.String(),
		Email:     in.Email,
	})

	if err != nil {
		return nil, errors.New("register success, but err when try login, login again")
	}

	return &users.StudentSignupResponse{
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
