package infrastructure

import (
	"github.com/negativations/modules/negativation/domain"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"testing"
	"time"
)

const (
	databaseName   = "negativation"
	collectionName = "negativations"
)

func TestNegativationRepositoryArangoDB(t *testing.T) {
	g := NewGomegaWithT(t)
	client, err := NewClient("localhost", 8529, "root", "somepassword")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := CreateDatabase(client, databaseName)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collection, err := CreateCollection(database, collectionName)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Get empty negativation by cpf when not stored", func(t *testing.T) {
		g := NewGomegaWithT(t)
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = CleanCollection(t, collection)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		cpf := domain.CPF("51537476467")

		negativations, err := sut.GetByCPF(cpf)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(negativations).Should(
			BeEmpty())
	})

	t.Run("Get negativations by cpf when stored", func(t *testing.T) {
		g := NewGomegaWithT(t)
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = CleanCollection(t, collection)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		firstNegativation, err := domain.NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		secondNegativation, err := domain.NewNegativation("77723018000146", "123 S.A.", "51537476467", 400.00, "5f206825-3cfe-412f-8302-cc1b24a179b0", "2015-10-12T20:32:51-03:00", "2020-10-12T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		thirdNegativation, err := domain.NewNegativation("04843574000182", "DBZ S.A.", "26658236674", 59.99, "3132f136-3889-4efb-bf92-e1efbb3fe15e", "2015-09-11T20:32:51-03:00", "2020-09-11T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = InsertNegativations(t, collection, firstNegativation, secondNegativation, thirdNegativation)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		cpf := domain.CPF("51537476467")

		negativations, err := sut.GetByCPF(cpf)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(negativations).Should(
			And(
				HaveLen(2),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("59291534000167"),
					"CompanyName":      BeEquivalentTo("ABC S.A."),
					"CustomerDocument": BeEquivalentTo("51537476467"),
					"Value":            BeEquivalentTo(1235.23),
					"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
				}))),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("77723018000146"),
					"CompanyName":      BeEquivalentTo("123 S.A."),
					"CustomerDocument": BeEquivalentTo("51537476467"),
					"Value":            BeEquivalentTo(400.00),
					"Contract":         BeEquivalentTo("5f206825-3cfe-412f-8302-cc1b24a179b0"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 10, 12, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 10, 12, 23, 32, 51, 0, time.UTC)),
				})))))
	})

	t.Run("Sync negativations", func(t *testing.T) {
		g := NewGomegaWithT(t)
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = CleanCollection(t, collection)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		firstNegativation, err := domain.NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		secondNegativation, err := domain.NewNegativation("77723018000146", "123 S.A.", "51537476467", 400.00, "5f206825-3cfe-412f-8302-cc1b24a179b0", "2015-10-12T20:32:51-03:00", "2020-10-12T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		thirdNegativation, err := domain.NewNegativation("04843574000182", "DBZ S.A.", "26658236674", 59.99, "3132f136-3889-4efb-bf92-e1efbb3fe15e", "2015-09-11T20:32:51-03:00", "2020-09-11T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = InsertNegativations(t, collection, firstNegativation, secondNegativation, thirdNegativation)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		fourthNegativation, err := domain.NewNegativation("77723018000146", "123 S.A.", "26658236674", 2500.00, "3132f136-3889-4efb-bf92-e1efbb3fe15e", "2015-09-11T20:32:51-03:00", "2020-09-11T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))
		fifthNegativation, err := domain.NewNegativation("70170935000100", "ASD S.A.", "25124543043", 10340.67, "d6628a0e-d4dd-4f14-8591-2ddc7f1bbeff", "2015-07-09T20:32:51-03:00", "2020-07-09T20:32:51-03:00")
		g.Expect(err).Should(
			Not(HaveOccurred()))

		err = sut.Synchronize(fourthNegativation, fifthNegativation)

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(GetNegativations(t, collection)).Should(
			And(
				HaveLen(2),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("70170935000100"),
					"CompanyName":      BeEquivalentTo("ASD S.A."),
					"CustomerDocument": BeEquivalentTo("25124543043"),
					"Value":            BeEquivalentTo(10340.67),
					"Contract":         BeEquivalentTo("d6628a0e-d4dd-4f14-8591-2ddc7f1bbeff"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 7, 9, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 7, 9, 23, 32, 51, 0, time.UTC)),
				}))),
				ContainElement(PointTo(MatchAllFields(Fields{
					"CompanyDocument":  BeEquivalentTo("77723018000146"),
					"CompanyName":      BeEquivalentTo("123 S.A."),
					"CustomerDocument": BeEquivalentTo("26658236674"),
					"Value":            BeEquivalentTo(2500.00),
					"Contract":         BeEquivalentTo("3132f136-3889-4efb-bf92-e1efbb3fe15e"),
					"DebtDate":         BeEquivalentTo(time.Date(2015, 9, 11, 23, 32, 51, 0, time.UTC)),
					"InclusionDate":    BeEquivalentTo(time.Date(2020, 9, 11, 23, 32, 51, 0, time.UTC)),
				})))))
	})
}
