package domain

import (
	"github.com/pkg/errors"
	"time"
)

var (
	InvalidDebtDateFormat      = errors.New("invalid debt date format")
	InvalidInclusionDateFormat = errors.New("invalid inclusion date format")
)

func NewNegativation(
	companyDocument string,
	companyName string,
	customerDocument string,
	value float64,
	contract string,
	rawDebtDate string,
	rawInclusionDate string,
) (*Negativation, error) {
	debtDate, err := time.Parse(time.RFC3339, rawDebtDate)
	if err != nil {
		return nil, InvalidDebtDateFormat
	}
	inclusionDate, err := time.Parse(time.RFC3339, rawInclusionDate)
	if err != nil {
		return nil, InvalidInclusionDateFormat
	}
	cpf, err := NewCPF(customerDocument)
	if err != nil {
		return nil, err
	}
	return &Negativation{
		CompanyDocument:  companyDocument,
		CompanyName:      companyName,
		CustomerDocument: cpf,
		Value:            value,
		Contract:         contract,
		DebtDate:         debtDate.UTC(),
		InclusionDate:    inclusionDate.UTC(),
	}, nil
}

type Negativation struct {
	CompanyDocument  string    `json:"companyDocument"`
	CompanyName      string    `json:"companyName"`
	CustomerDocument CPF       `json:"customerDocument"`
	Value            float64   `json:"value"`
	Contract         string    `json:"contract"`
	DebtDate         time.Time `json:"debtDate"`
	InclusionDate    time.Time `json:"inclusionDate"`
}

func (negativation *Negativation) UTC() {
	negativation.DebtDate = negativation.DebtDate.UTC()
	negativation.InclusionDate = negativation.InclusionDate.UTC()
}
