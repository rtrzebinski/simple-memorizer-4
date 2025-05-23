package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	gengrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client gengrpc.AuthServiceClient
	once   sync.Once
)

func startServer() {
	once.Do(func() {
		go main()
	})
}

func waitForServer(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := grpc.NewClient(
			"localhost:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(50 * time.Millisecond)
	}

	return context.DeadlineExceeded
}

func setupClient() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client = gengrpc.NewAuthServiceClient(conn)

	return conn, nil
}

func TestMain(m *testing.M) {
	// use dummies
	os.Setenv("DUMMIES", "TRUE")

	startServer()

	err := waitForServer(5 * time.Second)
	if err != nil {
		log.Fatalf("gRPC server did not start in time: %v", err)
	}

	conn, err := setupClient()
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	os.Exit(m.Run())
}

func TestSignInCall(t *testing.T) {
	ctx := context.Background()

	req := &gengrpc.SignInRequest{
		Email:    "foo@bar.com",
		Password: "password",
	}

	resp, err := client.SignIn(ctx, req)
	if err != nil {
		t.Fatalf("gRPC SignUp call failed: %v", err)
	}

	publicKeyBytes, _ := os.ReadFile("./../../keys/public.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	token, err := jwt.Parse(resp.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		slog.Info("decoded claims:", "sub", claims["sub"], "name", claims["name"], "email", claims["email"])
	} else {
		t.Fatalf("dnvalid token: %v", err)
	}

	assert.Equal(t, req.Email, claims["email"])
	assert.Equal(t, "name", claims["name"])
}

func TestRegisterCall(t *testing.T) {
	ctx := context.Background()

	req := &gengrpc.RegisterRequest{
		Email:    "foo@bar.com",
		Name:     "foo bar",
		Password: "password",
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		t.Fatalf("gRPC StoreUser call failed: %v", err)
	}

	publicKeyBytes, _ := os.ReadFile("./../../keys/public.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	token, err := jwt.Parse(resp.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		slog.Info("decoded claims:", "sub", claims["sub"], "name", claims["name"], "email", claims["email"])
	} else {
		t.Fatalf("tnvalid token: %v", err)
	}

	assert.Equal(t, req.Email, claims["email"])
	assert.Equal(t, req.Name, claims["name"])
}
