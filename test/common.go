package test

import (
	"bytes"
	"encoding/json"
	"github.com/akhettar/app-features-manager/mocks"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"testing"
)

const (

	// ValidToken used in integration test.
	ValidToken = "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJ3c28yLm9yZ1wvcHJvZHVjdHNcL2FtIiwiZXhwIjozMTI1NzI1MTUyLCJodHRwOlwvXC93c28yLm9yZ1wvY2xhaW1zXC91c2VydHlwZSI6IkFQUExJQ0FUSU9OIn0.eV7XNqQl361LBohJnK0rQ-HorvOFNf-ILujkVAT9yZ8"
	//ValidToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJpc3MiOiJ0ZXN0In0.TWqKCHLFEq2wxGyLQksr4WBg-YQ6T9-fM9XQsgHs-W8:"

	//InvalidToken
	InvalidToken = "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.SXsvSJVO80YrQi_HrIksQcVI5BXe0AhRxHYd7b2_dqq_NdTvF1bQCNt0IAeeOfS2"

	// CheckMark used for unit test highlight.
	CheckMark = "\u2713"

	// BallotX used for unit test highlight.
	BallotX = "\u2717"
)

// HttpRequest helper
func HttpRequest(jsonReq interface{}, endpoint string, method string, token string) (*http.Request, error) {
	req, err := http.NewRequest(method, endpoint, RequestBody(jsonReq))
	if err != nil {
		panic("Failed to marshall json request")
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)
	return req, err
}

// RequestBody helper
func RequestBody(req interface{}) *bytes.Buffer {
	jsonBytes, err := json.Marshal(req)
	if err != nil {
		panic("Failed to marshall json request")
	}
	return bytes.NewBuffer(jsonBytes)
}

// Ok assert helper
func Ok(err error, t *testing.T) {
	if err != nil {
		t.Fatal("\t\tShould be able to make the Post call.", BallotX, err)
	}
}

// GetMockUnleashClient initialises an instance of the unleash mock client
func GetMockUnleashClient(t *testing.T) *mocks.MockUnleashService {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockUnleash := mocks.NewMockUnleashService(mockCtrl)
	results := make(map[string]bool)
	results["BANK_AGGREGATION"] = false
	results["MIF"] = false
	results["MIF_LIMITED_COMPANY"] = false
	mockUnleash.EXPECT().FetchFeatureFlags(gomock.Any()).Return(results).AnyTimes()
	return mockUnleash
}
