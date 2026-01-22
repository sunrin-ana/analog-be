package interceptor

import (
	"os"
	"strings"

	"github.com/NARUBROWN/spine/core"
	"github.com/NARUBROWN/spine/interceptor/cors"
)

func NewCORSInterceptor() core.Interceptor {
	origins := getAllowedOrigins()

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
}

func getAllowedOrigins() []string {
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		return strings.Split(envOrigins, ",")
	}

	return []string{
		"http://localhost:3000",
		"http://localhost:8080",
		"https://ana.st",
		"https://*.ana.st",
	}
}
