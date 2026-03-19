package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func (m *Middlewares) RateLimit(next http.Handler) http.Handler {
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	if m.Config.Limiter.Enabled {
		go clearClients(&mu, clients)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.Config.Limiter.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		ip := realip.FromRequest(r)

		mu.Lock()

		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{
				limiter: rate.NewLimiter(2, 4),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			m.ErrorManager.RateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func clearClients(mu *sync.Mutex, clients map[string]*client) {
	for {
		time.Sleep(time.Minute)

		mu.Lock()

		for ip, client := range clients {
			if time.Since(client.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}

		mu.Unlock()
	}
}
