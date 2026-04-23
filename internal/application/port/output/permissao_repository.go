package output

import "github.com/renatofagalde/app-openfinance-fake/internal/application/domain"

type PermissaoRepository interface {
	FindByConsentId(consentId string) ([]domain.Permissao, error)
	FindByConsentIdAndPermission(consentId string, permission string) (*domain.Permissao, error)
	Save(permissao domain.Permissao) error
	UpdateLancar403(consentId string, permission string, lancar403 bool) error
}
