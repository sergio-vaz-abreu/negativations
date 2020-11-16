package domain

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"testing"
	"time"
)

const (
	symmetricKey      = "UkVDMgAAAC13PCVZAKOczZXUpvkhsC+xvwWnv3CLmlG0Wzy8ZBMnT+2yx/dg"
	encryptionContext = "context"
)

func TestEncryptNegativation(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
	g.Expect(err).Should(
		Not(HaveOccurred()))

	err = sut.Encrypt(symmetricKey, encryptionContext)

	g.Expect(err).Should(
		Not(HaveOccurred()))
	g.Expect(sut).Should(
		PointTo(MatchAllFields(Fields{
			"CompanyDocument":  BeEquivalentTo("59291534000167"),
			"CompanyName":      BeEquivalentTo("ABC S.A."),
			"CustomerDocument": BeEquivalentTo("FUNoBV5JaZEji6c="),
			"Value":            BeEquivalentTo(1235.23),
			"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
			"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
			"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
		})))
}

func TestDecryptNegativation(t *testing.T) {
	g := NewGomegaWithT(t)
	sut, err := NewNegativation("59291534000167", "ABC S.A.", "51537476467", 1235.23, "bc063153-fb9e-4334-9a6c-0d069a42065b", "2015-11-13T20:32:51-03:00", "2020-11-13T20:32:51-03:00")
	g.Expect(err).Should(
		Not(HaveOccurred()))
	err = sut.Encrypt(symmetricKey, encryptionContext)
	g.Expect(err).Should(
		Not(HaveOccurred()))

	err = sut.Decrypt(symmetricKey, encryptionContext)

	g.Expect(err).Should(
		Not(HaveOccurred()))
	g.Expect(sut).Should(
		PointTo(MatchAllFields(Fields{
			"CompanyDocument":  BeEquivalentTo("59291534000167"),
			"CompanyName":      BeEquivalentTo("ABC S.A."),
			"CustomerDocument": BeEquivalentTo("51537476467"),
			"Value":            BeEquivalentTo(1235.23),
			"Contract":         BeEquivalentTo("bc063153-fb9e-4334-9a6c-0d069a42065b"),
			"DebtDate":         BeEquivalentTo(time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC)),
			"InclusionDate":    BeEquivalentTo(time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC)),
		})))
}
