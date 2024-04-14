package token_manager

import "time"

type Config struct {
	TtlAccessToken time.Duration
	PrivateKey     string
	PublicKey      string
}

type Data struct {
	UserID string `json:"user_id"`
}
