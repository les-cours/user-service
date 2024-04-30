package resolvers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/learning"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/env"
	"github.com/les-cours/user-service/utils"
	"strings"
)

func (s *Server) InviteTeacher(ctx context.Context, in *users.InviteTeacherRequest) (*users.OperationStatus, error) {
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

	link := env.Settings.TeacherConfirmEndPoint + teacherID
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

	err := s.DB.QueryRow(`SELECT email,subjects FROM teachers_invitations WHERE teacher_id = $1;`, in.TeacherID).Scan(&email, &subjectsString)
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

	var userName = "T_" + in.Firstname + "_" + in.Lastname
	_, err = tx.ExecContext(context.Background(),
		`INSERT INTO accounts 
		(account_id,email,password,username, status, user_type) VALUES($1, $2,crypt($3,gen_salt('bf')),$4,$5,$6);`,
		in.TeacherID, email, in.GetPassword(), userName, "active", "teacher")

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
VALUES ($1,$2,$3)`, in.TeacherID, in.Firstname, in.Lastname, avatar)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	/*
		GENERATE CLASSROOMS FOR HER/HIS SUBJECTS ...
	*/

	subjects := strings.Split(subjectsString, ",")
	_, err = s.LearningService.CreateClassRooms(ctx, &learning.CreateClassRoomsRequest{
		TeacherID:  in.TeacherID,
		SubjectIDs: subjects,
	})

	if err != nil {
		tx.Rollback()
		s.Logger.Error(err.Error())
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
    false,       -- Example values for orgs_read
    false,       -- Example values for users_create
    false,       -- Example values for users_update
    false,       -- Example values for users_delete
    false,      -- Example values for users_read
    true,       -- Example values for learning_create
    true,      -- Example values for learning_update
    true,      -- Example values for learning_delete
    true        -- Example values for learning_read
);`, in.TeacherID)
	if err != nil {
		s.Logger.Error(err.Error())
		tx.Rollback()
		return nil, ErrInternal
	}

	/*
		Commit
	*/
	err = tx.Commit()
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}

	/*
		DELETE INVITATION
	*/

	_, err = s.DB.Exec(`DELETE FROM teachers_invitations where teacher_id = $1;`, in.TeacherID)
	if err != nil {
		s.Logger.Error(err.Error())
	}

	res, err := s.AuthService.Signup(ctx, &auth.SignUpRequest{
		AccountID: in.TeacherID,
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

func (s *Server) GetTeachersBySubject(ctx context.Context, in *users.GetTeacherBySubjectRequest) ([]*users.Teacher, error) {

	rows, err := s.DB.Query(`
SELECT 
    t.teacher_id,t.firstname,t.lastname
FROM teachers as t
    INNER JOIN 
    public.grades_subjects gs 
        on t.teacher_id = gs.grade_id
WHERE gs.subject_id = $1;
        `, in.SubjectID)

	var teacher *users.Teacher
	var teachers []*users.Teacher
	for rows.Next() {
		err = rows.Scan(&teacher.TeacherID, &teacher.Firstname, &teacher.Lastname)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}
