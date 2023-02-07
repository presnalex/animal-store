package storage

import (
	"context"

	"github.com/google/uuid"

	"github.com/presnalex/go-micro/v3/wrapper/requestid"
	"github.com/stretchr/testify/mock"
	"go.unistack.org/micro/v3/metadata"
)

type DbWrapperMock struct {
	mock.Mock
}

func (m *DbWrapperMock) SelectContext(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	argsCalles := []interface{}{ctx, dst, query}
	argsCalles = append(argsCalles, args...)
	arguments := m.Called(argsCalles...)
	return arguments.Error(0)
}

func NewDbWrapperMock() *DbWrapperMock {
	return &DbWrapperMock{}
}

func NewContext() context.Context {
	ctx := context.Background()
	uid, err := uuid.NewRandom()
	if err != nil {
		uid = uuid.Nil
	}
	id := uid.String()
	md := make(metadata.Metadata)
	ctx = metadata.NewIncomingContext(ctx, md)
	ctx = requestid.SetIncomingRequestId(ctx, id)
	return ctx
}
