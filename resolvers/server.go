package resolvers

import (
	"database/sql"
	book "github.com/les-cours/user-service/protobuf/book"
)

var instance *Server

type Server struct {
	DB *sql.DB
	book.UnimplementedBookServiceServer
}

func GetInstance(db *sql.DB) *Server {
	if instance != nil {
		return instance
	}

	instance = &Server{
		DB: db,
	}

	return instance
}
