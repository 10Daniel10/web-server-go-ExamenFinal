package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthKeys struct {
	SecretKey string
	PubKey    string
}

func NewAuthKeys(secretKey string, pubKey string) *AuthKeys {
	return &AuthKeys{
		SecretKey: secretKey,
		PubKey:    pubKey,
	}
}

func (v *AuthKeys) Validate(ctx *gin.Context) {
	secretKey := ctx.GetHeader("SECRET_KEY")
	pubKey := ctx.GetHeader("PUBLIC_KEY")

	if secretKey == v.SecretKey && pubKey == v.PubKey {
		ctx.Next()
	}

	ctx.AbortWithStatus(http.StatusUnauthorized)
}
