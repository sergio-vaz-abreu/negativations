package api

import (
	"github.com/negativations/modules/negativation/application"
	"github.com/negativations/modules/negativation/infrastructure"
	"github.com/pkg/errors"
)

func createNegativationController(config ArangoConfig) (*application.NegativationController, error) {
	client, err := infrastructure.NewClient(config.Host, config.Port, config.User, config.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create arango client")
	}
	database, err := infrastructure.CreateDatabase(client, "negativation")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create database")
	}
	repository, err := infrastructure.NewNegativationRepositoryArangoDB(database)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create negativation repository")
	}
	controller := application.NewNegativationController(repository)
	return controller, nil
}
