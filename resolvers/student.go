package resolvers

import (
	"database/sql"
	"errors"
	"github.com/les-cours/user-service/api/learning"
	"github.com/les-cours/user-service/api/users"
	"golang.org/x/net/context"
)

func (s *Server) InitStudent(ctx context.Context, in *users.IDRequest) (*users.Notifications, error) {
	res, err := s.LearningService.InitClassRooms(ctx, &learning.IDRequest{
		Id:     in.Id,
		UserID: in.Id,
	})
	if err != nil {
		return nil, err
	}

	var notifications = make([]*users.Notification, 0)
	for _, notification := range res.Notifications {
		notifications = append(notifications, &users.Notification{
			Id:      notification.Id,
			Title:   notification.Title,
			Content: notification.Content,
		})
	}

	return &users.Notifications{
		Notifications: notifications,
	}, nil
}

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
	students := &users.Students{}
	student := &users.Student{}

	for rows.Next() {
		err = rows.Scan(&student.StudentID, &student.Firstname, &student.Lastname, &student.Username, &student.DateOfBirth, &student.Gender, &student.Status, &student.Avatar, &student.NotificationStatus, &student.OnlineStatus, &student.CityID)
		if err != nil {
			s.Logger.Error(err.Error())
			return nil, ErrInternal
		}
		students.Students = append(students.Students, student)
	}
	return students, nil

}

func (s *Server) GetStudent(ctx context.Context, in *users.GetStudentRequest) (*users.Student, error) {
	s.Logger.Info("starting GetStudent ...")
	student := &users.Student{}

	err := s.DB.QueryRow(`
select 
    firstname,lastname,a.username,gender,a.status,avatar,notification_status,online_status,city_id,date_of_birth
    from  
        students inner join accounts a on a.account_id = students.student_id where student_id = $1;`, in.StudentID).Scan(&student.Firstname, &student.Lastname, &student.Username, &student.Gender, &student.Status, &student.Avatar, &student.NotificationStatus, &student.OnlineStatus, &student.CityID, &student.DateOfBirth)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound("student")
		}
		s.Logger.Error(err.Error())
		return nil, ErrInternal
	}
	student.StudentID = in.StudentID
	return student, nil
}

const Query = `select 
   student_id,firstname,lastname,a.username,date_of_birth,gender,a.status,avatar,notification_status,online_status,city_id
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
