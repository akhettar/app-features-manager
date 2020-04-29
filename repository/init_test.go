package repository

import (
	"context"
	"flag"
	"fmt"
	"github.com/akhettar/docker-db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"

	"net/http/httptest"
	"os"
	"testing"
)

const ProfileEnvVar = "PROFILE"

var (
	HttpServer          *httptest.Server
	RepositoryUnderTest *MongoRepository
)

// TestMain wraps all tests with the needed initialized mock DB and fixtures
// This test runs before other integration test. It starts an instance of mongo db in the background (provided you have mongo
// installed on the server on which this test will be running) and shuts it down.
func TestMain(m *testing.M) {

	flag.Parse()

	c := dbtest.StartMongoContainer()
	log.Printf("running mongo %s:%d", c.Host(), c.Port())

	uri := fmt.Sprintf("mongodb://%s:%d", c.Host(), c.Port())
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	RepositoryUnderTest = &MongoRepository{client, DBInfo{uri, DefaultDBName, DefaultCollection}}

	// Run the test suite
	retCode := m.Run()

	c.Destroy()

	// call with result of m.Run()
	os.Exit(retCode)
}
