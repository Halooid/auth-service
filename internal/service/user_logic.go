package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/halooid/backend/auth-service/internal/db"
	"github.com/google/uuid"
)

// UserBusinessLogic handles user domain operations
type UserBusinessLogic struct {
	repo db.Querier
}

func NewUserBusinessLogic(repo db.Querier) *UserBusinessLogic {
	return &UserBusinessLogic{repo: repo}
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
