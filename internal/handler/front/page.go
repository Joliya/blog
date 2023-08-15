package front

import (
	"blog/internal/dao"
	"blog/internal/pkg"
	"github.com/gin-gonic/gin"
)

func Page(c *gin.Context) {
	id := c.Query("id")
	page := dao.GetPage(id)
	data := make(map[string]interface{})
	data["title"] = page.Title
	data["description"] = page.Title
	data["page"] = page
	pkg.Render(data, c, "page")
}
