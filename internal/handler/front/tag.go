package front

import (
	"blog/internal/dao"
	"blog/internal/pkg"
	"github.com/gin-gonic/gin"
)

func Tag(c *gin.Context) {
	tags, _ := dao.GetTags()
	data := make(map[string]interface{})
	data["title"] = "标签"
	data["description"] = "蛋壳吧的博客标签"
	data["tags"] = tags
	pkg.Render(data, c, "tag")
}
