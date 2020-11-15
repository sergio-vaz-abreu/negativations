package application

import (
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/internal/infrastructure"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"testing"
	"time"
)

func TestNegativationController(t *testing.T) {
	g := NewGomegaWithT(t)
	client, err := infrastructure.NewClient()
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := infrastructure.CreateDatabase(client, "negativation")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	repo, err := infrastructure.NewNegativationRepositoryArangoDB(database)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collection, err := infrastructure.CreateCollection(database, "negativations")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	err = infrastructure.CleanCollection(t, collection)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	negativation, err := domain.NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	err = infrastructure.InsertNegativations(t, collection, negativation)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	sut := NewNegativationController(repo)

	t.Run("Validate cpf and get negativations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		cpf := "515.374.764-67"

		negativations, err := sut.GetByCPF(cpf)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(negativations).Should(
			ContainElement(PointTo(MatchAllFields(Fields{
				"CompanyDocument":  BeEquivalentTo("59291534000167"),
				"CompanyName":      BeEquivalentTo("ABC S.A."),
				"CustomerDocument": BeEquivalentTo("51537476467"),
				"Value":            BeEquivalentTo(1235.23),
				"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
				"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
				"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
			}))))
	})
}
