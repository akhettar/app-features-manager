package features

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Unleash/unleash-client-go"
	"github.com/Unleash/unleash-client-go/context"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strings"
)

// ListOfFlags list of feature flag comma separated to be populated from the vault server
var ListOfFlags string

const (
	// UnleashUsername is the username to access the unleash server
	UnleashUsername = "UNLEASH_USERNAME"

	// UnleashPassword is the password to access the unleash server
	UnleashPassword = "UNLEASH_PASSWORD"

	// UnleashBaseURL is the base URL of the UnleashService server
	UnleashBaseURL = "UNLEASH_BASE_URL"

	// UnleashAppName the name of the app to use within unleash
	UnleashAppName = "AppFeatureManager"

	// UnleashAPISuffix the api location to add to the base URL
	UnleashAPISuffix = "/api"

	// AuthorizationHeader the name of the header to add
	AuthorizationHeader = "Authorization"

	// AuthorizationBasicPrefix the basic prefix to add to authorization requests
	AuthorizationBasicPrefix = "Basic "

	// CommaSeparator separator for the feature flag in vault
	CommaSeparator = ","

	// FeatureFlagList a list of all the feature flags currently defined in vault
	FeatureFlagList = "FEATURE_FLAG_LIST"

	// UnleashDevURL the unleash dev url.
	UnleashDevURL = "http://localhost"
)

// Flags holding list of feature flags
type Flags map[string]bool

// UnleashService interface
type UnleashService interface {
	FetchFeatureFlags(customerID string) Flags
}

// UnleashClient the unleash client
type UnleashClient struct{}

// NewUnleashClient initialises an instance of the Unleash client
func NewUnleashClient() UnleashService {

	// Getting default values that will be used in the integration test
	ListOfFlags = "ITFeature,ITFeatureDisabled"
	password := fetchEnv(UnleashPassword)
	username := fetchEnv(UnleashUsername)
	url := UnleashDevURL

	headers := http.Header{
		AuthorizationHeader: []string{AuthorizationBasicPrefix +
			base64.StdEncoding.EncodeToString([]byte(username+":"+password))},
	}

	unleash.Initialize(
		unleash.WithListener(&unleash.DebugListener{}),
		unleash.WithAppName(UnleashAppName),
		unleash.WithUrl(url),
		unleash.WithCustomHeaders(headers),
	)
	return UnleashClient{}
}

// FetchFeatureFlags fetches feature flags for a given customerID
func (cl UnleashClient) FetchFeatureFlags(customerID string) Flags {

	ctx := context.Context{
		UserId: customerID,
	}
	results := make(map[string]bool)

	flagNames := strings.Split(ListOfFlags, CommaSeparator)
	ch := make(chan func() (string, bool), len(flagNames))
	defer close(ch)

	// Fire all the queries for given flags in the background
	for _, flag := range flagNames {
		go fetchFlagStatus(ctx, ch, flag)
	}

	// Read all the feature flag status
	for range flagNames {
		flag, status := (<-ch)()
		results[flag] = status
	}
	return results
}

// Fetches the status for a given flag from the unleash server
func fetchFlagStatus(ctx context.Context, c chan func() (string, bool), flag string) {
	c <- func() (string, bool) {
		log.Infof("Fetching feature flag value for %s", flag)
		status := unleash.IsEnabled(flag, unleash.WithContext(ctx))
		log.Infof("Found feature flag value for %s:%v", flag, status)
		return flag, status
	}
}

func fetchEnv(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	panic(errors.New(fmt.Sprintf("Required environment variable nout found %s", key)))
}
