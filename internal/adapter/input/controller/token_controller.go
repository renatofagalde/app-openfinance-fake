package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct{}

func NewTokenController() *TokenController {
	return &TokenController{}
}

// POST /auth/realms/:realm/protocol/openid-connect/token
func (t *TokenController) GetToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"access_token": "fake-token-openfinance-fake-12345",
		"token_type":   "Bearer",
		"expires_in":   300,
		"scope":        "openid",
	})
}
