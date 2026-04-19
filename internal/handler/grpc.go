package handler

import (
	"context"

	authv1 "github.com/halooid/backend/auth-service/gen/go/auth/v1"
	"github.com/halooid/backend/auth-service/internal/service"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	logic *service.UserBusinessLogic
}

func NewAuthHandler(logic *service.UserBusinessLogic) *AuthHandler {
	return &AuthHandler{logic: logic}
}

func (h *AuthHandler) GetOrCreateUser(ctx context.Context, req *authv1.GetOrCreateUserRequest) (*authv1.GetOrCreateUserResponse, error) {
	user, err := h.logic.GetOrCreateUser(ctx, req.ExternalId, req.Email, req.Username, req.TenantId)
	if err != nil {
		return nil, err
	}

	return &authv1.GetOrCreateUserResponse{
		User: &authv1.User{
			Id:         user.ID.String(),
			ExternalId: user.ExternalID,
			Email:      user.Email,
			Username:   user.Username,
			TenantId:   user.TenantID.String,
		},
	}, nil
}

// Implement other methods ...
func (h *AuthHandler) GetUserById(ctx context.Context, req *authv1.GetUserByIdRequest) (*authv1.GetUserByIdResponse, error) {
	user, err := h.logic.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &authv1.GetUserByIdResponse{
		User: &authv1.User{
			Id:         user.ID.String(),
			ExternalId: user.ExternalID,
			Email:      user.Email,
			Username:   user.Username,
			TenantId:   user.TenantID.String,
		},
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	resp, err := h.logic.Login(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		ExpiresIn:        resp.ExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
		TokenType:        resp.TokenType,
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	resp, err := h.logic.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		ExpiresIn:        resp.ExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
		TokenType:        resp.TokenType,
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	user, err := h.logic.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid: true,
		User: &authv1.User{
			Id:         user.ID.String(),
			ExternalId: user.ExternalID,
			Email:      user.Email,
			Username:   user.Username,
			TenantId:   user.TenantID.String,
		},
	}, nil
}
