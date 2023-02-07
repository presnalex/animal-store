package encoders

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	"net/http"
)

var ErrWrongResponseType = errors.New("JSONProto: wrong response message type")

type JSONProto struct {
	m jsonpb.Marshaler
}

func NewJSONProto() *JSONProto {
	return &JSONProto{m: jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     false,
	}}
}

func (e *JSONProto) Success(rw http.ResponseWriter, response interface{}) error {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	if v, ok := response.(proto.Message); ok {
		return errors.WithStack(e.m.Marshal(rw, v))
	}
	return ErrWrongResponseType
}

func (e *JSONProto) Error(rw http.ResponseWriter, err pb.Error, status int) error {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(status)
	return errors.WithStack(e.m.Marshal(rw, &pb.ErrorRsp{Error: &err}))
}
