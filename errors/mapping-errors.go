package errors

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	log "github.com/presnalex/go-micro/v3/logger"
	"go.unistack.org/micro/v3/metadata"
)

const (
	internalErrorCode = "1"
	badRequestCode    = "2"
	notFoundErrorCode = "3"
)

func MapError(ctx context.Context, err error) (result pb.Error, status int) {
	status = http.StatusInternalServerError

	switch e := errors.Unwrap(err).(type) {
	case *BadRequestError:
		result = pb.Error{Code: badRequestCode, Title: "bad request error", Details: e.Error()}
		status = http.StatusBadRequest
	case *NotFoundError:
		result = pb.Error{Code: notFoundErrorCode, Title: "not found", Details: e.Error()}
		status = http.StatusOK
	default:
		result = pb.Error{Code: internalErrorCode, Title: "internal error"}
	}

	md, _ := metadata.FromIncomingContext(ctx)
	id, ok := md.Get(log.LoggerField)
	if !ok {
		uid, err := uuid.NewRandom()
		if err != nil {
			uid = uuid.Nil
		}
		id = uid.String()
	}
	result.Uuid = id

	return
}
