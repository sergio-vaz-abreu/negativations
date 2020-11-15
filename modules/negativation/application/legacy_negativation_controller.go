package application

import (
	"github.com/negativations/modules/negativation/infrastructure"
)

func NewLegacyNegativationController(repository infrastructure.NegativationRepository, legacyRepository infrastructure.NegativationLegacyRepository) *LegacyNegativationController {
	return &LegacyNegativationController{repository: repository, legacyRepository: legacyRepository}
}

type LegacyNegativationController struct {
	repository       infrastructure.NegativationRepository
	legacyRepository infrastructure.NegativationLegacyRepository
}

func (ctrl *LegacyNegativationController) Synchronize() error {
	negativations, err := ctrl.legacyRepository.GetAll()
	if err != nil {
		return err
	}
	for _, negativation := range negativations {
		negativation.UTC()
	}
	return ctrl.repository.Synchronize(negativations...)
}
