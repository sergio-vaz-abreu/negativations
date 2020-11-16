package api

import (
	"github.com/negativations/modules/negativation/application"
	"github.com/negativations/modules/negativation/infrastructure"
	"github.com/pkg/errors"
)

func createNegativationRepository(config ArangoConfig) (infrastructure.NegativationRepository, error) {
	client, err := infrastructure.NewClient(config.Host, config.Port, config.User, config.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create arango client")
	}
	database, err := infrastructure.CreateDatabase(client, "negativation")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}
	repository, err := infrastructure.NewNegativationRepositoryArangoDB(database, "negativations")
	return repository, errors.Wrap(err, "failed to create negativation repository")
}

func createLegacyNegativationRepository(config LegacyConfig) (infrastructure.NegativationLegacyRepository, error) {
	repository, err := infrastructure.NewNegativationLegacyRepositoryAPI(config.Url)
	return repository, errors.Wrap(err, "failed to create legacy negativation repository")
}

func createNegativationController(repository infrastructure.NegativationRepository, config SymmetricKeyConfig) *application.NegativationController {
	return application.NewNegativationController(repository, config.SymmetricKey, config.EncryptionContext)
}

func createLegacyNegativationController(repository infrastructure.NegativationRepository, legacyRepository infrastructure.NegativationLegacyRepository, config SymmetricKeyConfig) *application.LegacyNegativationController {
	return application.NewLegacyNegativationController(repository, legacyRepository, config.SymmetricKey, config.EncryptionContext)
}
