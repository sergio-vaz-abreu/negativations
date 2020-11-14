package domain

import "time"

func NewNegativation(
	companyDocument string,
	companyName string,
	customerDocument CPF,
	value float64,
	contract string,
	debtDate time.Time,
	inclusionDate time.Time,
) *Negativation {
	return &Negativation{
		CompanyDocument:  companyDocument,
		CompanyName:      companyName,
		CustomerDocument: customerDocument,
		Value:            value,
		Contract:         contract,
		DebtDate:         debtDate,
		InclusionDate:    inclusionDate,
	}
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
