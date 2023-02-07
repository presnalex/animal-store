package handler

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	pb "github.com/presnalex/animal-store-proto/go/animal_store_proto"
	errs "github.com/presnalex/animal-store/errors"
	"github.com/presnalex/animal-store/handler/encoders"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAnimal(t *testing.T) {
	var err error
	//declare results
	var resultCode int
	var resultBody string
	//create mock objects
	writer, _ := NewWriter(encoders.NewJSONProto())

	strg := NewStorageMock()
	//Ok
	strg.On("GetAnimal",
		mock.AnythingOfType("*context.emptyCtx"),
		&pb.AnimalReq{AnimalId: "1"},
	).Return(&pb.AnimalRsp{
		AnimalId: &wrappers.Int32Value{Value: 1},
		Animal:   &wrappers.StringValue{Value: "Zebra"},
		Price:    &wrappers.Int32Value{Value: 5000},
	}, nil)
	//NotFound
	strg.On("GetAnimal",
		mock.AnythingOfType("*context.emptyCtx"),
		&pb.AnimalReq{AnimalId: "-1"},
	).Return(nil, errs.NewNotFoundError("animal not found"))

	hdl, err := NewHandler(writer, strg)
	if err != nil {
		t.Fatal(err)
	}

	//table with tests
	tests := []struct {
		name             string       //test case name
		testCaseFunction func() error //behaviour function that involves a test scenario
		expectCode       int          //expected result
		expectBody       string       //expected result
		expectErr        bool         //error expectation flag
	}{
		{
			name: "Ok",
			testCaseFunction: func() error {
				// Request
				req, err := http.NewRequest("GET", "/animal-store/v1/getAnimal", nil)
				if err != nil {
					return err
				}
				// Request parameters
				q := req.URL.Query()
				q.Add("animalId", "1")
				// assign encoded query string to http request
				req.URL.RawQuery = q.Encode()
				// create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(hdl.GetAnimal)
				// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
				handler.ServeHTTP(rr, req)
				resultCode = rr.Code
				resultBody = rr.Body.String()
				return err
			},
			expectCode: http.StatusOK,
			expectBody: `{"animalId":1,"animal":"Zebra","price":5000}`,
		},
		{
			name: "Notfound",
			testCaseFunction: func() error {
				req, err := http.NewRequest("GET", "/animal-store/v1/getAnimal", nil)
				if err != nil {
					return err
				}
				q := req.URL.Query()
				q.Add("animalId", "-1")
				req.URL.RawQuery = q.Encode()
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(hdl.GetAnimal)
				handler.ServeHTTP(rr, req)
				resultCode = rr.Code
				resultBody = rr.Body.String()
				return err
			},
			expectCode: http.StatusOK,
			expectBody: `((.|\n)*)"code":"3"((.|\n)*)"details":"animal not found"((.|\n)*)`,
		},
		{
			name: "BadRequest",
			testCaseFunction: func() error {
				req, err := http.NewRequest("GET", "/animal-store/v1/getAnimal", nil)
				if err != nil {
					return err
				}
				rr := httptest.NewRecorder()
				handler := http.HandlerFunc(hdl.GetAnimal)
				handler.ServeHTTP(rr, req)
				resultCode = rr.Code
				resultBody = rr.Body.String()
				return err
			},
			expectCode: http.StatusBadRequest,
			expectBody: `((.|\n)*)"code":"2"((.|\n)*)"title":"bad request error"((.|\n)*)`,
		},
	}
	//iterate our test cases and compare results
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = tt.testCaseFunction()

			assert.NoError(t, err)
			assert.Equal(t, tt.expectCode, resultCode)
			assert.Regexp(t, regexp.MustCompile(tt.expectBody), resultBody)
		})
	}
}
