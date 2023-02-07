package storage

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
)

func (m *Storage) mappingAnimal(ctx context.Context, dbResponse []*AnimalResponse) (*pb.AnimalRsp, error) {
	response := new(pb.AnimalRsp)

	for i := range dbResponse {
		if dbResponse[i].AnimalId.Valid {
			response.AnimalId = &wrappers.Int32Value{Value: dbResponse[i].AnimalId.Int32}
		}
		if dbResponse[i].Animal.Valid {
			response.Animal = &wrappers.StringValue{Value: dbResponse[i].Animal.String}
		}
		if dbResponse[i].Price.Valid {
			response.Price = &wrappers.Int32Value{Value: dbResponse[i].Price.Int32}
		}
	}
	return response, nil
}
