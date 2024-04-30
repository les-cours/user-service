package resolvers

import (
	"database/sql"
	"github.com/les-cours/user-service/api/learning"
	"go.uber.org/zap"

	"github.com/les-cours/user-service/api/auth"
)

var instance *Server

func GetInstance(SQLDB *sql.DB, authService auth.AuthServiceClient, learningService learning.LearningServiceClient, logger *zap.Logger) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB:              SQLDB,
		AuthService:     authService,
		LearningService: learningService,
		Logger:          logger,
	}
	return instance
}
