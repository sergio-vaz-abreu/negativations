package infrastructure

import (
	"github.com/arangodb/go-driver"
	"github.com/negativations/modules/negativation/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NegativationRepository interface {
	GetByCPF(cpf domain.CPF) ([]*domain.Negativation, error)
}

var (
	ReadNegativationsByCPFErr = errors.New("failed to get negativations by cpf")
)

func NewNegativationRepositoryArangoDB(database driver.Database) (*NegativationRepositoryArangoDB, error) {
	collection, err := CreateCollection(database, "negativations")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create negativations' collection")
	}
	return &NegativationRepositoryArangoDB{collection: collection}, nil
}

type NegativationRepositoryArangoDB struct {
	collection driver.Collection
}

func (repo *NegativationRepositoryArangoDB) GetByCPF(cpf domain.CPF) ([]*domain.Negativation, error) {
	query :=
		`FOR negativation in @@collection
FILTER negativation.customerDocument == @customerDocument
RETURN negativation
`
	cursor, err := repo.collection.Database().Query(nil, query, map[string]interface{}{
		"@collection":      repo.collection.Name(),
		"customerDocument": cpf,
	})
	if err != nil {
		logrus.
			WithField("cpf", cpf).
			WithError(err).
			Error(ReadNegativationsByCPFErr)
		return nil, ReadNegativationsByCPFErr
	}
	var negativations []*domain.Negativation
	for cursor.HasMore() {
		var negativation domain.Negativation
		_, err := cursor.ReadDocument(nil, &negativation)
		if err != nil {
			logrus.
				WithField("cpf", cpf).
				WithError(err).
				Error(ReadNegativationsByCPFErr)
			return nil, ReadNegativationsByCPFErr
		}
		negativations = append(negativations, &negativation)
	}
	return negativations, nil
}
