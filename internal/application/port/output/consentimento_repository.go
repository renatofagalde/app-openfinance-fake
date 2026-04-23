package output

import "github.com/renatofagalde/app-openfinance-fake/internal/application/domain"

type ConsentimentoRepository interface {
	FindByConsentId(consentId string) (*domain.Consentimento, error)
	Save(consentimento domain.Consentimento) error
	FindAll() ([]domain.Consentimento, error)
	UpdateStatus(consentId string, status string) error
}
