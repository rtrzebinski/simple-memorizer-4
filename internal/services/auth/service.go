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
	"golang.org/x/crypto/bcrypt"
)

const (
	pkPath       = "./../../../keys/private.pem"
	daysToExpire = 30
)

type Service struct {
	r Reader
	w Writer
}

func NewService(r Reader, w Writer) *Service {
	return &Service{
		r: r,
		w: w,
	}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (accessToken string, err error) {
	privateKey, err := pk()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	userID, err := s.w.Register(ctx, name, email, string(hashed))
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   userID,
		"name":  name,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * daysToExpire).Unix(),
	})

	accessToken, err = token.SignedString(privateKey)
	if err != nil {
		return accessToken, fmt.Errorf("failed to sign token: %w", err)
	}

	return accessToken, nil
}

func (s *Service) SignIn(ctx context.Context, email, password string) (accessToken string, err error) {
	privateKey, err := pk()
	if err != nil {
		return "", fmt.Errorf("failed to get private key: %w", err)
	}

	name, userID, err := s.r.SignIn(ctx, email, password)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   userID,
		"name":  name,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * daysToExpire).Unix(),
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
