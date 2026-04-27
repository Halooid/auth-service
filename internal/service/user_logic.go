package service

import (
	"context"
	"fmt"

	authv1 "github.com/halooid/backend/auth-service/gen/go/auth/v1"
	"github.com/halooid/backend/go-shared/auth"
)

// UserBusinessLogic handles user domain operations
type UserBusinessLogic struct {
	validator *auth.Validator
}

func NewUserBusinessLogic(validator *auth.Validator) *UserBusinessLogic {
	return &UserBusinessLogic{
		validator: validator,
	}
}

func (s *UserBusinessLogic) ValidateToken(ctx context.Context, token string) (*authv1.User, error) {
	claims, err := s.validator.Validate(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	return &authv1.User{
		Id:         claims.Subject, // Use Keycloak 'sub' as internal ID now
		ExternalId: claims.Subject,
		Email:      claims.Email,
		Username:   claims.PreferredUsername,
		TenantId:   claims.TenantID,
		FirstName:  claims.GivenName,
		LastName:   claims.FamilyName,
	}, nil
}

func (s *UserBusinessLogic) GetCurrentUser(ctx context.Context, token string) (*authv1.User, error) {
	return s.ValidateToken(ctx, token)
}
