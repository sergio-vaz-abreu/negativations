package infrastructure

import (
	"github.com/arangodb/go-driver"
	"github.com/negativations/modules/negativation/domain"
	"testing"
)

func CleanCollection(_ testing.TB, collection driver.Collection) error {
	return collection.Truncate(nil)
}

func InsertNegativations(_ testing.TB, collection driver.Collection, negativations ...*domain.Negativation) error {
	_, _, err := collection.CreateDocuments(driver.WithWaitForSync(nil, true), negativations)
	return err
}
