package repository

import (
	"fmt"
	"github.com/akhettar/app-features-manager/model"
	"github.com/akhettar/app-features-manager/test"
	"os"
	"testing"
	"time"
)

// TestMongoRepository_Insert should test successful insert
func TestMongoRepository_Insert(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Publish app status request to endpoint:  \"%s\"", "\\status")
		{
			status := model.ReleaseDAO{Status: "supported", Version: "1.0", Platform: "ios", Released: time.Now()}
			err := RepositoryUnderTest.Insert(status)
			if err == nil {
				t.Logf("\t\tThe insert should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe insert should have been successful %v", test.BallotX)
			}
		}
	}
}

// TestMongoRepository_Insert should test successful insert
func TestMongoRepository_FindSuccess(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query status of the app for given version: \"%s\"", "1.0")
		{
			// publish status
			publishStatus("1.0", "ios")

			// query
			result, err := RepositoryUnderTest.Find("1.0", "ios")
			if err == nil {
				t.Logf("\t\tThe query should have been successful %v", test.CheckMark)
			} else {
				t.Errorf("\t\tThe query should have been successful %v", test.BallotX)
			}

			if result.Status == "supported" {
				t.Logf("\t\tThe status should have been set to %v %v", result.Status, test.CheckMark)
			} else {
				t.Errorf("\t\tThe status should have been set to %v %v", result.Status, test.BallotX)
			}
		}
	}
}

// TestMongoRepository_Insert should test successful insert
func TestMongoRepository_AppStatusNotFound(t *testing.T) {

	t.Logf("Given the app status api is up and running")
	{
		t.Logf("\tWhen Sending Query status of the app for given version: \"%s\"", "1.0")
		{
			// query
			_, err := RepositoryUnderTest.Find("2.0", "ios")
			if err != nil && err.Error() == NotFoundErrorMessage {
				t.Logf("\t\tThe query should have failed with status %v %v", err.Error(), test.CheckMark)
			} else {
				t.Errorf("\t\tThe query should have failed with status %v %v", err.Error(), test.BallotX)
			}
		}
	}
}

func TestNewRepository(t *testing.T) {
	os.Setenv(ENVIRONMENT, "dev")
	go NewRepository()
}

// Helper function
func publishStatus(version, platform string) {
	status := model.ReleaseDAO{Status: "supported", Version: version, Platform: platform, Released: time.Now()}
	err := RepositoryUnderTest.Insert(status)
	if err != nil {
		fmt.Printf("\t\tThe insert should have been successful %v", test.CheckMark)
	} else {
		fmt.Errorf("\t\tThe insert should have been successful %v", test.CheckMark)
	}
}
