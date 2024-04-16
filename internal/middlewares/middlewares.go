package middlewares

type Middlewares struct {
	tokenManager tokenManager
}

func New(tokenManager tokenManager) *Middlewares {
	return &Middlewares{
		tokenManager: tokenManager,
	}
}
