package api

import (
	"encoding/base64"
	"errors"
	"github.com/akhettar/app-features-manager/features"
	"github.com/akhettar/app-features-manager/model"
	"github.com/akhettar/app-features-manager/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	// JwtIssuer the jwt issuer. This will get set per environment at runtime
	JwtIssuer = "default-issuer"

	// JwtSecret the JWT secret. This will get set per environment at runtime
	JwtSecret = "SNtvSaiXEnP3eG5Pm7ORCMiphN/pOGCpq6yxo1Sx1b0="

	// DB name
	DatabaseName = "app-feature-release-db"

	// CustomerID header
	CustomerID = "CUSTOMER_ID"
)

const (
	// ServiceName the service name used by the vault config.
	ServiceName = "app-status-api"

	// JSONMimeType the mime type
	JSONMimeType = "application/json; charset=utf-8"

	// ContentType the content type
	ContentType = "Content-Type"

	// DataDogAgentHostEnv the datadog agent host environment
	DataDogAgentHostEnv = "DD_AGENT_HOST"

	// DataDogAgentHostFallback the datadog agent hostname
	DataDogAgentHostFallback = "localhost"

	// DataDogServiceNameEnv the datadog service name
	DataDogServiceNameEnv = "DD_SERVICE_NAME"

	// DataDogTracerPort the datadog tracer port
	DataDogTracerPort = "8126"

	// JwtIssuerKey the key of the jwt issuer by which is stored in vault.
	JwtIssuerKey = "JWT_ISSUER"

	// JwtSecretKey the JWT secret key by which the value is stored in vault.
	JwtSecretKey = "JWT_SECRET"

	// IdentityKey the key id defined in the JWT token
	IdentityKey = "id"
)

// AppVersionHandler the app status handler
type AppVersionHandler struct {
	repository.Repository
	features.UnleashService
}

// NewAppStatusHandler creates an instance of AppVersionHandler
func NewAppStatusHandler(repo repository.Repository, unleashClient features.UnleashService) *AppVersionHandler {
	return &AppVersionHandler{repo, unleashClient}
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// CreateRouter creates an instance of Gin Router
func (handler *AppVersionHandler) CreateRouter() *echo.Echo {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// register validator
	validate := validator.New()
	validate.RegisterStructValidation(model.ReleaseRequestStructLevelValidation, model.ReleaseRequest{})
	e.Validator = &model.ReleaseRequestValidator{validate}

	// JWT Auth middleware
	key, err := base64.StdEncoding.DecodeString(fetchValue(JwtSecretKey, JwtSecret))
	if err != nil {
		log.Fatal(errors.New("failed to decode jwt secret key"))
	}
	middlewareFunc := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(key),
	})

	// Define the routes
	e.GET("/status/version/:version/:platform", handler.GetAppFeatures, middlewareFunc)
	e.POST("/status", handler.PublishAppStatus, middlewareFunc)
	e.GET("/health", handler.Health)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	return e
}

// @Summary Get App Status
// @ID get-app-status
// @Description Query app status for a given app release version
// @Accept  json
// @Produce  json
// @Param version path string true "app version"
// @Param platform path string true "App platform IOS, Android"
// @Success 200 {object} model.ReleaseResponse	"ok"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ReleaseResponse	"not found"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /status/version/{version}/{platform} [get]
func (handler *AppVersionHandler) GetAppFeatures(c echo.Context) error {
	version := c.Param(model.AppVersion)
	platform, err := model.Platform(c.Param(model.AppPlatform)).Value()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()})
	}
	log.Printf("Received request to retrieve app status for given version \"%v\" and platform  \"%v", version, platform)
	httpResponse := http.StatusOK
	result, err := handler.Find(version, platform)
	if err != nil {
		if err.Error() == repository.NotFoundErrorMessage {
			httpResponse = http.StatusNotFound
			result = model.ReleaseResponse{
				Status: model.Supported,
			}
		} else {
			return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()})
		}
	}

	// Fetch all the features
	result.Flags = handler.FetchFeatureFlags(c.Request().Header.Get(CustomerID))

	// return response to the client
	return c.JSON(httpResponse, result)
}

// @Summary Publish App status
// @ID publish-app-status
// @Description Publish a new app status
// @Accept  json
// @Produce  json
// @Param status-request body model.ReleaseRequest true "New App Status"
// @Success 204
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 500 {object} model.ErrorResponse "Internal server error"
// @Router /status/ [post]
func (handler *AppVersionHandler) PublishAppStatus(c echo.Context) error {

	// unmarshal the request
	request := new(model.ReleaseRequest)

	if err := c.Bind(&request); err != nil {
		log.Error(err.Error())
		return errorResponse("Failed to parse json request", http.StatusBadRequest, c)
	}

	if err := c.Validate(request); err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest, c)
	}
	log.Infof("Received request to publish app status for given version \"%v\" and platform  \"%v", request.Version, request.Platform)

	// persist the release
	err := handler.Insert(&model.ReleaseDAO{Version: request.Version, Platform: request.Platform, Released: time.Now(), Status: request.Status})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.ErrorResponse{Code: http.StatusInternalServerError, Message: "Insert failed"})
	}
	return c.NoContent(http.StatusNoContent)
}

// @Summary Health
// @ID health
// @Description Query the health of the service
// @Success 200 {object} model.EmptyBody "ok"
// @Failure 500 {object} model.EmptyBody "Server is down"
// @Router /health [get]
func (handler *AppVersionHandler) Health(c echo.Context) error {
	return c.String(200, "Success")
}

// Helper method
func errorResponse(msg string, status int, c echo.Context) error {
	return c.JSON(status, model.ErrorResponse{Message: msg, Code: status})
}

// Fetches environment variable, returns default if not set
func fetchValue(key, def string) string {
	if s, ok := os.LookupEnv(key); ok {
		return s
	}
	return def
}
