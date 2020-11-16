package application

import (
	"github.com/negativations/modules/negativation/infrastructure"
)

func NewLegacyNegativationController(repository infrastructure.NegativationRepository, legacyRepository infrastructure.NegativationLegacyRepository, symmetricKey string, encryptionContext string) *LegacyNegativationController {
	return &LegacyNegativationController{repository: repository, legacyRepository: legacyRepository, symmetricKey: symmetricKey, encryptionContext: encryptionContext}
}

type LegacyNegativationController struct {
	repository        infrastructure.NegativationRepository
	legacyRepository  infrastructure.NegativationLegacyRepository
	symmetricKey      string
	encryptionContext string
}

func (ctrl *LegacyNegativationController) Synchronize() error {
	negativations, err := ctrl.legacyRepository.GetAll()
	if err != nil {
		return err
	}
	for _, negativation := range negativations {
		negativation.UTC()
		err = negativation.Encrypt(ctrl.symmetricKey, ctrl.encryptionContext)
		if err != nil {
			return err
		}
	}
	return ctrl.repository.Synchronize(negativations...)
}
