package service

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int32  `json:"expires_in"`
	RefreshExpiresIn int32  `json:"refresh_expires_in"`
	TokenType        string `json:"token_type"`
}

type KeycloakClient struct {
	client       *gocloak.GoCloak
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
}

func NewKeycloakClient(baseURL, realm, clientID, clientSecret string) *KeycloakClient {
	return &KeycloakClient{
		client:       gocloak.NewClient(baseURL),
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (c *KeycloakClient) Login(ctx context.Context, username, password string) (*TokenResponse, error) {
	token, err := c.client.Login(ctx, c.ClientID, c.ClientSecret, c.Realm, username, password)
	if err != nil {
		return nil, fmt.Errorf("keycloak login failed: %w", err)
	}

	return &TokenResponse{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        int32(token.ExpiresIn),
		RefreshExpiresIn: int32(token.RefreshExpiresIn),
		TokenType:        token.TokenType,
	}, nil
}

func (c *KeycloakClient) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	token, err := c.client.RefreshToken(ctx, refreshToken, c.ClientID, c.ClientSecret, c.Realm)
	if err != nil {
		return nil, fmt.Errorf("keycloak refresh failed: %w", err)
	}

	return &TokenResponse{
		AccessToken:      token.AccessToken,
		RefreshToken:     token.RefreshToken,
		ExpiresIn:        int32(token.ExpiresIn),
		RefreshExpiresIn: int32(token.RefreshExpiresIn),
		TokenType:        token.TokenType,
	}, nil
}

func (c *KeycloakClient) GetUserInfo(ctx context.Context, accessToken string) (*gocloak.UserInfo, error) {
	userInfo, err := c.client.GetUserInfo(ctx, accessToken, c.Realm)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	return userInfo, nil
}

func (c *KeycloakClient) Register(ctx context.Context, username, email, password, firstName, lastName string) (string, error) {
	// 1. Get admin token using client credentials
	token, err := c.client.GetToken(ctx, c.Realm, gocloak.TokenOptions{
		GrantType:    gocloak.StringP("client_credentials"),
		ClientID:     &c.ClientID,
		ClientSecret: &c.ClientSecret,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}

	// 2. Create User
	user := gocloak.User{
		Username:      gocloak.StringP(username),
		Email:         gocloak.StringP(email),
		Enabled:       gocloak.BoolP(true),
		FirstName:     gocloak.StringP(firstName),
		LastName:      gocloak.StringP(lastName),
		EmailVerified: gocloak.BoolP(true), // TODO: change this and send verification email
	}

	userID, err := c.client.CreateUser(ctx, token.AccessToken, c.Realm, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user in keycloak: %w", err)
	}

	// 3. Set Password
	err = c.client.SetPassword(ctx, token.AccessToken, userID, c.Realm, password, false)
	if err != nil {
		return "", fmt.Errorf("failed to set user password in keycloak: %w", err)
	}

	return userID, nil
}
