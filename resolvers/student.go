package resolvers

import (
	"github.com/les-cours/user-service/api/users"

	"golang.org/x/net/context"
)

func (s *Server) GetStudent(ctx context.Context, request *users.GetStudentRequest) (*users.Student, error) {

	var accountID string
	var firstName string
	err := s.DB.QueryRow(`select account_id,firstname from  students Limit 1`).Scan(&accountID, &firstName)
	if err != nil {
		return nil, err
	}
	return &users.Student{
		Id:        accountID,
		Firstname: firstName,
	}, nil
}
