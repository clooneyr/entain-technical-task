package main

import (
	"database/sql"
	"log"
	"net"

	"git.neds.sh/matty/entain/sports/db"
	"git.neds.sh/matty/entain/sports/proto/sports"
	"git.neds.sh/matty/entain/sports/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func main() {
	// Initialize database connection
	dbConn, err := sql.Open("sqlite3", "./sports.db")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repository
	eventsRepo := db.NewEventsRepo(dbConn)
	if err := eventsRepo.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize service
	sportsService := service.NewSportsService(eventsRepo)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	sports.RegisterSportsServer(s, sportsService)

	log.Printf("gRPC server listening on: %s", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
