package infrastructure

import (
	"fmt"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNegativationLegacyRepositoryAPI(t *testing.T) {
	handler := gin.Default()
	handler.GET("negativation", func(context *gin.Context) {
		data := `
[
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
`
		context.Writer.WriteHeader(http.StatusOK)
		_, _ = context.Writer.WriteString(data)
	})
	server := httptest.NewServer(handler)
	baseUrl := fmt.Sprintf("http://%s", server.Listener.Addr())

	t.Run("Getting all negativations from legacy api", func(t *testing.T) {
		g := NewGomegaWithT(t)
		sut, err := NewNegativationLegacyRepositoryAPI(baseUrl)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		maintenances, err := sut.GetAll()

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(maintenances).Should(
			And(
				HaveLen(1),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("59291534000167"),
					"CompanyName":      BeEquivalentTo("ABC S.A."),
					"CustomerDocument": BeEquivalentTo("51537476467"),
					"Value":            BeEquivalentTo(1235.23),
					"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
				})))))
	})
}
