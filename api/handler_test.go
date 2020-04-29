package api

import (
	"encoding/json"
	"errors"
	"github.com/akhettar/app-features-manager/mocks"
	"github.com/akhettar/app-features-manager/model"
	"github.com/akhettar/app-features-manager/test"
	"github.com/golang/mock/gomock"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

// Publish App status with a valid jwt token
func TestPublishAppStatus_WithInvalidAppPlatformReturnBadRequestResponse(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		platform := "dummy"
		t.Logf("\tWhen Sending Publish App status request to endpoint with unsupported platform value:  \"%s\"", platform)
		{
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()
			version := "1.0"

			body := model.ReleaseRequest{Version: version, Platform: platform, Status: "deprecated"}
			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusBadRequest {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusBadRequest, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusBadRequest, test.BallotX, w.Code)
			}
		}
	}
}

// Publish App status with a valid jwt token
func TestPublishAppStatus_WithGivenStatusStatusSuccess(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Publish App status request to endpoint with valid token:  \"%s\"", "\\status")
		{
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()
			version := "1.0"
			platform := "ios"
			body := model.ReleaseRequest{Version: version, Platform: platform, Status: model.Deprecated}
			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusNoContent {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusNoContent, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusNoContent, test.BallotX, w.Code)
			}

			appStatus := queryAppStatus(version, platform, t, mockUnleash)
			if appStatus.Status == model.Deprecated {
				t.Logf("\t\tApp status should have \"%s\" status. %v", appStatus.Status, test.CheckMark)
			} else {
				t.Errorf("\t\tApp status should have \"%s\" status. %v", appStatus.Status, test.BallotX)
			}
		}
	}
}

// Publish App status with a valid jwt token
func TestPublishAppStatus_WithInvalidAppStatus(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		status := "dummystatus"
		t.Logf("\tWhen Sending Publish App status request to endpoint with unsupported status value:  \"%s\"", status)
		{
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()
			version := "1.0"
			platform := "ios"
			body := model.ReleaseRequest{Version: version, Platform: platform, Status: status}
			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusBadRequest {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusBadRequest, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusBadRequest, test.BallotX, w.Code)
			}
		}
	}
}

// Publish App status with an invalid jwt token
func TestPublishAppStatus_UnautorisedForInvalidJWTToken(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Publish App status request to endpoint with an invalid JWT token:  \"%s\"", "\\status")
		{
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			body := model.ReleaseRequest{Version: "1.0", Platform: "ios"}
			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.InvalidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusUnauthorized {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusUnauthorized, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusUnauthorized, test.BallotX, w.Code)
			}
		}
	}
}

// Publish App status with an invalid payload
func TestPublishAppStatus_WithInvalidRequest(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Publish App status request to endpoint with invalid payalod:  \"%s\"", "\\status")
		{
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			body := model.ReleaseRequest{Version: "1.0"}
			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusBadRequest {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusBadRequest, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusBadRequest, test.BallotX, w.Code)
			}
		}
	}
}

// Query status for a given version
func TestQueryAppStatus_Success(t *testing.T) {
	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query app status request with valid token:  \"%s\"", "\\version\\1.0\\ios")
		{
			mockUnleash := test.GetMockUnleashClient(t)
			// publish app status and assert
			version := "1.0"
			platform := "ios"
			publishAppStatus(version, platform, t, mockUnleash)

			// send query request

			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusOK {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusNoContent, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusNoContent, test.BallotX, w.Code)
			}

			var response model.ReleaseResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ReleaseResponse{Status: model.Supported}
			if reflect.DeepEqual(response, expectedResponse) {
				t.Logf("\t\tShould receive app status: \"%s\" . %v", expectedResponse.Status, test.CheckMark)
			} else {
				t.Logf("\t\tShould receive app status: \"%s\" . %v", expectedResponse.Status, test.BallotX)
			}
		}
	}
}

// Query status for a given version
func TestQueryAppStatus_WithPlatformShouldReturnBadRequestError(t *testing.T) {
	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query app status request with valid token:  \"%s\"", "\\version\\1.0\\ios")
		{
			// publish app status and assert
			mockUnleash := test.GetMockUnleashClient(t)
			version := "1.0"
			platform := "dummyPlatform"
			publishAppStatus(version, platform, t, mockUnleash)

			// send query request
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusBadRequest {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusBadRequest, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusBadRequest, test.BallotX, w.Code)
			}
		}
	}
}

func TestQueryAppStatus_ShouldReturnTheLatestStatus(t *testing.T) {
	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query app status request with valid token:  \"%s\"", "\\version\\1.0\\ios")
		{
			// publish app status and assert
			version := "1.0"
			platform := "ios"

			// First entry with default status `Supported`
			mockUnleash := test.GetMockUnleashClient(t)
			publishAppStatus(version, platform, t, mockUnleash)

			// Second entry with deprecated status
			time.Sleep(100 * time.Millisecond)
			body := model.ReleaseRequest{Version: version, Platform: platform, Status: model.Deprecated}
			publishAppStatusWithBody(version, platform, body, t, mockUnleash)

			// send query request
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusOK {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusNoContent, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusNoContent, test.BallotX, w.Code)
			}

			var response model.ReleaseResponse
			json.NewDecoder(w.Body).Decode(&response)
			expectedResponse := model.ReleaseResponse{Status: model.Deprecated}
			if reflect.DeepEqual(response, expectedResponse) {
				t.Logf("\t\tShould receive app status: \"%s\" . %v", expectedResponse.Status, test.CheckMark)
			} else {
				t.Logf("\t\tShould receive app status: \"%s\" . %v", expectedResponse.Status, test.BallotX)
			}
		}
	}
}

// Query status for a given version
func TestQueryAppStatus_NotFound(t *testing.T) {
	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query app status request with none registered app version: \"%s\"", "\\version\\1.0\\ios")
		{
			// publish app status and assert
			version := "1.0"
			platform := "ios"
			mockUnleash := test.GetMockUnleashClient(t)
			publishAppStatus(version, platform, t, mockUnleash)

			// send query request
			platform = "android"
			handler := NewAppStatusHandler(Repository, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			if w.Code == http.StatusNotFound {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusNotFound, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusNotFound, test.BallotX, w.Code)
			}
		}
	}
}

// Simulate insert failure
func TestAppHandler_PublishAppStatusWithFailureToPersistInDataStore(t *testing.T) {

	t.Logf("Given the app status service is up and running")
	{
		t.Logf("\tWhen Sending Publish App status request to endpoint:  \"%s\"", "\\tags")
		{
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Insert failed"
			body := model.ReleaseRequest{Version: "1.0", Platform: "ios", Status: "deprecated"}
			err := errors.New(expectedErrorMessage)

			mockRepo.EXPECT().Insert(gomock.Any()).Return(err).Times(1)

			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(mockRepo, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			if w.Code == http.StatusInternalServerError {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusInternalServerError, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusInternalServerError, test.BallotX, w.Code)
			}

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			if response.Message == expectedResponse.Message {
				t.Logf("\t\tThe body response should  contain a message \"%s\" . %v", expectedResponse.Message, test.CheckMark)
			} else {
				t.Errorf("\t\tThe body response should contain a message \"%s\". %v %v", response.Message, test.BallotX, expectedResponse.Message)
			}
		}
	}
}

func TestAppHandler_QueryAppStatusResultsInInternalServerError(t *testing.T) {

	t.Logf("Given the tag service is up and running")
	{
		t.Logf("\tWhen Sending Query App Status request to endpoint:  \"%s\"", "\\tags")
		{
			version := "1.0"
			platform := "ios"

			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockRepo := mocks.NewMockRepository(mockCtrl)

			expectedErrorMessage := "Failed to query datastore"
			err := errors.New(expectedErrorMessage)

			mockRepo.EXPECT().Find(gomock.Any(), gomock.Any()).Return(model.ReleaseResponse{}, err).Times(1)

			// send query request
			mockUnleash := test.GetMockUnleashClient(t)
			handler := NewAppStatusHandler(mockRepo, mockUnleash)
			router := handler.CreateRouter()

			req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// check call success
			test.Ok(err, t)

			// Assert response code status
			if w.Code == http.StatusInternalServerError {
				t.Logf("\t\tShould receive a \"%d\" status. %v", http.StatusInternalServerError, test.CheckMark)
			} else {
				t.Errorf("\t\tShould receive a \"%d\" status. %v %v", http.StatusInternalServerError, test.BallotX, w.Code)
			}

			var response model.ErrorResponse
			json.NewDecoder(w.Body).Decode(&response)

			expectedResponse := model.ErrorResponse{Code: http.StatusInternalServerError, Message: expectedErrorMessage}

			// check body response matches the expected response
			if response.Message == expectedResponse.Message {
				t.Logf("\t\tThe body response should  contain a message \"%s\" . %v", expectedResponse.Message, test.CheckMark)
			} else {
				t.Errorf("\t\tThe body response should contain a message \"%s\". %v %v", response.Message, test.BallotX, expectedResponse.Message)
			}
		}
	}

}

// Helper function
func publishAppStatus(version, platform string, t *testing.T, mockUnleash *mocks.MockUnleashService) {
	body := model.ReleaseRequest{Version: version, Platform: platform}
	publishAppStatusWithBody(version, platform, body, t, mockUnleash)
}

// Helper function
func publishAppStatusWithBody(version, platform string, body interface{}, t *testing.T, mockUnleash *mocks.MockUnleashService) {
	handler := NewAppStatusHandler(Repository, mockUnleash)
	router := handler.CreateRouter()
	req, err := test.HttpRequest(body, "/status", http.MethodPost, test.ValidToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// check call success
	test.Ok(err, t)
}

func queryAppStatus(version, platform string, t *testing.T, mockUnleash *mocks.MockUnleashService) model.ReleaseResponse {
	// send query request
	handler := NewAppStatusHandler(Repository, mockUnleash)
	router := handler.CreateRouter()

	req, err := test.HttpRequest(nil, "/status/version/"+version+"/"+platform, http.MethodGet, test.ValidToken)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// check call success
	test.Ok(err, t)

	var response model.ReleaseResponse
	json.NewDecoder(w.Body).Decode(&response)
	return response
}
