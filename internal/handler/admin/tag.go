package admin

import (
	"blog/internal/dao"
	"blog/internal/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
)

func TagList(c *gin.Context) {
	tags, err := dao.GetTags()
	if err != nil {
		fmt.Println("get ags err:", err)
		return
	}
	data := make(map[string]interface{})
	data["tags"] = tags
	pkg.AdminRender(data, c, "tag_list")
}
