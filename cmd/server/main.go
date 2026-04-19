package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	authv1 "github.com/halooid/backend/auth-service/gen/go/auth/v1"
	"github.com/halooid/backend/auth-service/internal/db"
	"github.com/halooid/backend/auth-service/internal/handler"
	"github.com/halooid/backend/auth-service/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	queries := db.New(conn)
	logic := service.NewUserBusinessLogic(queries)
	authHandler := handler.NewAuthHandler(logic)

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)

	// Register reflection for debugging
	reflection.Register(grpcServer)

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Auth Service starting on gRPC port %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
