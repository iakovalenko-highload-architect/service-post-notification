package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt"
)

type JWT struct {
}

func New() *JWT {
	return &JWT{}
}

func (j *JWT) CreateToken(privateKey string, ttl time.Duration, payload interface{}) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed decode base64 private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "failed parse private key")
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["sub"] = payload
	claims["exp"] = now.Add(ttl).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "failed create token")
	}

	return token, nil
}

func (j *JWT) ExtractTokenData(publicKey string, token string) (interface{}, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed decode base64 public key")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed parse public key")
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			return nil, errors.Wrapf(err, "token %s was expired", token)
		default:
			return nil, errors.Wrapf(err, "invalid refresh token  %s", token)
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.Wrap(err, "failed extract token data")
	}

	return claims["sub"], nil
}
