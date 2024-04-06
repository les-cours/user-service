package resolvers

import (
	"database/sql"
	"errors"
	"log"

	"github.com/les-cours/user-service/api/users"
	"golang.org/x/net/context"
)

func (s *Server) GetUser(ctx context.Context, user *users.GetUserRequest) (*users.User, error) {
	var (
		userID,
		accountID,
		username,
		firstname,
		lastname,
		email,
		avatar,
		userType,
		accountStatus,
		planID string
		err error
	)

	//if user.IsTeacher {
	//	//change query to teacher
	//}

	err = s.DB.QueryRow(`SELECT account_id FROM accounts WHERE email = $1 AND password = crypt($2,password)`,
		user.GetUsername(),
		user.GetPassword(),
	).Scan(&accountID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, Err("wrong password or email")
		}
		return nil, ErrInternal
	}

	//AccountStatus AccountStatus `json:"account"`
	//Permissions   Permissions   `json:"permissions"`

	err = s.DB.QueryRow(
		`
		SELECT 
		students.id, students.username, accounts.email, students.firstname, students.lastname,  students.avatar , accounts.status, accounts.plan_id,accounts.user_type
		FROM 
		students
		INNER JOIN accounts ON accounts.account_id = students.account_id
		WHERE
	 	(accounts.account_id = $1)
		`, &accountID).Scan(
		&userID,
		&username,
		&email,
		&firstname,
		&lastname,
		&avatar,
		&accountStatus,
		&planID,
		&userType)
	if err != nil {
		log.Printf("Failed to scan agent: %v", err)
		return nil, ErrInternal
	}

	//get permision
	var writeComment, live, settings bool
	err = s.DB.QueryRow(`
	SELECT write_comment,live,settings FROM permissions
	WHERE account_id = $1;
	`, &accountID).Scan(&writeComment, live, settings)
	if err != nil {
		return nil, ErrInternal
	}

	//get plan (TO DO  )
	//ADD TABLE SUBSCRIPTION AND GET PLAN INFORMATION

	return &users.User{
		Id:        userID,
		AccountID: accountID,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Plan: &users.Plan{
			PlanID:      planID,
			Name:        "PLAN NAME HERE",
			PeriodEndAt: 0,
			Active:      false,
			Require:     "",
		},
		Avatar: avatar,
		Permissions: &users.Permissions{
			WriteComment: writeComment,
			Live:         live,
			Settings:     settings,
			AccountId:    accountID,
		},
	}, nil
}

func (s *Server) GetUserByID(ctx context.Context, user *users.GetUserByIDRequest) (*users.User, error) {
	var (
		userID,
		accountID,
		username,
		firstname,
		lastname,
		email,
		avatar,
		userType,
		accountStatus,
		planID string
		err error
	)

	//if user.IsTeacher {
	//	//change query to teacher
	//}

	err = s.DB.QueryRow(
		`
		SELECT 
		students.id, students.username, accounts.email, students.firstname, students.lastname,  students.avatar , accounts.status, accounts.plan_id,accounts.user_type
		FROM 
		students
		INNER JOIN accounts ON accounts.account_id = students.account_id
		WHERE
	 	(accounts.account_id = $1)
		`, &user.UserID).Scan(
		&userID,
		&username,
		&email,
		&firstname,
		&lastname,
		&avatar,
		&accountStatus,
		&planID,
		&userType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidInput("id", "doesn't exist")
		}
		return nil, ErrInternal
	}

	//get permision
	var writeComment, live, settings bool
	err = s.DB.QueryRow(`
	SELECT write_comment,live,settings FROM permissions
	WHERE account_id = $1;
	`, &accountID).Scan(&writeComment, live, settings)
	if err != nil {
		return nil, ErrInternal
	}

	//get plan (TO DO  )
	//ADD TABLE SUBSCRIPTION AND GET PLAN INFORMATION

	return &users.User{
		Id:        userID,
		AccountID: accountID,
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Plan: &users.Plan{
			PlanID:      planID,
			Name:        "PLAN NAME HERE",
			PeriodEndAt: 0,
			Active:      false,
			Require:     "",
		},
		Avatar: avatar,
		Permissions: &users.Permissions{
			WriteComment: writeComment,
			Live:         live,
			Settings:     settings,
			AccountId:    accountID,
		},
	}, nil
}
