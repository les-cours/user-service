package resolvers

import (
	"database/sql"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
	"github.com/sendgrid/sendgrid-go"
)

var instance *Server

func GetInstance(SQLDB *sql.DB, authService auth.AuthServiceClient, sendgridClient *sendgrid.Client, logger *zap.Logger) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:             SQLDB,
		AuthService:    authService,
		SendGridClient: sendgridClient,
		Logger:         logger,
	}
	return instance
}
