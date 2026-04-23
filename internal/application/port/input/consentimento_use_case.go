package input

import "github.com/renatofagalde/app-openfinance-fake/internal/application/domain"

type ConsentimentoUseCase interface {
	GetByConsentId(consentId string) (*domain.Consentimento, error)
	Save(consentimento domain.Consentimento) error
	GetAll() ([]domain.Consentimento, error)
	UpdateStatus(consentId string, status string) error
}
