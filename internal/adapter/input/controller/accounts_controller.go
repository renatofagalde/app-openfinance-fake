package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/port/input"
)

type AccountsController struct {
	consentimentoUseCase input.ConsentimentoUseCase
	permissaoUseCase     input.PermissaoUseCase
}

func NewAccountsController(
	consentimentoUseCase input.ConsentimentoUseCase,
	permissaoUseCase input.PermissaoUseCase,
) *AccountsController {
	return &AccountsController{
		consentimentoUseCase: consentimentoUseCase,
		permissaoUseCase:     permissaoUseCase,
	}
}

func (a *AccountsController) verificarConsentimento(c *gin.Context) (*domain.Consentimento, bool) {
	consentId := c.GetHeader("x-consent-id")
	if consentId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "x-consent-id header obrigatório"})
		return nil, false
	}

	consentimento, err := a.consentimentoUseCase.GetByConsentId(consentId)
	if err != nil || consentimento == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"code": "CONSENTIMENTO_NAO_ENCONTRADO", "title": "Consentimento não encontrado"}},
		})
		return nil, false
	}

	if consentimento.ConsentStatus != domain.StatusAuthorised {
		c.JSON(http.StatusUnauthorized, gin.H{
			"errors": []gin.H{{"code": "CONSENTIMENTO_NAO_AUTORIZADO", "title": "Consentimento não autorizado"}},
		})
		return nil, false
	}

	return consentimento, true
}

func (a *AccountsController) verificarPermissao(c *gin.Context, consentId string, permission string) bool {
	permissao, err := a.permissaoUseCase.GetByConsentId(consentId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"errors": []gin.H{{"code": "FORBIDDEN", "title": "Acesso negado"}},
		})
		return false
	}

	for _, p := range permissao {
		if p.Permission == permission {
			if p.Lancar403 {
				c.JSON(http.StatusForbidden, gin.H{
					"errors": []gin.H{{"code": "FORBIDDEN", "title": "Permissão negada para " + permission}},
				})
				return false
			}
			return true
		}
	}

	c.JSON(http.StatusForbidden, gin.H{
		"errors": []gin.H{{"code": "FORBIDDEN", "title": "Permission " + permission + " não concedida"}},
	})
	return false
}

// GET /accounts
func (a *AccountsController) GetAccounts(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"brandName":   "Fake Bank",
				"companyCnpj": "00000000000191",
				"type":        "CONTA_DEPOSITO_A_VISTA",
				"compeCode":   "001",
				"branchCode":  "0001",
				"number":      "12345678",
				"checkDigit":  "9",
				"accountId":   "fake-account-id-001",
			},
		},
		"meta": gin.H{"totalRecords": 1, "totalPages": 1, "requestDateTime": "2026-04-22T10:00:00Z"},
	})
}

// GET /accounts/:accountId
func (a *AccountsController) GetAccountById(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"compeCode":  "001",
			"branchCode": "0001",
			"number":     "12345678",
			"checkDigit": "9",
			"type":       "CONTA_DEPOSITO_A_VISTA",
			"subtype":    "INDIVIDUAL",
			"currency":   "BRL",
		},
		"meta": gin.H{"totalRecords": 1, "totalPages": 1, "requestDateTime": "2026-04-22T10:00:00Z"},
	})
}

// GET /accounts/:accountId/balances
func (a *AccountsController) GetBalances(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_BALANCES_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"availableAmount":             gin.H{"amount": "1000.00", "currency": "BRL"},
			"blockedAmount":               gin.H{"amount": "0.00", "currency": "BRL"},
			"automaticallyInvestedAmount": gin.H{"amount": "0.00", "currency": "BRL"},
			"updateDateTime":              "2026-04-22T10:00:00Z",
		},
		"meta": gin.H{"totalRecords": 1, "totalPages": 1, "requestDateTime": "2026-04-22T10:00:00Z"},
	})
}

// GET /accounts/:accountId/transactions
func (a *AccountsController) GetTransactions(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_TRANSACTIONS_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"transactionId":                  "fake-txn-001",
				"completedAuthorisedPaymentType": "TRANSACAO_EFETIVADA",
				"creditDebitType":                "DEBITO",
				"transactionName":                "Fake Transaction",
				"type":                           "PIX",
				"transactionAmount":              gin.H{"amount": "100.00", "currency": "BRL"},
				"transactionDateTime":            "2026-04-22T10:00:00Z",
			},
		},
		"meta": gin.H{"requestDateTime": "2026-04-22T10:00:00Z"},
	})
}

// GET /accounts/:accountId/transactions-current
func (a *AccountsController) GetTransactionsCurrent(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_TRANSACTIONS_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": []gin.H{
			{
				"transactionId":                  "fake-txn-current-001",
				"completedAuthorisedPaymentType": "TRANSACAO_EFETIVADA",
				"creditDebitType":                "CREDITO",
				"transactionName":                "Fake Current Transaction",
				"type":                           "TED",
				"transactionAmount":              gin.H{"amount": "500.00", "currency": "BRL"},
				"transactionDateTime":            "2026-04-22T10:00:00Z",
			},
		},
		"meta": gin.H{"requestDateTime": "2026-04-22T10:00:00Z"},
	})
}

// GET /accounts/:accountId/overdraft-limits
func (a *AccountsController) GetOverdraftLimits(c *gin.Context) {
	consentimento, ok := a.verificarConsentimento(c)
	if !ok {
		return
	}
	if !a.verificarPermissao(c, consentimento.ConsentId, "ACCOUNTS_OVERDRAFT_LIMITS_READ") {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"overdraftContractedLimit":  gin.H{"amount": "5000.00", "currency": "BRL"},
			"overdraftUsedLimit":        gin.H{"amount": "0.00", "currency": "BRL"},
			"unarrangedOverdraftAmount": gin.H{"amount": "0.00", "currency": "BRL"},
		},
		"meta": gin.H{"totalRecords": 1, "totalPages": 1, "requestDateTime": "2026-04-22T10:00:00Z"},
	})
}
