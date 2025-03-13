package middleware

import (
	"UserCrud/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	ValidateAndExtractJwt() gin.HandlerFunc
}

const (
	JWTClaimsContextKey = "JWTClaimsContextKey"
)

type authMiddleware struct {
	jwt util.JwtUtil
}

func (a *authMiddleware) ValidateAndExtractJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Authorization header is empty",
			})
			return
		}
		header := strings.Fields(authHeader)
		if len(header) != 2 && header[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Authorization header is invalid",
			})
			return
		}
		accessToken := header[1]
		claims, err := a.jwt.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
		}

		c.Set(JWTClaimsContextKey, claims)
		c.Next()
	}
}

func NewAuthMiddleware(jwt util.JwtUtil) AuthMiddleware {
	return &authMiddleware{jwt: jwt}
}
