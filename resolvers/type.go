package resolvers

import (
	"database/sql"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
	"github.com/sendgrid/sendgrid-go"
)

type Server struct {
	DB             *sql.DB
	AuthService    auth.AuthServiceClient
	SendGridClient *sendgrid.Client
	Logger         *zap.Logger
	users.UnimplementedUserServiceServer
}
