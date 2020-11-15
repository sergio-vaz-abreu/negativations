package api

import (
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/infrastructure"
	. "github.com/onsi/gomega"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestApi(t *testing.T) {
	g := NewGomegaWithT(t)
	client, err := infrastructure.NewClient("localhost", 8529, "root", "somepassword")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := infrastructure.CreateDatabase(client, "negativation")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collection, err := infrastructure.CreateCollection(database, "negativations")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	sut, err := LoadAPI(8090, "http://localhost", ArangoConfig{
		Host:     "localhost",
		Port:     8529,
		User:     "root",
		Password: "somepassword",
	})
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
		err = infrastructure.InsertNegativations(t, collection, negativation)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		g := NewGomegaWithT(t)

		httpResponse, err := http.Get("http://localhost:8090/negativation?cpf=515.374.764-67")

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

		httpResponse, err := http.Post("http://localhost:8090/negativation/synchronize", "application/json", nil)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(httpResponse).Should(
			HaveHTTPStatus(http.StatusOK))
		g.Consistently(appErr).Should(
			Not(Receive()))
	})
}

func ReadHttpResponseBody(reader io.Reader) ([]byte, error) {
	return ioutil.ReadAll(reader)
}
