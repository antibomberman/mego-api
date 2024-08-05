package clients

import (
	pb "github.com/antibomberman/mego-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	User pb.UserServiceClient
	Auth pb.AuthServiceClient
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &UserClient{User: pb.NewUserServiceClient(conn), Auth: pb.NewAuthServiceClient(conn)}, nil
}
