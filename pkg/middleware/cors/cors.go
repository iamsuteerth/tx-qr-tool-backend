package cors

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupCORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		allowedOrigins := getAllowedOrigins()

		if isOriginAllowed(origin, allowedOrigins) {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, X-Api-Key, X-Request-ID")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

func getAllowedOrigins() []string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		return []string{"*"}
	}
	return strings.Split(origins, ",")
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}
