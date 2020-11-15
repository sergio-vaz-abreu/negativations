package domain

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestCreateValidCPF(t *testing.T) {
	rawCpfs := []string{
		"529.982.247-25",
		"529.982.247.25",
		"529.982.24725",
		"529.98224725",
		"529982247-25",
		"52998224725",
		" 52998224725 ",
		" 529 982 247 25 ",
	}
	for _, rawCpf := range rawCpfs {
		t.Run("Create valid cpf from all formats", func(t *testing.T) {
			g := NewGomegaWithT(t)

			cpf, err := NewCPF(rawCpf)

			g.Expect(err).Should(
				Not(HaveOccurred()))
			g.Expect(cpf).Should(
				BeEquivalentTo("52998224725"))
		})
	}
}

type invalidCases struct {
	cpf         string
	expectedErr error
}

func TestDoNotCreateInvalidCPF(t *testing.T) {

	invalidCases := []invalidCases{
		{"529.982.2e47-25", NonValidCharacter},
		{"529,982.247.25", NonValidCharacter},
		{"529.982.247;25", NonValidCharacter},
		{"529.982.247-255", InvalidNumberOfCharacter},
		{"5529.982.247-25", InvalidNumberOfCharacter},
		// verification digit not valid
	}
	for _, invalidCase := range invalidCases {
		t.Run("Do not create invalid cpf", func(t *testing.T) {
			g := NewGomegaWithT(t)

			_, err := NewCPF(invalidCase.cpf)

			g.Expect(err).Should(
				MatchError(invalidCase.expectedErr))
		})
	}
}

func TestEncryptCpf(t *testing.T) {
	g := NewGomegaWithT(t)
	sut := CPF("51537476467")

	encryptedCpf, err := sut.Encrypt(symmetricKey, encryptionContext)

	g.Expect(err).Should(
		Not(HaveOccurred()))
	g.Expect(encryptedCpf).Should(
		BeEquivalentTo("FUNoBV5JaZEji6c="))
}

func TestDecryptCpf(t *testing.T) {
	g := NewGomegaWithT(t)
	sut := CPF("FUNoBV5JaZEji6c=")

	decryptedCpf, err := sut.Decrypt(symmetricKey, encryptionContext)

	g.Expect(err).Should(
		Not(HaveOccurred()))
	g.Expect(decryptedCpf).Should(
		BeEquivalentTo("51537476467"))
}
