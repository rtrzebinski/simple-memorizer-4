package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	pkPath = "./../../../keys/private.pem"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Register(_ context.Context, name, email, password string) (accessToken string, err error) {
	privateKey, err := pk()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  name,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	accessToken, err = token.SignedString(privateKey)
	if err != nil {
		return accessToken, fmt.Errorf("failed to sign token: %w", err)
	}

	return accessToken, nil
}

func (s *Service) SignIn(_ context.Context, email, password string) (accessToken string, err error) {
	privateKey, err := pk()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  "", // todo: fetch name from db
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	accessToken, err = token.SignedString(privateKey)
	if err != nil {
		return accessToken, fmt.Errorf("failed to sign token: %w", err)
	}

	return accessToken, nil
}

func pk() (*rsa.PrivateKey, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get caller info")
	}

	dir := filepath.Dir(filename)

	privateKeyBytes, err := os.ReadFile(filepath.Join(dir, pkPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	return jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
}
