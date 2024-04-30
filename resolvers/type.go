package resolvers

import (
	"database/sql"
	"github.com/les-cours/user-service/api/learning"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
)

type Server struct {
	DB              *sql.DB
	AuthService     auth.AuthServiceClient
	LearningService learning.LearningServiceClient
	Logger          *zap.Logger
	users.UnimplementedUserServiceServer
}
