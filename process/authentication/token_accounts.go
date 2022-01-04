package authentication

import (
	"errors"
	"fmt"
	"log"
	"time"

	repository "integration/process/repository"

	model "integration/process/authentication/model"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Token interface {
	CreateAccessToken(cred *model.CredentialModel) (*TokenDetails, error)
	VerifyAccessToken(tokenString string) (*AccessDetails, error)
	StoreAccessToken(userName string, tokenDetail *TokenDetails) error
	FetchAccessToken(accessDetails *AccessDetails) (string, error)
	RevokeAccessToken(accessDetails *AccessDetails) (int64, error)
}

type TokenConfig struct {
	ApplicationName     string
	JwtSignatureKey     string
	JwtSigningMethod    *jwt.SigningMethodHMAC
	AccessTokenLifeTime time.Duration
	Client              repository.ICacheRepository
}
type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	AtExpires   int64
}

type AccessDetails struct {
	AccessUuid string
	UserName   string
}
type token struct {
	Config TokenConfig
}

func NewTokenService(config TokenConfig) Token {
	return &token{
		Config: config,
	}
}

func (t *token) CreateAccessToken(cred *model.CredentialModel) (*TokenDetails, error) {
	td := &TokenDetails{}
	now := time.Now().UTC()
	end := now.Add(t.Config.AccessTokenLifeTime)

	td.AtExpires = end.Unix()
	td.AccessUuid = uuid.New().String()

	claims := model.MyClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer: t.Config.ApplicationName,
		},
		Username:   cred.Username,
		Email:      cred.Email,
		AccessUUID: td.AccessUuid,
	}
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = end.Unix()
	token := jwt.NewWithClaims(
		t.Config.JwtSigningMethod,
		claims,
	)
	newToken, err := token.SignedString([]byte(t.Config.JwtSignatureKey))
	td.AccessToken = newToken
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (t *token) StoreAccessToken(userName string, tokenDetail *TokenDetails) error {
	at := time.Unix(tokenDetail.AtExpires, 0)
	now := time.Now()
	err := t.Config.Client.Set(tokenDetail.AccessUuid, userName, at.Sub(now))
	if err != nil {
		return err
	}
	err = t.Config.Client.Set(userName, 1, at.Sub(now))
	if err != nil {
		return err
	}
	return nil
}
func (t *token) FetchAccessToken(accessDetails *AccessDetails) (string, error) {
	if accessDetails != nil {
		userId, err := t.Config.Client.Get(accessDetails.AccessUuid)
		if err != nil {
			return "", err
		}
		return userId, nil
	} else {
		return "", errors.New("Invalid Access")
	}

}

func (t *token) RevokeAccessToken(accessDetails *AccessDetails) (int64, error) {
	if accessDetails != nil {
		number, err := t.Config.Client.Delete(accessDetails.AccessUuid)
		if err != nil {
			return -1, err
		}
		return number, nil
	} else {
		return -1, errors.New("Invalid Access")
	}

}
func (t *token) VerifyAccessToken(tokenString string) (*AccessDetails, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Signing method invalid")
		} else if method != t.Config.JwtSigningMethod {
			return nil, fmt.Errorf("Signing method invalid")
		}

		return []byte(t.Config.JwtSignatureKey), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims["iss"] != t.Config.ApplicationName {
		log.Println("Token Invalid")
		return nil, err
	}
	accessUUID := claims["AccessUUID"].(string)
	userName := claims["Username"].(string)
	return &AccessDetails{
		AccessUuid: accessUUID,
		UserName:   userName,
	}, nil
}
