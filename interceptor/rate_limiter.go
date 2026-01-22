package interceptor

import (
	"analog-be/pkg"
	"strings"
	"sync"
	"time"

	"github.com/NARUBROWN/spine/core"
	"golang.org/x/time/rate"
)

type RateLimitInterceptor struct {
	visitors map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

func NewRateLimitInterceptor() *RateLimitInterceptor {
	rps := 10
	burst := 20

	rl := &RateLimitInterceptor{
		visitors: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}

	// 3분마다
	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimitInterceptor) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.mu.Lock()
		rl.visitors[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimitInterceptor) cleanupVisitors() {
	for {
		time.Sleep(3 * time.Minute)
		rl.mu.Lock()
		rl.visitors = make(map[string]*rate.Limiter)
		rl.mu.Unlock()
	}
}

func (rl *RateLimitInterceptor) PreHandle(ctx core.ExecutionContext, _ core.HandlerMeta) error {
	ip := "unknown"
	if reqCtx := ctx.Context(); reqCtx != nil {
		if ipVal := reqCtx.Value("remote_addr"); ipVal != nil {
			if ipStr, ok := ipVal.(string); ok {
				ip = ipStr
			}
		}
	}

	if ip == "unknown" {
		if forwarded := ctx.Header("X-Forwarded-For"); forwarded != "" {
			ips := strings.Split(forwarded, ",")
			if len(ips) > 0 {
				ip = strings.TrimSpace(ips[0])
			}
		} else if realIP := ctx.Header("X-Real-IP"); realIP != "" {
			ip = realIP
		}
	}

	limiter := rl.getVisitor(ip)

	if !limiter.Allow() {
		if rwAny, ok := ctx.Get("spine.response_writer"); ok {
			if rw, ok := rwAny.(core.ResponseWriter); ok {
				rw.WriteJSON(429, map[string]string{
					"error":   "RATE_LIMIT_EXCEEDED",
					"message": "Rate limit exceeded",
				})
				return core.ErrAbortPipeline
			}
		}
		return pkg.NewAppError(429, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded")
	}

	return nil
}

func (rl *RateLimitInterceptor) PostHandle(core.ExecutionContext, core.HandlerMeta) {}

func (rl *RateLimitInterceptor) AfterCompletion(core.ExecutionContext, core.HandlerMeta, error) {}
