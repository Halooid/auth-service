package handler

import (
	"context"
	"strings"

	authv1 "github.com/halooid/backend/auth-service/gen/go/auth/v1"
	"github.com/halooid/backend/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	resp, err := h.logic.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &authv1.RegisterResponse{
		User: &authv1.User{
			Id:         resp.User.ID.String(),
			ExternalId: resp.User.ExternalID,
			Email:      resp.User.Email,
			Username:   resp.User.Username,
			TenantId:   resp.User.TenantID.String,
		},
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		ExpiresIn:        1, // We don't have these in RegisterResult yet, but could add. For now, matching proto.
		RefreshExpiresIn: 1,
		TokenType:        "Bearer",
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

func (h *AuthHandler) GetCurrentUser(ctx context.Context, req *authv1.GetCurrentUserRequest) (*authv1.GetCurrentUserResponse, error) {
	// Extract token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "metadata missing")
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "authorization header missing")
	}

	// Assuming "Bearer <token>"
	token := strings.TrimPrefix(authHeader[0], "Bearer ")

	user, err := h.logic.GetCurrentUser(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to get current user: %v", err)
	}

	return &authv1.GetCurrentUserResponse{
		User: &authv1.User{
			Id:         user.ID.String(),
			ExternalId: user.ExternalID,
			Email:      user.Email,
			Username:   user.Username,
			TenantId:   user.TenantID.String,
		},
	}, nil
}
