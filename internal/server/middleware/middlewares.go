package middleware

import "github.com/shanto-323/backend-scaffold/internal/server"

type Middlewares struct {
	*Global
	*RateLimit
	*ContextEnhancer
	*Tracer
}

func New(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobal(s),
		RateLimit:       NewRateLimit(s),
		ContextEnhancer: NewContextEnhancer(s),
		Tracer:          NewTracer(s),
	}
}
