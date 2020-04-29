package repository

import (
	"errors"
	"github.com/akhettar/app-features-manager/model"
	"github.com/labstack/gommon/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"os"
)

const (

	// MongoURI environment variable
	MongoURI = "MONGO_URI"

	// DefaultMongoHost for running the app as a standalone server.
	DefaultMongoHost = "mongodb://localhost"

	// ENVIRONMENT OS variable
	ENVIRONMENT = "ENVIRONMENT"

	// DefaultEnvironment for running the app locally using docker-compose
	DefaultEnvironment = "default"

	// NotFoundErrorMessage the not found error message
	NotFoundErrorMessage = "not found"

	// DefaultDBName the default database name
	DefaultDBName = "app_status_api_dev"

	// DefaultCollection the default collection if not set in vault
	DefaultCollection = "status"

	// Collection the vault entry of the db collection
	Collection = "MONGO_COLLECTION"
)

// DBInfo the database info
type DBInfo struct {

	// URL the db url
	URL string

	// Database the database name
	Database string

	// Collection the collection name
	Collection string
}

// MongoRepository type
type MongoRepository struct {
	*mongo.Client
	DBInfo
}

// Repository interface
type Repository interface {
	Insert(body interface{}) error
	Find(version, platform string) (model.ReleaseResponse, error)
}

// NewRepository function to create an instance of Mongo repository
func NewRepository() Repository {

	// default value
	url := GetEnv(MongoURI, DefaultMongoHost)
	log.Infof("Initialising Document database session ...")
	clientOptions := options.Client().ApplyURI(url)

	dbname := DefaultDBName
	if clientOptions.Auth != nil && clientOptions.Auth.AuthSource != "" {
		dbname = clientOptions.Auth.AuthSource
	}
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal("Failed to initialise mongo vaultClient")
		panic(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	dbInfo := DBInfo{url, dbname, GetEnv(Collection, DefaultCollection)}
	log.Printf("Connected to Document DB %s, %s, %s", clientOptions.Hosts, dbInfo.Database, dbInfo.Collection)
	return &MongoRepository{client, dbInfo}
}

// Insert into data store
func (repo *MongoRepository) Insert(body interface{}) error {
	_, err := repo.Client.Database(repo.DBInfo.Database).Collection(repo.DBInfo.Collection).InsertOne(context.TODO(), &body)
	return err
}

// Find query the status of the app for given version shall return the latest
func (repo *MongoRepository) Find(version, platform string) (model.ReleaseResponse, error) {

	var results []*model.ReleaseDAO
	findOptions := options.Find()
	sortMap := make(map[string]interface{})
	sortMap["released"] = 1
	findOptions = findOptions.SetSort(sortMap)

	// query
	cursor, err := repo.Client.Database(repo.DBInfo.Database).Collection(repo.DBInfo.Collection).Find(context.TODO(),
		bson.M{model.AppVersion: version, model.AppPlatform: platform}, findOptions)

	if err != nil {
		return model.ReleaseResponse{}, err
	}

	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {
		var result model.ReleaseDAO
		err := cursor.Decode(&result)
		if err != nil {
			log.Error("Failed to decode document queried from the DB")
			return model.ReleaseResponse{}, err
		}
		results = append(results, &result)
	}

	// decode the response
	var result model.ReleaseDAO
	if len(results) == 0 {
		return model.ReleaseResponse{Status: result.Status}, errors.New(NotFoundErrorMessage)
	}
	return model.ReleaseResponse{Status: results[len(results)-1].Status}, err
}

// GetEnv env variable or fall back to default
func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
