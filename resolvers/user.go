package resolvers

import (
	"database/sql"
	"errors"
	"log"

	"github.com/les-cours/user-service/api/users"
	"golang.org/x/net/context"
)

func (s *Server) GetUser(ctx context.Context, in *users.GetUserRequest) (*users.User, error) {

	var accountID string

	err := s.DB.QueryRow(`SELECT account_id FROM accounts WHERE email = $1 AND password = crypt($2,password)`,
		in.GetUsername(),
		in.GetPassword(),
	).Scan(&accountID)

	if err != nil {
		s.Logger.Error(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, Err("wrong password or email")
		}
		return nil, err
	}

	var role = "student"
	if in.IsTeacher {
		role = "teacher"
	}
	user, err := s.GetUserByID(ctx, &users.GetUserByIDRequest{
		AccountID: accountID,
		UserRole:  role,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return user, nil
}

func (s *Server) GetUserByID(ctx context.Context, user *users.GetUserByIDRequest) (*users.User, error) {
	var (
		username,
		firstname,
		lastname,
		email,
		avatar,
		userType,
		accountStatus string
		err error
	)
	if user.UserRole == "teacher" {

		err = s.DB.QueryRow(
			`
		SELECT 
		accounts.username, accounts.email, teachers.firstname, teachers.lastname,  teachers.avatar , accounts.status, accounts.user_type
		FROM 
		teachers
		INNER JOIN accounts ON accounts.account_id =teachers.teacher_id
		WHERE
	 	(accounts.account_id = $1)
		`, &user.AccountID).Scan(
			&username,
			&email,
			&firstname,
			&lastname,
			&avatar,
			&accountStatus,
			&userType)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}

		var writeComment, live, upload bool
		err = s.DB.QueryRow(`
	SELECT write_comment,live,upload FROM permissions
	WHERE account_id = $1;
	`, &user.AccountID).Scan(&writeComment, &live, &upload)
		if err != nil {
			return nil, err
		}

		//get plan (TO DO  )
		//ADD TABLE SUBSCRIPTION AND GET PLAN INFORMATION

		return &users.User{
			Id:        user.AccountID,
			AccountID: user.AccountID,
			Username:  username,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Avatar:    avatar,
			UserType:  userType,
			Permissions: &users.Permissions{
				WriteComment: writeComment,
				Live:         live,
				Upload:       upload,
			},
		}, nil
	}

	//student ...
	err = s.DB.QueryRow(
		`
		SELECT 
		accounts.username, accounts.email, students.firstname, students.lastname,  students.avatar , accounts.status,accounts.user_type
		FROM 
		students
		INNER JOIN accounts ON accounts.account_id = students.student_id
		WHERE
	 	(accounts.account_id = $1)
		`, &user.AccountID).Scan(
		&username,
		&email,
		&firstname,
		&lastname,
		&avatar,
		&accountStatus,
		&userType)
	if err != nil {
		log.Println("err SELECT students", err)
		if errors.Is(err, sql.ErrNoRows) {

			return nil, ErrInvalidInput("id", "doesn't exist")
		}
		return nil, err
	}

	//get plan (TO DO  )
	//ADD TABLE SUBSCRIPTION AND GET PLAN INFORMATION

	return &users.User{
		Id:        user.AccountID,
		AccountID: user.AccountID,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Avatar:    avatar,
		UserType:  userType,
	}, nil
}
