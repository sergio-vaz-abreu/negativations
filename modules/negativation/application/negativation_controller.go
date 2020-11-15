package application

import (
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/infrastructure"
)

func NewNegativationController(repository infrastructure.NegativationRepository) *NegativationController {
	return &NegativationController{repository: repository}
}

type NegativationController struct {
	repository infrastructure.NegativationRepository
}

func (ctrl *NegativationController) GetByCPF(rawCpf string) ([]*domain.Negativation, error) {
	cpf, err := domain.NewCPF(rawCpf)
	if err != nil {
		return nil, err
	}
	return ctrl.repository.GetByCPF(cpf)
}
