package service

import (
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/port/output"
)

type PermissaoService struct {
	repo output.PermissaoRepository
}

func NewPermissaoService(repo output.PermissaoRepository) *PermissaoService {
	return &PermissaoService{repo: repo}
}

func (s *PermissaoService) GetByConsentId(consentId string) ([]domain.Permissao, error) {
	return s.repo.FindByConsentId(consentId)
}

func (s *PermissaoService) Save(permissao domain.Permissao) error {
	return s.repo.Save(permissao)
}

func (s *PermissaoService) UpdateLancar403(consentId string, permission string, lancar403 bool) error {
	return s.repo.UpdateLancar403(consentId, permission, lancar403)
}
