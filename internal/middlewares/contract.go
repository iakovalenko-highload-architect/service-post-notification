package middlewares

type tokenManager interface {
	ExtractUserID(token string) (string, error)
}
