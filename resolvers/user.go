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

		var orgsCreate, orgsUpdate, orgsDelete, orgsRead bool
		var usersCreate, usersUpdate, usersDelete, usersRead bool
		var learningCreate, learningUpdate, learningDelete, learningRead bool

		err = s.DB.QueryRow(`SELECT  orgs_create, orgs_update, orgs_delete, orgs_read,
               users_create, users_update, users_delete, users_read,
               learning_create, learning_update, learning_delete, learning_read
        FROM permissions
	WHERE account_id = $1;
	`, &user.AccountID).Scan(
			&orgsCreate,
			&orgsUpdate,
			&orgsDelete,
			&orgsRead,
			&usersCreate,
			&usersUpdate,
			&usersDelete,
			&usersRead,
			&learningCreate,
			&learningUpdate,
			&learningDelete,
			&learningRead,
		)

		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}

		return &users.User{
			Id:        user.AccountID,
			AccountID: user.AccountID,
			Username:  username,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Avatar:    avatar,
			UserType:  userType,
			CREATE: &users.Permissions{
				Orgs:     orgsCreate,
				Learning: learningCreate,
				Users:    usersCreate,
			},
			READ: &users.Permissions{
				Orgs:     orgsRead,
				Learning: learningRead,
				Users:    usersRead,
			},
			UPDATE: &users.Permissions{
				Orgs:     orgsUpdate,
				Learning: learningUpdate,
				Users:    usersUpdate,
			},
			DELETE: &users.Permissions{
				Orgs:     orgsDelete,
				Learning: learningDelete,
				Users:    usersDelete,
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
