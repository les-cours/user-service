package resolvers

import (
	"database/sql"
	"github.com/les-cours/user-service/api/learning"
	"github.com/les-cours/user-service/api/payment"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
	"github.com/les-cours/user-service/api/users"
)

type Server struct {
	DB              *sql.DB
	AuthService     auth.AuthServiceClient
	LearningService learning.LearningServiceClient
	PaymentService  payment.PaymentServiceClient
	Logger          *zap.Logger
	users.UnimplementedUserServiceServer
}
