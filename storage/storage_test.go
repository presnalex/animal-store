package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test function GetAnimal, more detailed test that covered sql query
func TestGetAnimal(t *testing.T) {
	var err error
	//create result object
	animalResult := new(pb.AnimalRsp)
	//create all mocked object for testing function GetAnimal
	dbWrapper := NewDbWrapperMock()
	storage := NewStorage(dbWrapper)
	//Ok
	dbWrapper.On("SelectContext",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*[]*storage.AnimalResponse"),
		getAnimal,
		"1",
	).Run(func(args mock.Arguments) {
		argDbResp := args.Get(1).(*[]*AnimalResponse)
		*argDbResp = []*AnimalResponse{
			&AnimalResponse{
				AnimalId: sql.NullInt32{Int32: 1, Valid: true},
				Animal:   sql.NullString{String: "Zebra", Valid: true},
				Price:    sql.NullInt32{Int32: 5000, Valid: true},
			},
		}
	}).Return(nil)
	//NotFound
	dbWrapper.On("SelectContext",
		mock.AnythingOfType("*context.valueCtx"),
		mock.AnythingOfType("*[]*storage.AnimalResponse"),
		getAnimal,
		"-1",
	).Return(nil)

	//table with tests
	tests := []struct {
		name             string       //test case name
		testCaseFunction func() error //behaviour function that involves a test scenario
		expected         pb.AnimalRsp //expected result
		expectErr        bool         //error expectation flag
	}{
		{
			name: "Ok",
			testCaseFunction: func() error {
				req := new(pb.AnimalReq)
				req.AnimalId = "1"
				//test the function
				animalResult, err = storage.GetAnimal(context.Background(), req)
				if err != nil {
					return err
				}
				return nil
			},
			expected: pb.AnimalRsp{
				AnimalId: &wrappers.Int32Value{Value: 1},
				Animal:   &wrappers.StringValue{Value: "Zebra"},
				Price:    &wrappers.Int32Value{Value: 5000},
			},
		},
		{
			name: "Not found err",
			testCaseFunction: func() error {
				req := new(pb.AnimalReq)
				req.AnimalId = "-1"
				//test the function
				animalResult, err = storage.GetAnimal(context.Background(), req)
				if err != nil {
					return err
				}
				return nil
			},
			expectErr: true,
		},
	}
	//iterate our test cases and compare results
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.testCaseFunction()

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, "animal not found", err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, *animalResult)
			}
		})
	}
}
