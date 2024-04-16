package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wagslane/go-rabbitmq"

	"service-post-notification/internal/usecase/token_manager"
)

func init() {
	mustInitEnv()
}

func mustInitEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}

func MustInitTokenManagerConfig() token_manager.Config {
	var privateKey, publicKey string
	var ok bool

	if privateKey, ok = os.LookupEnv("TOKEN_PRIVATE_KEY"); !ok {
		panic("TOKEN_PRIVATE_KEY not set")
	}

	if publicKey, ok = os.LookupEnv("TOKEN_PUBLIC_KEY"); !ok {
		panic("TOKEN_PUBLIC_KEY not set")
	}

	return token_manager.Config{
		TtlAccessToken: token_manager.TtlAccessTokenDefault,
		PrivateKey:     privateKey,
		PublicKey:      publicKey,
	}
}

func MustInitRabbit() *rabbitmq.Conn {
	rabbit, err := rabbitmq.NewConn(
		fmt.Sprintf(
			"amqp://%s:%s@%s",
			os.Getenv("RM_USERNAME"),
			os.Getenv("RM_PASSWORD"),
			os.Getenv("RM_HOST"),
		),
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		panic(err)
	}

	return rabbit
}
