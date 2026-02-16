package middleware

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
	"go-industry-server/pkg/response"
)

type client struct {
	limiter *rate.Limiter
}

var (
	mu      sync.Mutex
	clients = make(map[string]*client)
)

func RateLimit(r float64, b int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				ip = req.RemoteAddr
			}

			mu.Lock()
			if _, ok := clients[ip]; !ok {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(r), b),
				}
			}
			lim := clients[ip].limiter
			mu.Unlock()

			if !lim.Allow() {
				w.Header().Set("Retry-After", "1")
				response.Error(w, http.StatusTooManyRequests, "rate limit exceeded")
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}

