package resolvers

import (
	"context"
	"github.com/les-cours/user-service/api/users"
	"github.com/les-cours/user-service/utils"
)

func (s *Server) InviteTeacher(ctx context.Context, in *users.InviteTeacherRequest) (*users.OperationStatus, error) {
	teacherID := utils.GenerateUUIDString()
	var subjectsString = in.Subjects[0]
	for i := 1; i < len(in.Subjects); i++ {
		subjectsString += "," + in.Subjects[i]
	}
	_, err := s.DB.Exec(`
INSERT INTO 
    teachers_invitations 
    (teacher_id, email, subjects) 
VALUES ($1,$2,$3)`, teacherID, in.Email, subjectsString)
	if err != nil {
		return nil, err
	}

	link := "localhost:3001/teacher/confirm?agentID=" + teacherID
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
		s.Logger.Error(err.Error())
		return nil, err
	}

	return &users.OperationStatus{
		Completed: true,
	}, nil
}

func (s *Server) TeacherSignup(ctx context.Context, in *users.TeacherSignupRequest) (*users.TeacherSignupResponse, error) {

	return nil, nil
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
