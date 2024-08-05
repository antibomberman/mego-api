package clients

import (
	pb "github.com/antibomberman/mego-protos/gen/go/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostClient struct {
	pb.PostServiceClient
}

func NewPostClient(address string) (*PostClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PostClient{pb.NewPostServiceClient(conn)}, nil
}
