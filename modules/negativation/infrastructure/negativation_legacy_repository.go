package infrastructure

import (
	"encoding/json"
	"fmt"
	"github.com/negativations/modules/negativation/domain"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

type NegativationLegacyRepository interface {
	GetAll() ([]*domain.Negativation, error)
}

var (
	ReadNegativationError = errors.New("failed to read negativation from legacy system")
)

func NewNegativationLegacyRepositoryAPI(baseUrl string) (*NegativationLegacyRepositoryAPI, error) {
	_, err := url.ParseRequestURI(baseUrl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse legacy api base url")
	}
	return &NegativationLegacyRepositoryAPI{baseUrl: baseUrl}, nil
}

type NegativationLegacyRepositoryAPI struct {
	baseUrl string
}

func (r *NegativationLegacyRepositoryAPI) GetAll() ([]*domain.Negativation, error) {
	httpResponse, err := http.Get(fmt.Sprintf("%s/%s", r.baseUrl, "negativation"))
	if err != nil {
		logrus.
			WithError(err).
			Error("failed to send request to legacy api")
		return nil, ReadNegativationError
	}
	body, err := ioutil.ReadAll(httpResponse.Body)
	if body != nil {
		defer httpResponse.Body.Close()
	}
	if err != nil {
		logrus.
			WithField("statusCode", httpResponse.StatusCode).
			WithError(err).
			Error("failed to read http response body")
		return nil, ReadNegativationError
	}
	if !IsSuccess(httpResponse) {
		logrus.
			WithField("statusCode", httpResponse.StatusCode).
			WithField("httpResponse", string(body)).
			Error("failed to have a success status code")
		return nil, ReadNegativationError
	}
	var negativations []*domain.Negativation
	err = json.Unmarshal(body, &negativations)
	if err != nil {
		logrus.
			WithField("statusCode", httpResponse.StatusCode).
			WithField("httpResponse", string(body)).
			WithError(err).
			Error("failed to decode negativations")
		return nil, ReadNegativationError
	}
	return negativations, nil
}

func IsSuccess(httpResponse *http.Response) bool {
	return httpResponse.StatusCode >= 200 && httpResponse.StatusCode < 300
}
