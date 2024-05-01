package token_controller

import (
	"cep-gin-clean-arch/internal/entity"
	"cep-gin-clean-arch/utils"

	"github.com/gin-gonic/gin"
)

const (
	erroAoGerarToken = "Erro interno ao tentar gerar o token JWT"
)

type GerarTokenHandler struct {
	GerarTokenInterface entity.GerarTokenInterface
}

func NewGerarTokenHandler(gerarTokenRepository entity.GerarTokenInterface) *GerarTokenHandler {
	return &GerarTokenHandler{
		GerarTokenInterface: gerarTokenRepository,
	}
}

func (h *GerarTokenHandler) GerarTokenJWT(c *gin.Context) {

	token, err := h.GerarTokenInterface.GenerateTokenJWT()
	if err != nil {
		utils.GravarErroNoSentry(err, c)
		c.AbortWithStatusJSON(500, erroAoGerarToken)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})

}
