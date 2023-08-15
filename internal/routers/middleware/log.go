/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 22:38
 * @Description:
 */

package middleware

import (
	"blog/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(context *gin.Context) {
		t := time.Now()
		context.Next()
		logger.Info(fmt.Sprintf("[%d] %s %v", context.Status, context.Request.RequestURI, time.Since(t)))
	}
}
