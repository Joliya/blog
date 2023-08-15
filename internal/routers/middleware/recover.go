/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 22:41
 * @Description:
 */

package middleware

import (
	"blog/conf"
	"blog/pkg/logger"
	"blog/pkg/utils"
	"fmt"
	"github.com/convee/goblog/pkg/ding"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"runtime/debug"
)

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			r := recover()
			if utils.IsNotNil(r) {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Error("http_router_panic", zap.Any("err", r), zap.Stack(string(debug.Stack())))
				if !conf.Conf.App.DisableDingDing {
					_, _ = ding.SendAlert(fmt.Sprintf("http_router_panic:err:%v;stack:%s", r, string(debug.Stack())), false)
				}
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}()
	})
}

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				logger.Error("http_router_panic", zap.Any("err", r), zap.Stack(string(debug.Stack())))
				if !conf.Conf.App.DisableDingDing {
					_, _ = ding.SendAlert(fmt.Sprintf("http_router_panic:err:%v;stack:%s", r, string(debug.Stack())), false)
				}
				c.Error(err)

			}
		}()
		c.Next()
	}
}
