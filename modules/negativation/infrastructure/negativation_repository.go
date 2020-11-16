package infrastructure

import (
	"github.com/arangodb/go-driver"
	"github.com/negativations/modules/negativation/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type NegativationRepository interface {
	GetByCPF(cpf domain.CPF) ([]*domain.Negativation, error)
	Synchronize(negativations ...*domain.Negativation) error
}

var (
	ReadNegativationsByCPFErr = errors.New("failed to get negativations by cpf")
	RemoveOldNegativationsErr = errors.New("failed to remove old negativations")
	SaveNegativationsErr      = errors.New("failed to save negativation")
)

func NewNegativationRepositoryArangoDB(database driver.Database, collectionsName string) (*NegativationRepositoryArangoDB, error) {
	collection, err := CreateCollection(database, collectionsName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create negativations' collection")
	}
	_, _, err = collection.EnsurePersistentIndex(nil, []string{"customerDocument"}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create customerDocument index")
	}
	return &NegativationRepositoryArangoDB{collection: collection}, nil
}

type NegativationRepositoryArangoDB struct {
	collection driver.Collection
}

func (repo *NegativationRepositoryArangoDB) GetByCPF(cpf domain.CPF) ([]*domain.Negativation, error) {
	query :=
		`
FOR negativation in @@collection
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

func (repo *NegativationRepositoryArangoDB) Synchronize(negativations ...*domain.Negativation) error {
	contracts := repo.getContractList(negativations)
	err := repo.removeOldNegativations(contracts)
	if err != nil {
		return err
	}
	for _, negativation := range negativations {
		err = repo.upsertNegativation(negativation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *NegativationRepositoryArangoDB) removeOldNegativations(contracts []string) error {
	query :=
		`
FOR negativation IN @@collection
FILTER negativation.contract not in @contracts
REMOVE {_key: negativation._key} in @@collection OPTIONS {waitForSync: true}
`
	_, err := repo.collection.Database().Query(nil, query, map[string]interface{}{
		"@collection": repo.collection.Name(),
		"contracts":   contracts,
	})
	if err != nil {
		logrus.
			WithError(err).
			Error(RemoveOldNegativationsErr)
		return RemoveOldNegativationsErr
	}
	return nil
}

func (repo *NegativationRepositoryArangoDB) getContractList(negativations []*domain.Negativation) []string {
	var contracts []string
	for _, negativation := range negativations {
		contracts = append(contracts, negativation.Contract)
	}
	return contracts
}

func (repo *NegativationRepositoryArangoDB) upsertNegativation(negativation *domain.Negativation) error {
	query :=
		`
UPSERT {contract: @contract}
INSERT @negativation
UPDATE @negativation in @@collection OPTIONS {waitForSync: true}
`
	_, err := repo.collection.Database().Query(nil, query, map[string]interface{}{
		"@collection":  repo.collection.Name(),
		"negativation": negativation,
		"contract":     negativation.Contract,
	})
	if err != nil {
		logrus.
			WithError(err).
			Error(SaveNegativationsErr)
		return SaveNegativationsErr
	}
	return nil
}
