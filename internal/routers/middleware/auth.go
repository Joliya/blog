/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 22:36
 * @Description:
 */

package middleware

import (
	"blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthWithCookie() gin.HandlerFunc {
	return func(context *gin.Context) {
		if cookie, err := context.Request.Cookie("email"); utils.IsNotNil(err) || cookie.Value == "" {
			context.Redirect(http.StatusFound, "/login")
			return
		}
		context.Next()
	}
}
