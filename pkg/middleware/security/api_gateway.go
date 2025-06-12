package security

import (
    "net/http"
    "os"
    
    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

func APIKeyAuthMiddleware() gin.HandlerFunc {
    apiKey := os.Getenv("API_GATEWAY_KEY")
    
    if apiKey == "" {
        log.Fatal().Msg("API_GATEWAY_KEY environment variable is not set")
    }
    
    return func(c *gin.Context) {
        if c.Request.URL.Path == "/health" {
            c.Next()
            return
        }
        
        requestKey := c.GetHeader("X-Api-Key")
        if requestKey == "" {
            log.Warn().
                Str("path", c.Request.URL.Path).
                Str("method", c.Request.Method).
                Msg("Missing API key")
            
            c.JSON(http.StatusUnauthorized, gin.H{
                "status":  "ERROR",
                "code":    "MISSING_API_KEY",
                "message": "API key is required",
            })
            c.Abort()
            return
        }
        
        if requestKey != apiKey {
            log.Warn().
                Str("path", c.Request.URL.Path).
                Str("method", c.Request.Method).
                Msg("Invalid API key")
            
            c.JSON(http.StatusForbidden, gin.H{
                "status":  "ERROR",
                "code":    "INVALID_API_KEY",
                "message": "Invalid API key",
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
