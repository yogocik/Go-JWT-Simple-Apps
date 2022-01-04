package utility

import (
	"errors"
	"fmt"
	response "integration/presentation/model/response"
	authenticator "integration/process/authentication"
	"strings"

	"github.com/gin-gonic/gin"
)

type authHeader struct {
	AuthorizationHeader string `header:"Authorization"`
}

type AuthTokenMiddleware struct {
	acctToken authenticator.Token
}

func NewTokenValidator(acctToken authenticator.Token) *AuthTokenMiddleware {
	return &AuthTokenMiddleware{
		acctToken: acctToken,
	}
}

func (a *AuthTokenMiddleware) RevokeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/login" || c.Request.URL.Path == "/api/user" {
			c.Next()
		} else {
			h := authHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				errModel := response.NewUnauthorizedError(err)
				c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("No valid token provided",
					nil, errModel))
				return
			}
			tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
			fmt.Println(tokenString)
			if tokenString == "" {
				errModel := response.NewUnauthorizedError(errors.New("Empty Bearer Token"))
				c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("No valid token provided",
					nil, errModel))
				return
			}
			token, err := a.acctToken.VerifyAccessToken(tokenString)
			numCode, err := a.acctToken.RevokeAccessToken(token)
			if numCode == -1 || err != nil {
				errModel := response.NewInternalServerError(err)
				c.AbortWithStatusJSON(errModel.StatusCode, response.NewResponse("Something Wrong",
					nil, errModel))
				return
			}
			fmt.Println(numCode, err)
		}
	}
}

func (a *AuthTokenMiddleware) RequireToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/api/user/login" {
			c.Next()
		} else {
			h := authHeader{}
			if err := c.ShouldBindHeader(&h); err != nil {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
			tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
			fmt.Println(tokenString)
			if tokenString == "" {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
			token, err := a.acctToken.VerifyAccessToken(tokenString)
			userId, err := a.acctToken.FetchAccessToken(token)
			if userId == "" || err != nil {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
			fmt.Println(token)
			if token != nil {
				c.Set("user-id", userId)
				c.Next()
			} else {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
		}
	}
}
