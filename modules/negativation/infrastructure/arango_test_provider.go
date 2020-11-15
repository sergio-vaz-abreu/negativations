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

func GetNegativations(_ testing.TB, collection driver.Collection) ([]*domain.Negativation, error) {
	query :=
		`
FOR n in negativations
RETURN n
`
	cursor, err := collection.Database().Query(nil, query, nil)
	if err != nil {
		return nil, err
	}
	var negativations []*domain.Negativation
	for cursor.HasMore() {
		var negativation domain.Negativation
		_, err := cursor.ReadDocument(nil, &negativation)
		if err != nil {
			return nil, err
		}
		negativations = append(negativations, &negativation)
	}
	return negativations, nil
}
