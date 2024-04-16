package token_manager

import (
	"fmt"
	"time"
)

const (
	TtlAccessTokenDefault = time.Hour * 24
)

type TokenManager struct {
	config     Config
	jwtManager jwtManager
}

func New(tokenCreator jwtManager, config Config) *TokenManager {
	return &TokenManager{
		jwtManager: tokenCreator,
		config:     config,
	}
}

func (m *TokenManager) CreateAuthToken(userID string) (string, error) {
	accessToken, err := m.jwtManager.CreateToken(
		m.config.PrivateKey,
		m.config.TtlAccessToken,
		Data{
			UserID: userID,
		})
	if err != nil {
		return "", fmt.Errorf("failed create access token: %w", err)
	}

	return accessToken, nil
}

func (m *TokenManager) ExtractUserID(token string) (string, error) {
	payload, err := m.jwtManager.ExtractTokenData(
		m.config.PublicKey,
		token,
	)
	if err != nil {
		return "", fmt.Errorf("failed extract token data: %w", err)
	}

	return payload.(map[string]interface{})["user_id"].(string), nil
}
