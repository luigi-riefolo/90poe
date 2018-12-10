package api

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	storermock "github.com/luigi-riefolo/nlp/svc-storer/pb/mocks"
)

// mocks
var (
	storerClientMock = &storermock.StorerClient{}
)

var (
	svc = &IngestorService{
		config:       testConf,
		storerClient: storerClientMock,
	}
	errTest = fmt.Errorf("test error")

	testConf = Config{
		DataFile:    "../data/data.csv",
		Environment: "test",
		Version:     "test",
	}
)

// mocks calls
var (
	storerClientStoreEntryHandler = storerClientMock.On("StoreEntryHandler",
		mock.AnythingOfType("*context.timerCtx"),
		mock.AnythingOfType("*pb.StoreEntryRequest"))
)

// TestCase represents the a test case description.
type TestCase struct {
	Name             string
	ExpectedResponse []byte
	ExpectedError    error
	ExpectedHTTPCode int32
	PreFn            func()
	AfterFn          func()
}

var (
	testCases = []TestCase{
		TestCase{
			Name:             "successful_ingestion",
			ExpectedHTTPCode: http.StatusOK,
		},
		TestCase{
			Name: "store_client_failure",
			PreFn: func() {
				storerClientStoreEntryHandler.Return(nil, errTest)
			},
			AfterFn: func() {
				storerClientStoreEntryHandler.Return(&empty.Empty{}, nil)
			},
		},
	}
)

func TestIngestFunction(t *testing.T) {

	storerClientStoreEntryHandler.Return(&empty.Empty{}, nil)

	for _, tc := range testCases {

		tc := tc

		t.Run(fmt.Sprintf("%s", tc.Name), func(t *testing.T) {

			if tc.PreFn != nil {
				tc.PreFn()
			}

			err := svc.ingest()

			if tc.ExpectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "unexpected error: %#v", err)
			}

			if tc.AfterFn != nil {
				tc.AfterFn()
			}
		})
	}
}
