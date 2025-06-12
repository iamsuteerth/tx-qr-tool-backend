package request_id

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.GetRequestID(c)
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}
