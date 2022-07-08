package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mrinjamul/gnote/models"
	"github.com/mrinjamul/gnote/utils"
)

var (
	jwtKey string
)

func init() {
	jwtKey = utils.GetEnv("JWT_SECRET")
	if jwtKey == "" {
		panic("JWT_SECRET not set")
	}
}

// JWTAuth is a middleware for validating JWT tokens
func JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check if token is present
		if ctx.Request.Header.Get("Authorization") == "" {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}
		// get token from header
		tokenString, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}

		claims := &models.Claims{}
		// check if token is expired
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad token",
			})
			ctx.Abort()
			return
		}

		// check if token is expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// JWTAuthCookie is a middleware for validating JWT tokens through cookie
func JWTAuthCookie() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get cookie "token"
		tokenString, err := ctx.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "cookie not found",
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad cookie",
			})
		}

		claims := &models.Claims{}
		// check if token is expired
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad token",
			})
			ctx.Abort()
			return
		}

		// check if token is expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// JWTAuth is a middleware for validating JWT tokens
func JWTAuthAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check if token is present
		if ctx.Request.Header.Get("Authorization") == "" {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}
		// get token from header
		tokenString, err := utils.ParseToken(ctx.Request.Header.Get("Authorization"))
		if err != nil {
			ctx.JSON(401, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}

		claims := &models.Claims{}
		// check if token is expired
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			ctx.Abort()
			return
		}
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
				})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "bad token",
			})
			ctx.Abort()
			return
		}

		// check if token is expired
		if time.Now().Unix() > claims.ExpiresAt.Unix() {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "token expired",
			})
			ctx.Abort()
			return
		}

		// check if user is admin
		if claims.Role != "admin" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
