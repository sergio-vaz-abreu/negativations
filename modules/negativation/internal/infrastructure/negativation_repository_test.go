package infrastructure

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
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
	client, err := GetClient()
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := GetDatabase(client, databaseName)
	g.Expect(err).Should(
		Not(HaveOccurred()))
	collectionName, err := GetCollection(database, collectionName)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Get empty negativation by cpf when not stored", func(t *testing.T) {
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = CleanCollection(collectionName)
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
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))
		err = CleanCollection(collectionName)
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
		err = InsertNegativations(collectionName, firstNegativation, secondNegativation, thirdNegativation)
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
}

func GetClient() (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"tcp://localhost:8529"},
	})
	if err != nil {
		return nil, err
	}
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication("root", "somepassword"),
	})
	return client, err
}

func GetDatabase(client driver.Client, dbName string) (driver.Database, error) {
	exists, err := client.DatabaseExists(nil, dbName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return client.CreateDatabase(driver.WithWaitForSync(nil, true), dbName, nil)
	}
	return client.Database(nil, dbName)
}

func CleanCollection(collection driver.Collection) error {
	return collection.Truncate(nil)
}

func GetCollection(database driver.Database, collectionName string) (driver.Collection, error) {
	exists, err := database.CollectionExists(nil, collectionName)
	if err != nil {
		return nil, err
	}
	if !exists {
		return database.CreateCollection(nil, collectionName, nil)
	}
	return database.Collection(nil, collectionName)
}

func InsertNegativations(collection driver.Collection, negativations ...*domain.Negativation) error {
	_, _, err := collection.CreateDocuments(driver.WithWaitForSync(nil, true), negativations)
	return err
}
