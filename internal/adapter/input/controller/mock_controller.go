package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/domain"
	"github.com/renatofagalde/app-openfinance-fake/internal/application/port/input"
)

type MockController struct {
	consentimentoUseCase input.ConsentimentoUseCase
	permissaoUseCase     input.PermissaoUseCase
}

func NewMockController(
	consentimentoUseCase input.ConsentimentoUseCase,
	permissaoUseCase input.PermissaoUseCase,
) *MockController {
	return &MockController{
		consentimentoUseCase: consentimentoUseCase,
		permissaoUseCase:     permissaoUseCase,
	}
}

// GET /mock/consentimentos
func (m *MockController) ListConsentimentos(c *gin.Context) {
	items, err := m.consentimentoUseCase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /mock/consentimentos
func (m *MockController) SaveConsentimento(c *gin.Context) {
	var body domain.Consentimento
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.consentimentoUseCase.Save(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

// PUT /mock/consentimentos/:consentId/status
func (m *MockController) UpdateConsentimentoStatus(c *gin.Context) {
	consentId := c.Param("consentId")
	var body struct {
		Status string `json:"consent_status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.consentimentoUseCase.UpdateStatus(consentId, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"consent_id": consentId, "consent_status": body.Status})
}

// GET /mock/permissoes/:consentId
func (m *MockController) ListPermissoes(c *gin.Context) {
	consentId := c.Param("consentId")
	items, err := m.permissaoUseCase.GetByConsentId(consentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// POST /mock/permissoes
func (m *MockController) SavePermissao(c *gin.Context) {
	var body domain.Permissao
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.permissaoUseCase.Save(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, body)
}

// PUT /mock/permissoes/:consentId/:permission/lancar403
func (m *MockController) UpdateLancar403(c *gin.Context) {
	consentId := c.Param("consentId")
	permission := c.Param("permission")
	var body struct {
		Lancar403 bool `json:"lancar_403"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := m.permissaoUseCase.UpdateLancar403(consentId, permission, body.Lancar403); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"consent_id": consentId,
		"permission": permission,
		"lancar_403": body.Lancar403,
	})
}
