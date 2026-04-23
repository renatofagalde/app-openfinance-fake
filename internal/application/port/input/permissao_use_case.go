package input

import "github.com/renatofagalde/app-openfinance-fake/internal/application/domain"

type PermissaoUseCase interface {
	GetByConsentId(consentId string) ([]domain.Permissao, error)
	Save(permissao domain.Permissao) error
	UpdateLancar403(consentId string, permission string, lancar403 bool) error
}
