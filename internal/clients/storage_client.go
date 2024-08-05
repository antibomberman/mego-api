package clients

import (
	pb "github.com/antibomberman/mego-protos/gen/go/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StorageClient struct {
	pb.StorageServiceClient
}

func NewStorageClient(address string) (*StorageClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &StorageClient{pb.NewStorageServiceClient(conn)}, nil
}
