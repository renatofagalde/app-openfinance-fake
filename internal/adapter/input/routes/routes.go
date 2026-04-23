package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/renatofagalde/app-openfinance-fake/internal/adapter/input/controller"
)

func Setup(
	r *gin.Engine,
	tokenCtrl *controller.TokenController,
	accountsCtrl *controller.AccountsController,
	mockCtrl *controller.MockController,
) {
	// Token — mesmo path que o Keycloak real
	r.POST("/auth/realms/:realm/protocol/openid-connect/token", tokenCtrl.GetToken)

	// Open Finance endpoints — mesmos paths que o robô C# chama
	r.GET("/accounts", accountsCtrl.GetAccounts)
	r.GET("/accounts/:accountId", accountsCtrl.GetAccountById)
	r.GET("/accounts/:accountId/balances", accountsCtrl.GetBalances)
	r.GET("/accounts/:accountId/transactions", accountsCtrl.GetTransactions)
	r.GET("/accounts/:accountId/transactions-current", accountsCtrl.GetTransactionsCurrent)
	r.GET("/accounts/:accountId/overdraft-limits", accountsCtrl.GetOverdraftLimits)

	// Mock endpoints — para manipular dados via Postman
	mock := r.Group("/mock")
	{
		mock.GET("/consentimentos", mockCtrl.ListConsentimentos)
		mock.POST("/consentimentos", mockCtrl.SaveConsentimento)
		mock.PUT("/consentimentos/:consentId/status", mockCtrl.UpdateConsentimentoStatus)

		mock.GET("/permissoes/:consentId", mockCtrl.ListPermissoes)
		mock.POST("/permissoes", mockCtrl.SavePermissao)
		mock.PUT("/permissoes/:consentId/:permission/lancar403", mockCtrl.UpdateLancar403)
	}
}
