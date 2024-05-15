package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
	"time"
)

func (s *Server) DoesUserNameExist(ctx context.Context, in *users.DoesUserNameExistRequest) (*users.DoesUserNameExistResponse, error) {
	exists := false

	err := s.DB.QueryRow(
		`SELECT exists (SELECT 1 FROM accounts WHERE username = $1 LIMIT 1);
	`, in.Username).Scan(&exists)

	if err != nil {
		s.Logger.Error(err.Error())
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
		s.Logger.Error(err.Error())
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
		return nil, ErrInvalidInput("email", "exist already.")
	}

	/*
		SQL
	*/

	var tx *sql.Tx
	tx, err = s.DB.BeginTx(context.Background(), nil)
	defer tx.Rollback()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	var userName = in.GetFirstname() + "_" + in.GetLastname()
	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO accounts 
		(account_id,email,password,username, status, user_type) VALUES($1, $2,crypt($3,gen_salt('bf')),$4,$5,$6);`,
		accountID, in.GetEmail(), in.GetPassword(), userName, "inactive", "student")

	if err != nil {
		s.Logger.Error(err.Error())
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
		(student_id, firstname, lastname,avatar,grade_id,date_of_birth,gender,city_id)
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
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	//generate 5 digital

	code := utils.GenerateConfirmationCode()

	//save them in confirmation_table
	expiresAt := time.Now().Add(time.Minute * 60).Unix()
	_, err = tx.Exec(`INSERT INTO email_confirmation (account_id,code,expires_at) values ($1,$2,$3)`, accountID, code, expiresAt)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, err
	}
	//

	var emailData = struct {
		CompanyName string
		Receiver    string
		Code        int
	}{
		CompanyName: in.Firstname + " " + in.Lastname,
		Receiver:    in.Email,
		Code:        code,
	}

	var emailSubject = "Account Registration confirmation"
	var emailTemplate = "registration-confirmation"

	err = utils.GenerateEmail(in.Email, emailSubject, emailTemplate, emailData)

	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	// Generate token and send it to the user

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res, err := s.AuthService.Signup(ctx, &auth.SignUpRequest{
		AccountID: accountID.String(),
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
	}, nil
}

func (s *Server) EmailConfirmation(ctx context.Context, in *users.EmailConfirmationRequest) (*users.OperationStatus, error) {
	var code, expiresAt int64
	err := s.DB.QueryRow(`SELECT code, expires_at FROM email_confirmation WHERE account_id = $1 LIMIT 1;`, in.AccountID).Scan(
		&code, &expiresAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("student")
		}
		return nil, ErrInternal
	}

	if expiresAt-time.Now().Unix() < 0 {
		return nil, ErrInvalidInput("code", "expires")
	}

	if code != in.Code {
		return nil, ErrInvalidInput("code", "wrong")
	}

	_, err = s.DB.Exec(`UPDATE accounts SET status = 'active' WHERE  account_id = $1;`, in.AccountID)
	if err != nil {
		return nil, ErrInternal
	}
	_, err = s.DB.Exec(`DELETE FROM email_confirmation WHERE account_id = $1;`, in.AccountID)
	if err != nil {
		return nil, ErrInternal
	}

	return &users.OperationStatus{
		Completed: true,
	}, nil
}
