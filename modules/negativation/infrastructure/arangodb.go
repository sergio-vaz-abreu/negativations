package infrastructure

import (
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func NewClient(host string, port int, user, password string) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{fmt.Sprintf("tcp://%s:%d", host, port)},
	})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(user, password),
	})
	return client, err
}

func CreateDatabase(client driver.Client, dbName string) (driver.Database, error) {
	exists, err := client.DatabaseExists(nil, dbName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return client.CreateDatabase(driver.WithWaitForSync(nil, true), dbName, nil)
	}
	return client.Database(nil, dbName)
}

func CreateCollection(database driver.Database, collectionName string) (driver.Collection, error) {
	exists, err := database.CollectionExists(nil, collectionName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return database.CreateCollection(nil, collectionName, nil)
	}
	return database.Collection(nil, collectionName)
}
