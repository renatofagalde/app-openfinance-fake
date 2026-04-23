package service

import (
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/port/output"
)

type ConsentimentoService struct {
	repo output.ConsentimentoRepository
}

func NewConsentimentoService(repo output.ConsentimentoRepository) *ConsentimentoService {
	return &ConsentimentoService{repo: repo}
}

func (s *ConsentimentoService) GetByConsentId(consentId string) (*domain.Consentimento, error) {
	return s.repo.FindByConsentId(consentId)
}

func (s *ConsentimentoService) Save(consentimento domain.Consentimento) error {
	return s.repo.Save(consentimento)
}

func (s *ConsentimentoService) GetAll() ([]domain.Consentimento, error) {
	return s.repo.FindAll()
}

func (s *ConsentimentoService) UpdateStatus(consentId string, status string) error {
	return s.repo.UpdateStatus(consentId, status)
}
