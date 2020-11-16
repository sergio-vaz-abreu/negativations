package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoadAPI(config ApplicationConfig) (*Api, error) {
	repository, err := createNegativationRepository(config.ArangoConfig)
	if err != nil {
		return nil, err
	}
	legacyRepository, err := createLegacyNegativationRepository(config.LegacyConfig)
	if err != nil {
		return nil, err
	}
	negativationController := createNegativationController(repository, config.SymmetricKeyConfig)
	legacyNegativationController := createLegacyNegativationController(repository, legacyRepository, config.SymmetricKeyConfig)
	httpServer := createHttpServer(config.ApiConfig, negativationController, legacyNegativationController)
	return &Api{httpServer: httpServer}, nil
}

type Api struct {
	httpServer *http.Server
}

func (api *Api) Run() <-chan error {
	out := make(chan error)
	go func() {
		if err := api.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			out <- errors.Wrap(err, "failed to listen and serve api")
		}
	}()
	return out
}

func (api *Api) Shutdown() {
	err := api.httpServer.Shutdown(nil)
	if err != nil {
		logrus.
			WithError(err).
			Error("failed shutting down api")
	}
}
