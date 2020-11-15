package application

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/infrastructure"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	symmetricKey      = "UkVDMgAAAC13PCVZAKOczZXUpvkhsC+xvwWnv3CLmlG0Wzy8ZBMnT+2yx/dg"
	encryptionContext = "context"
)

func TestLegacyNegativationController(t *testing.T) {
	g := NewGomegaWithT(t)
	handler := gin.Default()
	handler.GET("negativation", func(context *gin.Context) {
		data := `
[
  {
    "companyDocument": "59291534000167",
    "companyName": "ABC S.A.",
    "customerDocument": "51537476467",
    "value": 2500.23,
    "contract": "bc063153-fb9e-4334-9a6c-0d069a42065b",
    "debtDate": "2015-11-13T20:32:51-03:00",
    "inclusionDate": "2020-11-13T20:32:51-03:00"
  },
  {
    "companyDocument": "77723018000146",
    "companyName": "123 S.A.",
    "customerDocument": "51537476467",
    "value": 400.00,
    "contract": "5f206825-3cfe-412f-8302-cc1b24a179b0",
    "debtDate": "2015-10-12T20:32:51-03:00",
    "inclusionDate": "2020-10-12T20:32:51-03:00"
  }
]
`
		context.Writer.WriteHeader(http.StatusOK)
		_, _ = context.Writer.WriteString(data)
	})
	server := httptest.NewServer(handler)
	baseUrl := fmt.Sprintf("http://%s", server.Listener.Addr())
	client, err := infrastructure.NewClient("localhost", 8529, "root", "somepassword")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := infrastructure.CreateDatabase(client, "negativation")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	repo, err := infrastructure.NewNegativationRepositoryArangoDB(database)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	legacyRepo, err := infrastructure.NewNegativationLegacyRepositoryAPI(baseUrl)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collection, err := infrastructure.CreateCollection(database, "negativations")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	err = infrastructure.CleanCollection(t, collection)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	firstNegativation, err := domain.NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	secondNegativation, err := domain.NewNegativation("04843574000182", "DBZ S.A.", "26658236674", 59.99, "3132f136-3889-4efb-bf92-e1efbb3fe15e", "2015-09-11T20:32:51-03:00", "2020-09-11T20:32:51-03:00")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	err = infrastructure.InsertNegativations(t, collection, firstNegativation, secondNegativation)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Synchronize negativations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		sut := NewLegacyNegativationController(repo, legacyRepo, symmetricKey, encryptionContext)

		err := sut.Synchronize()

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(infrastructure.GetNegativations(t, collection)).Should(
			And(
				HaveLen(2),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("59291534000167"),
					"CompanyName":      BeEquivalentTo("ABC S.A."),
					"CustomerDocument": BeEquivalentTo("FUNoBV5JaZEji6c="),
					"Value":            BeEquivalentTo(2500.23),
					"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
				}))),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("77723018000146"),
					"CompanyName":      BeEquivalentTo("123 S.A."),
					"CustomerDocument": BeEquivalentTo("FUNoBV5JaZEji6c="),
					"Value":            BeEquivalentTo(400.00),
					"Contract":         BeEquivalentTo("5f206825-3cfe-412f-8302-cc1b24a179b0"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 10, 12, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 10, 12, 23, 32, 51, 0, time.UTC)),
				})))))
	})

}
