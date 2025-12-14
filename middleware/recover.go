package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/karlorz/auth-gateway/common"
	"github.com/gin-gonic/gin"
)

func RelayPanicRecover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				common.SysLog(fmt.Sprintf("panic detected: %v", err))
				common.SysLog(fmt.Sprintf("stacktrace from panic: %s", string(debug.Stack())))
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"message": fmt.Sprintf("Panic detected, error: %v. Please submit an issue here: https://github.com/karlorz/auth-gateway", err),
						"type":    "auth_gateway_panic",
					},
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
