package resolvers

import (
	"context"
	"github.com/les-cours/user-service/api/users"
)

func InviteTeacher(ctx context.Context, in *users.InviteTeacherRequest) error {
	/*
		email
		subjects
	*/

	return nil
}

func (s *Server) GetTeachersBySubject(ctx context.Context, in *users.GetTeacherBySubjectRequest) ([]*users.Teacher, error) {

	rows, err := s.DB.Query(`
SELECT 
    t.teacher_id,t.username,t.firstname,t.lastname
FROM teachers as t
    INNER JOIN 
    public.grades_subjects gs 
        on t.teacher_id = gs.grade_id
WHERE gs.subject_id = $1;
        `, in.SubjectID)

	var teacher *users.Teacher
	var teachers []*users.Teacher
	for rows.Next() {
		err = rows.Scan(&teacher.TeacherID, &teacher.Username, &teacher.Firstname, &teacher.Lastname)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}
