package handler

import (
	"context"

	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	"github.com/stretchr/testify/mock"
)

type StorageMock struct {
	mock.Mock
}

func (s *StorageMock) GetAnimal(ctx context.Context, req *pb.AnimalReq) (*pb.AnimalRsp, error) {
	arguments := s.Called(ctx, req)
	if arguments.Get(0) != nil {
		return arguments.Get(0).(*pb.AnimalRsp), arguments.Error(1)
	}
	return nil, arguments.Error(1)
}

func NewStorageMock() *StorageMock {
	return &StorageMock{}
}
