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
	repository, err := infrastructure.NewNegativationRepositoryArangoDB(database)
	return repository, errors.Wrap(err, "failed to create negativation repository")
}

func createLegacyNegativationRepository(baseUrl string) (infrastructure.NegativationLegacyRepository, error) {
	repository, err := infrastructure.NewNegativationLegacyRepositoryAPI(baseUrl)
	return repository, errors.Wrap(err, "failed to create legacy negativation repository")
}

func createNegativationController(repository infrastructure.NegativationRepository) *application.NegativationController {
	return application.NewNegativationController(repository)
}

func createLegacyNegativationController(repository infrastructure.NegativationRepository, legacyRepository infrastructure.NegativationLegacyRepository) *application.LegacyNegativationController {
	return application.NewLegacyNegativationController(repository, legacyRepository)
}
