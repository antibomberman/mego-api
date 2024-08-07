package clients

import (
	pb "github.com/antibomberman/mego-protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	pb.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &AuthClient{pb.NewAuthServiceClient(conn)}, nil
}
