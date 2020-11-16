package application

import (
	"github.com/negativations/modules/negativation/domain"
	"github.com/negativations/modules/negativation/infrastructure"
)

func NewNegativationController(repository infrastructure.NegativationRepository, symmetricKey string, encryptionContext string) *NegativationController {
	return &NegativationController{repository: repository, symmetricKey: symmetricKey, encryptionContext: encryptionContext}
}

type NegativationController struct {
	repository        infrastructure.NegativationRepository
	symmetricKey      string
	encryptionContext string
}

func (ctrl *NegativationController) GetByCPF(rawCpf string) ([]*domain.Negativation, error) {
	cpf, err := domain.NewCPF(rawCpf)
	if err != nil {
		return nil, err
	}
	encryptedCpf, err := cpf.Encrypt(ctrl.symmetricKey, ctrl.encryptionContext)
	if err != nil {
		return nil, err
	}
	negativations, err := ctrl.repository.GetByCPF(encryptedCpf)
	if err != nil {
		return nil, err
	}
	for _, negativation := range negativations {
		err = negativation.Decrypt(ctrl.symmetricKey, ctrl.encryptionContext)
		if err != nil {
			return nil, err
		}
	}
	return negativations, nil
}
