package storage

import (
	"context"

	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	errs "github.com/presnalex/animal-store/errors"
	dbwrapper "github.com/presnalex/go-micro/v3/database/wrapper"
	log "github.com/presnalex/go-micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
)

type IDbWrapper interface {
	SelectContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error
}

type Storage struct {
	DBPrimary IDbWrapper
}

func (s *Storage) GetAnimal(ctx context.Context, req *pb.AnimalReq) (*pb.AnimalRsp, error) {
	logger := log.FromIncomingContext(ctx)
	logger.Debug(ctx, "GetAnimal request")

	md, _ := metadata.FromIncomingContext(ctx)
	ctx = metadata.NewOutgoingContext(ctx, md)

	dbResponse := make([]*AnimalResponse, 0)

	err := s.DBPrimary.SelectContext(dbwrapper.QueryContext(ctx, "getAnimal"), &dbResponse, getAnimal, req.AnimalId)
	if err != nil {
		logger.Errorf(ctx, "unable to get animal from db, lol: %v", err)
		return nil, errs.NewInternalError(err)
	}
	if len(dbResponse) == 0 {
		logger.Debug(ctx, "empty response from db")
		return nil, errs.NewNotFoundError("animal not found")
	}

	response, err := s.mappingAnimal(ctx, dbResponse)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// NewStorage ...
func NewStorage(db IDbWrapper) *Storage {
	return &Storage{DBPrimary: db}
}
