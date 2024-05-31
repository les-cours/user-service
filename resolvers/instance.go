package resolvers

import (
	"database/sql"
	"github.com/les-cours/user-service/api/learning"
	"github.com/les-cours/user-service/api/payment"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
)

var instance *Server

func GetInstance(SQLDB *sql.DB, authService auth.AuthServiceClient, learningService learning.LearningServiceClient, paymentService payment.PaymentServiceClient, logger *zap.Logger) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:              SQLDB,
		AuthService:     authService,
		LearningService: learningService,
		PaymentService:  paymentService,
		Logger:          logger,
	}
	return instance
}
