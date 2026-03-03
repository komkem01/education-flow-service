package auth

import "github.com/gin-gonic/gin"

const ginAuthClaimsKey = "auth.claims"

func SetClaimsToGin(ctx *gin.Context, claims *TokenClaims) {
	ctx.Set(ginAuthClaimsKey, claims)
}

func GetClaimsFromGin(ctx *gin.Context) (*TokenClaims, bool) {
	value, ok := ctx.Get(ginAuthClaimsKey)
	if !ok {
		return nil, false
	}

	claims, ok := value.(*TokenClaims)
	if !ok || claims == nil {
		return nil, false
	}

	return claims, true
}
