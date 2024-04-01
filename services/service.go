package services

import (
	"log"
	"net"

	"github.com/les-cours/user-service/database"
	"github.com/les-cours/user-service/protobuf/book"
	"github.com/les-cours/user-service/resolvers"
	"google.golang.org/grpc"
)

const GrpcPort = "1113"

func Start() {
	var err error
	lis, err := net.Listen("tcp", ":"+GrpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %v: %v", GrpcPort, err)
	}

	db, err := db.StartDatabase()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var s = resolvers.GetInstance(db)

	var grpcServer = grpc.NewServer()
	book.RegisterBookServiceServer(grpcServer, s)

	log.Printf("Starting grpc server on port " + GrpcPort)

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to start gRPC server on port %v: %v", GrpcPort, err)
	}

}
