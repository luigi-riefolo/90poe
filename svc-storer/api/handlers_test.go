package api

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/luigi-riefolo/90poe/svc-storer/pb"
	"github.com/stretchr/testify/assert"
)

var (
	data = []*pb.Entry{
		&pb.Entry{
			Id:           "1",
			Name:         "Kirk",
			Email:        "ornare@sedtortor.net",
			MobileNumber: "(013890) 37420",
		},
		&pb.Entry{
			Id:           "2",
			Name:         "Cain",
			Email:        "volutpat@semmollisdui.com",
			MobileNumber: "(016977) 2245",
		},
	}
)

var (
	ctx = context.TODO()
	svc = &StorerService{
		config:      testConf,
		userEntries: map[string]*pb.Entry{},
	}
	errTest = fmt.Errorf("test error")

	testConf = Config{
		Environment: "test",
		Version:     "test",
	}
)

// TestCase represents the a test case description.
type TestCase struct {
	Name             string
	Req              *pb.StoreEntryRequest
	ExpectedResponse []byte
	ExpectedError    error
	ExpectedHTTPCode int32
	PreFn            func()
	AfterFn          func()
}

var (
	testCases = []TestCase{
		TestCase{
			Name: "successful_store",
			Req: &pb.StoreEntryRequest{
				Entry: data[0],
			},
			ExpectedHTTPCode: http.StatusOK,
		},
	}
)

func TestStoreEntryHandler(t *testing.T) {

	for _, tc := range testCases {

		tc := tc

		t.Run(fmt.Sprintf("%s", tc.Name), func(t *testing.T) {

			if tc.PreFn != nil {
				tc.PreFn()
			}

			_, err := svc.StoreEntryHandler(ctx, tc.Req)

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
