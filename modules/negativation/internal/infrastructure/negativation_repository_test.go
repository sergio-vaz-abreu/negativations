package infrastructure

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	. "github.com/onsi/gomega"
	"testing"
)

func TestNegativationRepositoryArangoDB(t *testing.T) {
	g := NewGomegaWithT(t)
	client, err := GetClient()
	g.Expect(err).Should(
		Not(HaveOccurred()))
	database, err := GetDatabase(client, "negativation")
	g.Expect(err).Should(
		Not(HaveOccurred()))

	t.Run("Get empty negativation by cpf when not stored", func(t *testing.T) {
		sut, err := NewNegativationRepositoryArangoDB(database)
		g.Expect(err).Should(
			Not(HaveOccurred()))

		negativations, err := sut.GetByCPF("00000000000")

		g.Expect(err).Should(
			Not(HaveOccurred()))
		g.Expect(negativations).Should(
			BeEmpty())
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
