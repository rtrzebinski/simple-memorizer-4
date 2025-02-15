package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	authgrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client authgrpc.AuthServiceClient
	once   sync.Once
)

func startServer() {
	once.Do(func() {
		go main()
	})
}

func waitForServer(address string, timeout time.Duration) error {
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

	client = authgrpc.NewAuthServiceClient(conn)

	return conn, nil
}

func TestMain(m *testing.M) {
	startServer()

	err := waitForServer(":50051", 1*time.Second)
	if err != nil {
		log.Fatalf("gRPC server did not start in time: %v", err)
	}

	conn, err := setupClient()
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	os.Exit(m.Run())
}

func TestRegisterCall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &authgrpc.RegisterRequest{
		Email:    "foo@bar.com",
		Name:     "foo bar",
		Password: "foobar123",
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		t.Fatalf("gRPC Register call failed: %v", err)
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
		fmt.Println("Decoded Claims:", claims)
	} else {
		t.Fatalf("Invalid token: %v", err)
	}

	assert.Equal(t, req.Email, claims["email"])
	assert.Equal(t, req.Name, claims["name"])
}

func TestSignInCall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &authgrpc.SignInRequest{
		Email:    "foo@bar.com",
		Password: "foobar123",
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
		fmt.Println("Decoded Claims:", claims)
	} else {
		t.Fatalf("Invalid token: %v", err)
	}

	assert.Equal(t, req.Email, claims["email"])
	assert.Equal(t, "", claims["name"])
}
