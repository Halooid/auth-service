package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/halooid/backend/auth-service/internal/db"
	"github.com/halooid/backend/go-shared/auth"
	"github.com/google/uuid"
)

// UserBusinessLogic handles user domain operations
type UserBusinessLogic struct {
	repo      db.Querier
	keycloak  *KeycloakClient
	validator *auth.Validator
}

func NewUserBusinessLogic(repo db.Querier, keycloak *KeycloakClient, validator *auth.Validator) *UserBusinessLogic {
	return &UserBusinessLogic{
		repo:      repo,
		keycloak:  keycloak,
		validator: validator,
	}
}

func (s *UserBusinessLogic) GetOrCreateUser(ctx context.Context, extID, email, username, tenantID string) (*db.User, error) {
	user, err := s.repo.GetUserByExternalID(ctx, extID)
	if err == nil {
		return &user, nil
	}

	// For simplicity, if not found, create. In production, check error type.
	newUser, err := s.repo.CreateUser(ctx, db.CreateUserParams{
		ExternalID: extID,
		Email:      email,
		Username:   username,
		TenantID:   sql.NullString{String: tenantID, Valid: tenantID != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}

func (s *UserBusinessLogic) GetUserByID(ctx context.Context, id string) (*db.User, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}

func (s *UserBusinessLogic) Login(ctx context.Context, username, password string) (*TokenResponse, error) {
	return s.keycloak.Login(ctx, username, password)
}

type RegisterResult struct {
	User         *db.User
	AccessToken  string
	RefreshToken string
}

func (s *UserBusinessLogic) Register(ctx context.Context, username, email, password string) (*RegisterResult, error) {
	// 1. Create in Keycloak
	extID, err := s.keycloak.Register(ctx, username, email, password)
	if err != nil {
		return nil, fmt.Errorf("keycloak registration failed: %w", err)
	}

	// 2. Sync to local DB
	user, err := s.GetOrCreateUser(ctx, extID, email, username, "")
	if err != nil {
		return nil, fmt.Errorf("local sync failed: %w", err)
	}

	// 3. Log in to get tokens
	tokens, err := s.keycloak.Login(ctx, username, password)
	if err != nil {
		return nil, fmt.Errorf("login after registration failed: %w", err)
	}

	return &RegisterResult{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *UserBusinessLogic) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	return s.keycloak.RefreshToken(ctx, refreshToken)
}

func (s *UserBusinessLogic) ValidateToken(ctx context.Context, token string) (*db.User, error) {
	claims, err := s.validator.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Sync or get user
	return s.GetOrCreateUser(ctx, claims.Subject, claims.Email, claims.PreferredUsername, claims.TenantID)
}

func (s *UserBusinessLogic) GetCurrentUser(ctx context.Context, token string) (*db.User, error) {
	return s.ValidateToken(ctx, token)
}
