package resolvers

import (
	"database/sql"

	"github.com/les-cours/user-service/api/auth"
	"github.com/sendgrid/sendgrid-go"
)

var instance *Server

func GetInstance(SQLDB *sql.DB, authService auth.AuthServiceClient, sendgridClient *sendgrid.Client) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:             SQLDB,
		AuthService:    authService,
		SendGridClient: sendgridClient,
	}
	return instance
}
