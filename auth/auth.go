package auth

import (
	"context"

	"github.com/jvitoroc/todo-go/config"
	"google.golang.org/api/idtoken"
)

type GoogleAuth struct {
	ClientID string
}

type AuthService struct {
	Google *GoogleAuth
}

func NewAuthService(cfg *config.Config) *AuthService {
	auth := &AuthService{}
	auth.AddGoogleAuth(cfg.Auth.GoogleClientID)

	return auth
}

func (auth *AuthService) AddGoogleAuth(clientId string) {
	auth.Google = &GoogleAuth{
		ClientID: clientId,
	}
}

func (g *GoogleAuth) Validate(idToken string) (*idtoken.Payload, error) {
	return idtoken.Validate(
		context.Background(),
		idToken,
		g.ClientID,
	)
}
