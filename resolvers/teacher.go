package resolvers

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/learning"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/env"
	"github.com/les-cours/user-service/utils"
	"strings"
	"time"
)

func (s *Server) InviteTeacher(ctx context.Context, in *users.InviteTeacherRequest) (*users.OperationStatus, error) {
	var notExist bool
	s.DB.QueryRow(`SELECT true FROM teachers_invitations WHERE email = $1 `, in.Email).Scan(&notExist)
	if notExist {
		return nil, Err("already invited")
	}
	s.DB.QueryRow(`SELECT true FROM accounts WHERE email = $1 `, in.Email).Scan(&notExist)
	if notExist {
		return nil, Err("already teacher")
	}

	teacherID := utils.GenerateUUIDString()
	var subjectsString = in.Subjects[0]
	for i := 1; i < len(in.Subjects); i++ {
		subjectsString += "," + in.Subjects[i]
	}
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}
	_, err = tx.Exec(`
INSERT INTO 
    teachers_invitations 
    (teacher_id, email, subjects) 
VALUES ($1,$2,$3)`, teacherID, in.Email, subjectsString)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	link := env.Settings.TeacherConfirmEndPoint + base64.URLEncoding.EncodeToString([]byte(teacherID))
	//Send Email :
	var emailData = struct {
		Receiver string
		Link     string
	}{
		Receiver: in.Email,
		Link:     link,
	}

	var emailSubject = "Invite Teacher"
	var emailTemplate = "teacher-invitation"

	err = utils.GenerateEmail(in.Email, emailSubject, emailTemplate, emailData)

	if err != nil {
		tx.Rollback()
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	err = tx.Commit()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	return &users.OperationStatus{
		Completed: true,
	}, nil
}

func (s *Server) TeacherSignup(ctx context.Context, in *users.TeacherSignupRequest) (*users.TeacherSignupResponse, error) {

	var subjectsString string
	var email string

	decodedID, err := base64.URLEncoding.DecodeString(in.TeacherID)
	if err != nil {
		return nil, ErrNotFound("invitation")
	}

	teacherID := string(decodedID)

	err = s.DB.QueryRow(`SELECT email,subjects FROM teachers_invitations WHERE teacher_id = $1;`, teacherID).Scan(&email, &subjectsString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, Err("you are not invited.")
		}
		return nil, err
	}

	tx, err := s.DB.BeginTx(ctx, nil)

	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	var userName = "t_" + in.Firstname + "_" + in.Lastname
	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO accounts 
		(account_id,email,password,username, status, user_type) VALUES($1, $2,crypt($3,gen_salt('bf')),$4,$5,$6);`,
		teacherID, email, in.GetPassword(), userName, "active", "teacher")

	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	avatar, err := utils.GenerateAvatar(in.Firstname, in.Lastname)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	_, err = tx.Exec(`INSERT INTO 
    teachers 
    (teacher_id, firstname, lastname,avatar) 
VALUES ($1,$2,$3,$4)`, teacherID, in.Firstname, in.Lastname, avatar)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	/*
		Permission
	*/
	_, err = tx.Exec(`
INSERT INTO permissions (
    account_id,
    orgs_create,
    orgs_update,
    orgs_delete,
    orgs_read,
    users_create,
    users_update,
    users_delete,
    users_read,
    learning_create,
    learning_update,
    learning_delete,
    learning_read
) VALUES (
    $1, -- Replace 'account1' with the actual account ID
    false,       -- Example values for orgs_create
    false,      -- Example values for orgs_update
    false,      -- Example values for orgs_delete
    true,       -- Example values for orgs_read
    false,       -- Example values for users_create
    false,       -- Example values for users_update
    false,       -- Example values for users_delete
    true,      -- Example values for users_read
    true,       -- Example values for learning_create
    true,      -- Example values for learning_update
    true,      -- Example values for learning_delete
    true        -- Example values for learning_read
);`, teacherID)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	/*
		DELETE INVITATION
	*/

	_, err = tx.Exec(`DELETE FROM teachers_invitations where teacher_id = $1;`, teacherID)
	if err != nil {
		tx.Rollback()
		s.Logger.Error(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	/*
		GENERATE CLASSROOMS FOR HER/HIS SUBJECTS ...
	*/

	subjects := strings.Split(subjectsString, ",")
	_, err = s.LearningService.CreateClassRooms(ctx, &learning.CreateClassRoomsRequest{
		TeacherID:  teacherID,
		SubjectIDs: subjects,
	})

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	res, err := s.AuthService.Signup(ctx, &auth.SignUpRequest{
		AccountID: teacherID,
		UserRole:  "teacher",
	})

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &users.TeacherSignupResponse{
		Token: res.AccessToken.Token,
	}, nil
}

func (s *Server) UpdateTeacher(ctx context.Context, in *users.UpdateTeacherRequest) (*users.Teacher, error) {

	_, err := s.DB.Exec(` UPDATE teachers
        SET city_id = $2,
            firstname = $3,
            lastname = $4,
            gender = $5,
            date_of_birth = $6,
            description = $7,
            avatar = $8
        WHERE teacher_id = $1`, in.TeacherID, in.CityID, in.Firstname, in.Lastname, in.Gender, in.DateOfBirth, in.Description, in.Avatar)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("teacher")
		}
		return nil, err
	}

	return &users.Teacher{
		TeacherID:   in.TeacherID,
		CityID:      in.CityID,
		Firstname:   in.Firstname,
		Lastname:    in.Lastname,
		Gender:      in.Gender,
		DateOfBirth: in.DateOfBirth,
		Description: in.Description,
		Avatar:      in.Avatar,
	}, nil
}

func (s *Server) DeleteTeacher(ctx context.Context, in *users.IDRequest) (*users.OperationStatus, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}
	_, err = tx.Exec(`UPDATE  teachers SET deleted_at = $2 WHERE teacher_id = $1`, in.Id, time.Now().Unix())
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	//then set all his/her classRoom ==> disable
	_, err = s.LearningService.DeleteClassRoomsByTeacher(ctx, &learning.IDRequest{
		Id: in.Id,
	})
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	err = tx.Commit()

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	return &users.OperationStatus{
		Completed: true,
	}, nil
}

func (s *Server) GetTeachers(ctx context.Context, in *users.Empty) (*users.Teachers, error) {
	//all,grad,subject,classroom
	var rows *sql.Rows
	var err error

	teachers := &users.Teachers{}

	rows, err = getTeachers(s.DB)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		var teacher = new(users.Teacher)
		err = rows.Scan(&teacher.TeacherID, &teacher.Firstname, &teacher.Lastname, &teacher.Username, &teacher.DateOfBirth, &teacher.Gender, &teacher.Status, &teacher.Avatar, &teacher.OnlineStatus, &teacher.CityID)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		teachers.Teachers = append(teachers.Teachers, teacher)
	}
	return teachers, nil

}

func (s *Server) GetTeacher(ctx context.Context, in *users.IDRequest) (*users.Teacher, error) {

	teacher := &users.Teacher{}

	err := s.DB.QueryRow(`
select 
    firstname,lastname,a.username,gender,a.status,avatar,online_status,city_id,date_of_birth
    from  
        teachers inner join accounts a on a.account_id = teachers.teacher_id where teacher_id = $1;`, in.Id).Scan(&teacher.Firstname, &teacher.Lastname, &teacher.Username, &teacher.Gender, &teacher.Status, &teacher.Avatar, &teacher.OnlineStatus, &teacher.CityID, &teacher.DateOfBirth)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("teacher")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	teacher.TeacherID = in.Id
	return teacher, nil
}

func getTeachers(db *sql.DB) (*sql.Rows, error) {
	const Query = `select 
   teacher_id,firstname,lastname,a.username,date_of_birth,gender,a.status,avatar,online_status,city_id
    from  
        teachers inner join public.accounts a on a.account_id = teachers.teacher_id`
	rows, err := db.Query(Query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
