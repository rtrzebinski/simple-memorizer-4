package grpc

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
)

const (
	pkPath = "./../../../../keys/private.pem"
)

type Server struct {
	grpc.UnimplementedAuthServiceServer
}

func (s *Server) Register(_ context.Context, req *grpc.RegisterRequest) (*grpc.RegisterResponse, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	privateKeyBytes, _ := os.ReadFile(filepath.Join(dir, pkPath))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  req.Name,
		"email": req.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &grpc.RegisterResponse{
		AccessToken: accessToken,
	}, nil
}

func (s *Server) SignIn(_ context.Context, req *grpc.SignInRequest) (*grpc.SignInResponse, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	privateKeyBytes, _ := os.ReadFile(filepath.Join(dir, pkPath))
	privateKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "1234567890",
		"name":  "", // todo: fetch name from db
		"email": req.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &grpc.SignInResponse{
		AccessToken: accessToken,
	}, nil
}
