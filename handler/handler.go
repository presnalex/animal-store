package handler

import (
	"context"
	"errors"
	"net/http"

	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	errs "github.com/presnalex/animal-store/errors"
	log "github.com/presnalex/go-micro/v3/logger"
)

type Handler struct {
	writer  IWriter
	storage IStorage
}

type IStorage interface {
	GetAnimal(ctx context.Context, req *pb.AnimalReq) (*pb.AnimalRsp, error)
}
type IWriter interface {
	Response(tx context.Context, rw http.ResponseWriter, value interface{})
}

func (s *Handler) GetAnimal(w http.ResponseWriter, r *http.Request) {
	logger := log.FromIncomingContext(r.Context())
	logger.Debugf(r.Context(), "GetAnimal request started: %+v", r)

	request := new(pb.AnimalReq)

	q := r.URL.Query()

	animalId := q.Get("animalId")
	if animalId == "" {
		s.writer.Response(r.Context(), w, errs.NewBadRequestError(errors.New("empty param animalId")))
		return
	}

	request.AnimalId = animalId

	response, err := s.storage.GetAnimal(r.Context(), request)
	if err != nil {
		s.writer.Response(r.Context(), w, err)
		return
	}

	s.writer.Response(r.Context(), w, response)
	return
}

func NewHandler(writer IWriter, storage IStorage) (*Handler, error) {
	if writer == nil {
		return nil, errors.New("empty writer")
	}
	return &Handler{
		writer:  writer,
		storage: storage,
	}, nil
}
