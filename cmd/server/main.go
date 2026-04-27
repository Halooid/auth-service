package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	authv1 "github.com/halooid/backend/auth-service/gen/go/auth/v1"
	"github.com/halooid/backend/auth-service/internal/handler"
	"github.com/halooid/backend/auth-service/internal/service"
	"github.com/halooid/backend/go-shared/auth"
	"github.com/halooid/backend/go-shared/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	keycloakBaseURL := os.Getenv("KEYCLOAK_BASE_URL")
	keycloakRealm := os.Getenv("KEYCLOAK_REALM")

	if keycloakBaseURL == "" || keycloakRealm == "" {
		log.Fatal("KEYCLOAK_BASE_URL and KEYCLOAK_REALM are required")
	}

	ctx := context.Background()
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakBaseURL, keycloakRealm)
	validator, err := auth.NewValidator(ctx, jwksURL, 15*time.Minute)
	if err != nil {
		log.Fatalf("failed to initialize auth validator: %v", err)
	}

	logic := service.NewUserBusinessLogic(validator)
	authHandler := handler.NewAuthHandler(logic)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logging.UnaryServerInterceptor()),
	)
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
