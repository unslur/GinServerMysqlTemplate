// TemplateInterface
package main

import (
	"github.com/gin-gonic/gin"
)

type s struct {
	cryfrisdsfr string
	cryadadsf   int64
}

func queryLastPlaces(c *gin.Context) {

	defer func() {
		cylog.Flush()
		err := recover()
		if err != nil {
			cylog.Error("有错误发生，正在回滚", err)

		}

	}()
	cylog.Info("=========================查询GPS开始：%s %s", c.Request.URL.Path, c.ClientIP())
	defer cylog.Info("=========================查询GPS结束\n\n")

}
