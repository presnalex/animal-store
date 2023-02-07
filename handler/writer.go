package handler

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	errs "github.com/presnalex/animal-store/errors"
	log "github.com/presnalex/go-micro/v3/logger"
)

type encoder interface {
	Success(rw http.ResponseWriter, response interface{}) error
	Error(rw http.ResponseWriter, err pb.Error, status int) error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type Writer struct {
	encoder encoder
}

func NewWriter(encoder encoder) (*Writer, error) {
	if encoder == nil {
		return nil, errors.New("empty encoder")
	}

	return &Writer{
		encoder: encoder,
	}, nil
}

func (w *Writer) Response(ctx context.Context, rw http.ResponseWriter, value interface{}) {
	var err error
	logger := log.FromIncomingContext(ctx)

	if v, ok := value.(error); ok {
		err = w.error(ctx, rw, v)
	} else {
		err = w.success(rw, value)
	}

	if err != nil {
		logger.Fatal(ctx, "writer.Response", err)
	}
}

func (w *Writer) error(ctx context.Context, rw http.ResponseWriter, err error) error {
	logger := log.FromIncomingContext(ctx)
	e, status := errs.MapError(ctx, err)

	logger.Errorf(ctx, "error: %s, code: %s, http status: %d, uuid: %s", err, e.Code, status, e.Uuid)
	logger.Debugf(ctx, "error details: %v", e.Details)
	if err, ok := err.(stackTracer); ok {
		logger.Errorf(ctx, "error stacktrace: %+v, uuid: %s", err.StackTrace(), e.Uuid)
	}

	return w.encoder.Error(rw, e, status)
}

func (w *Writer) success(rw http.ResponseWriter, value interface{}) error {
	return w.encoder.Success(rw, value)
}
