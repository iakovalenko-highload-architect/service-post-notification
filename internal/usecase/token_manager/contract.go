package token_manager

import "time"

type jwtManager interface {
	CreateToken(privateKey string, ttl time.Duration, payload interface{}) (string, error)
	ExtractTokenData(publicKey string, token string) (interface{}, error)
}
