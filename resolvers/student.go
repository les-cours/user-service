package resolvers

import (
	"database/sql"
	"github.com/les-cours/user-service/api/users"

	"golang.org/x/net/context"
)

func (s *Server) GetStudents(ctx context.Context, in *users.GetStudentsRequest) (*users.Students, error) {
	//all,grad,subject,classroom
	var rows *sql.Rows
	var err error

	switch in.FilterType {
	case "all":
		rows, err = getStudents(s.DB)
	case "grad":
		rows, err = getStudentsByGrad(s.DB, in.FilterValue)
	case "classroom":
		rows, err = getStudentsByClassRoom(s.DB, in.FilterValue)

	default:
		return nil, ErrInvalidInput("filterValue", "invalid")
	}

	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	var students *users.Students

	var student *users.Student
	for rows.Next() {
		rows.Scan(&student.Firstname, &student.Lastname, &student.Username, &student.DateOfGirth, &student.Gender, &student.Status, &student.DefaultAvatar, &student.CityId)
		students.Students = append(students.Students)
	}
	return students, nil

}

func (s *Server) GetStudent(ctx context.Context, in *users.GetStudentRequest) (*users.Student, error) {

	var student *users.Student

	err := s.DB.QueryRow(`
select 
    firstname,lastname,a.username, date_of_birth,gender,a.status,avatar,notification_status,online_status,default_avatar,city_id
    from  
        students inner join public.accounts a on a.account_id = students.student_id;`).
		Scan(&student.Firstname, &student.Lastname, &student.Username, &student.DateOfGirth, &student.Gender, &student.Status, &student.DefaultAvatar, &student.CityId)
	if err != nil {
		s.Logger.Error(err.Error())
		return nil, err
	}
	return student, nil
}

const Query = `select 
    firstname,lastname,a.username,date_of_birth,gender,a.status,avatar,notification_status,online_status,default_avatar,city_id
    from  
        students inner join public.accounts a on a.account_id = students.student_id`

func getStudents(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query(Query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func getStudentsByGrad(db *sql.DB, gradID string) (*sql.Rows, error) {
	rows, err := db.Query(Query+" where grade_id = $1;", gradID)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func getStudentsByClassRoom(db *sql.DB, classID string) (*sql.Rows, error) {
	rows, err := db.Query(Query+"inner join classrooms_students c on c.student_id = students.student_id where classroom_id = $1;", classID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
