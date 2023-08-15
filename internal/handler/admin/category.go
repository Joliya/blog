package admin

import (
	"blog/internal/dao"
	"blog/internal/model"
	"blog/internal/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CategoryList(c *gin.Context) {
	categories, err := dao.GetCategories()
	if err != nil {
		fmt.Println("get categories err:", err)
		return
	}
	data := make(map[string]interface{})
	data["categories"] = categories
	pkg.AdminRender(data, c, "category_list")
}

func CategoryAdd(c *gin.Context) {
	data := make(map[string]interface{})
	id, _ := strconv.Atoi(c.PostForm("id"))
	var category model.Category
	if id > 0 {
		category = dao.GetCategory(id)
	}
	categories, _ := dao.GetCategories()
	data["categories"] = categories

	if category.Id > 0 {

		data["id"] = category.Id
		data["name"] = category.Name
	}
	pkg.AdminRender(data, c, "category_add")
}

func CategoryDelete(c *gin.Context) {
	var category model.Category
	category.Id, _ = strconv.Atoi(c.Query("id"))
	_, err := dao.DeleteCategory(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "删除失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}

func CategorySave(c *gin.Context) {
	var category model.Category
	category.Id, _ = strconv.Atoi(c.PostForm("id"))
	category.Name = c.PostForm("name")
	_, err := dao.SaveCategory(category)
	if err != nil {
		data := make(map[string]interface{})
		data["msg"] = "添加或修改失败，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	c.Redirect(http.StatusFound, "/admin/category")
}
