package app

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"cbr/renderer"

	"golang.org/x/time/rate"
)

const (
	ratePerSecond = 5
	rateBurst     = 10
	visitorTTL    = 10 * time.Minute
	maxVisitors   = 4096
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimit returns a middleware allowing ratePerSecond requests per second
// with rateBurst headroom, tracked per client IP.
func RateLimit() func(http.Handler) http.Handler {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			mu.Lock()
			v, ok := visitors[ip]
			if !ok {
				v = &visitor{limiter: rate.NewLimiter(rate.Limit(ratePerSecond), rateBurst)}
				visitors[ip] = v
				if len(visitors) > maxVisitors {
					for ip, v := range visitors {
						if time.Since(v.lastSeen) > visitorTTL {
							delete(visitors, ip)
						}
					}
				}
			}
			v.lastSeen = time.Now()
			allowed := v.limiter.Allow()
			mu.Unlock()

			if !allowed {
				w.Header().Set("Retry-After", "1")
				renderer.JSON(w, http.StatusTooManyRequests, ErrorResponse{Error: "rate limit exceeded"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP trusts X-Forwarded-For only when the peer is loopback, i.e. a
// local reverse proxy that appended the client address; anything else could
// be spoofed by the client directly.
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if parsed := net.ParseIP(host); parsed == nil || !parsed.IsLoopback() {
		return host
	}

	xff := r.Header.Get("X-Forwarded-For")
	if i := strings.LastIndex(xff, ","); i >= 0 {
		xff = xff[i+1:]
	}
	xff = strings.TrimSpace(xff)
	if xff != "" && len(xff) <= net.IPv6len*2 {
		return xff
	}
	return host
}
