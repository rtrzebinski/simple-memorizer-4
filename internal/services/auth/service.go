package auth

import (
	"context"
	"fmt"

	"github.com/Nerzal/gocloak/v13"
)

type Config struct {
	Realm        string
	ClientID     string
	ClientSecret string
}

type Tokens struct {
	AccessToken      string
	IDToken          string
	ExpiresIn        int
	RefreshExpiresIn int
	RefreshToken     string
	TokenType        string
}

type Service struct {
	kc  *gocloak.GoCloak
	cfg Config
}

func NewService(kc *gocloak.GoCloak, cfg Config) *Service {
	return &Service{
		kc:  kc,
		cfg: cfg,
	}
}

func (s *Service) Register(ctx context.Context, firstName, lastName, email, password string) (Tokens, error) {
	adminTok, err := s.kc.LoginClient(ctx, s.cfg.ClientID, s.cfg.ClientSecret, s.cfg.Realm)
	if err != nil {
		return Tokens{}, fmt.Errorf("login client: %w", err)
	}

	empty := []string{}
	u := gocloak.User{
		Username:        gocloak.StringP(email),
		Email:           gocloak.StringP(email),
		FirstName:       gocloak.StringP(firstName),
		LastName:        gocloak.StringP(lastName),
		Enabled:         gocloak.BoolP(true),
		EmailVerified:   gocloak.BoolP(true),
		RequiredActions: &empty,
	}

	userID, err := s.kc.CreateUser(ctx, adminTok.AccessToken, s.cfg.Realm, u)
	if err != nil {
		return Tokens{}, fmt.Errorf("create user: %w", err)
	}

	err = s.kc.SetPassword(ctx, adminTok.AccessToken, userID, s.cfg.Realm, password, false)
	if err != nil {
		return Tokens{}, fmt.Errorf("set password: %w", err)
	}

	err = s.kc.UpdateUser(ctx, adminTok.AccessToken, s.cfg.Realm, gocloak.User{
		ID:              &userID,
		RequiredActions: &empty,
		EmailVerified:   gocloak.BoolP(true),
		Enabled:         gocloak.BoolP(true),
	})
	if err != nil {
		return Tokens{}, fmt.Errorf("clear required actions: %w", err)
	}

	t, err := s.kc.Login(ctx, s.cfg.ClientID, s.cfg.ClientSecret, s.cfg.Realm, email, password)
	if err != nil {
		return Tokens{}, fmt.Errorf("login new user: %w", err)
	}

	return Tokens{
		AccessToken:      t.AccessToken,
		IDToken:          t.IDToken,
		ExpiresIn:        t.ExpiresIn,
		RefreshExpiresIn: t.RefreshExpiresIn,
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

func (s *Service) SignIn(ctx context.Context, email, password string) (Tokens, error) {
	t, err := s.kc.Login(ctx, s.cfg.ClientID, s.cfg.ClientSecret, s.cfg.Realm, email, password)
	if err != nil {
		return Tokens{}, fmt.Errorf("login user: %w", err)
	}

	return Tokens{
		AccessToken:      t.AccessToken,
		IDToken:          t.IDToken,
		ExpiresIn:        t.ExpiresIn,
		RefreshExpiresIn: t.RefreshExpiresIn,
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (Tokens, error) {
	t, err := s.kc.RefreshToken(ctx, refreshToken, s.cfg.ClientID, s.cfg.ClientSecret, s.cfg.Realm)
	if err != nil {
		return Tokens{}, fmt.Errorf("refresh: %w", err)
	}

	return Tokens{
		AccessToken:      t.AccessToken,
		IDToken:          t.IDToken,
		ExpiresIn:        t.ExpiresIn,
		RefreshExpiresIn: t.RefreshExpiresIn,
		RefreshToken:     t.RefreshToken,
		TokenType:        t.TokenType,
	}, nil
}

// Revoke the refresh token and all access tokens derived from it
func (s *Service) Revoke(ctx context.Context, refreshToken string) error {
	err := s.kc.RevokeToken(ctx, s.cfg.Realm, s.cfg.ClientID, s.cfg.ClientSecret, refreshToken)
	if err != nil {
		return fmt.Errorf("revoke: %w", err)
	}

	return nil
}
