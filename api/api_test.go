package api

import (
	"fmt"
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/infrastructure"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestApi(t *testing.T) {
	g := NewGomegaWithT(t)
	config := MakeConfig()
	client, err := infrastructure.NewClient(config.ArangoConfig.Host, config.ArangoConfig.Port, config.ArangoConfig.User, config.ArangoConfig.Password)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := infrastructure.CreateDatabase(client, "negativation")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collection, err := infrastructure.CreateCollection(database, "negativations")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	sut, err := LoadAPI(config)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	appErr := sut.Run()
	g.Consistently(appErr).Should(
		Not(Receive()))

	t.Run("Get negativations", func(t *testing.T) {
		err = infrastructure.CleanCollection(t, collection)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		negativation, err := domain.NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = negativation.Encrypt(config.SymmetricKeyConfig.SymmetricKey, config.SymmetricKeyConfig.EncryptionContext)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = infrastructure.InsertNegativations(t, collection, negativation)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g := NewGomegaWithT(t)

		httpResponse, err := http.Get(fmt.Sprintf("http://localhost:%d/negativation?cpf=515.374.764-67", config.ApiConfig.Port))

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(httpResponse).Should(
			HaveHTTPStatus(http.StatusOK))
		g.Expect(ReadHttpResponseBody(httpResponse.Body)).Should(
			MatchJSON(`
{
  "status":"success",
  "data":[
    {
      "companyDocument":"59291534000167",
	  "companyName":"ABC S.A.",
	  "customerDocument":"51537476467",
	  "value":1235.23,
	  "contract":"bc063153-fb9e-4334-9a6c-0d069a42065b",
	  "debtDate":"2015-11-13T23:32:51Z",
	  "inclusionDate":"2020-11-13T23:32:51Z"
	}
  ]
}
`))
		g.Consistently(appErr).Should(
			Not(Receive()))
	})

	t.Run("Synchronize negativations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		err := infrastructure.CleanCollection(t, collection)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		httpResponse, err := http.Post(fmt.Sprintf("http://localhost:%d/negativation/synchronize", config.ApiConfig.Port), "application/json", nil)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(httpResponse).Should(
			HaveHTTPStatus(http.StatusOK))
		g.Expect(infrastructure.GetNegativations(t, collection)).Should(
			And(
				HaveLen(5),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("77723018000146"),
					"CompanyName":      BeEquivalentTo("123 S.A."),
					"CustomerDocument": BeEquivalentTo("FUNoBV5JaZEji6c="),
					"Value":            BeEquivalentTo(400.00),
					"Contract":         BeEquivalentTo("5f206825-3cfe-412f-8302-cc1b24a179b0"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 10, 12, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 10, 12, 23, 32, 51, 0, time.UTC)),
				})))))
		g.Consistently(appErr).Should(
			Not(Receive()))
	})
}

func MakeConfig() ApplicationConfig {
	return ApplicationConfig{
		SymmetricKeyConfig: SymmetricKeyConfig{
			SymmetricKey:      "UkVDMgAAAC13PCVZAKOczZXUpvkhsC+xvwWnv3CLmlG0Wzy8ZBMnT+2yx/dg",
			EncryptionContext: "context",
		},
		ApiConfig: ApiConfig{
			Port: 8090,
		},
		ArangoConfig: ArangoConfig{
			Host:     "localhost",
			Port:     8529,
			User:     "root",
			Password: "somepassword",
		},
		LegacyConfig: LegacyConfig{
			Url: "http://localhost:3000",
		},
	}

}

func ReadHttpResponseBody(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}
