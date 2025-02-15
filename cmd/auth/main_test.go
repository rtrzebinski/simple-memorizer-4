package main

import (
	"context"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	authgrpc "github.com/rtrzebinski/simple-memorizer-4/generated/proto/grpc"
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

func TestSignUpCall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &authgrpc.SignUpRequest{
		Email:    "foo@bar.com",
		Password: "foobar",
	}

	resp, err := client.SignUp(ctx, req)
	if err != nil {
		t.Fatalf("gRPC SignUp call failed: %v", err)
	}

	t.Logf("AccessToken: %s", resp.AccessToken)
}

func TestSignInCall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &authgrpc.SignInRequest{
		Email:    "foo@bar.com",
		Password: "foobar",
	}

	resp, err := client.SignIn(ctx, req)
	if err != nil {
		t.Fatalf("gRPC SignIn call failed: %v", err)
	}

	t.Logf("AccessToken: %s", resp.AccessToken)
}
