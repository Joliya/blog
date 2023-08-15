/**
 * @Author: jinpeng zhang
 * @Date: 2023/8/13 22:50
 * @Description:
 */

package admin

import (
	"blog/internal/dao"
	"blog/internal/model"
	"blog/internal/pkg"
	"blog/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func Login(c *gin.Context) {
	data := make(map[string]interface{})
	pkg.AdminRender(data, c, "index")
}

func Register(c *gin.Context) {
	data := make(map[string]interface{})
	pkg.AdminRender(data, c, "register")
}

func Logout(c *gin.Context) {
	c.SetCookie("email", "", -1, "/", "", false, false)
	c.Redirect(http.StatusFound, "/login")
}

func Signup(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	repassword := c.PostForm("repassword")
	if password != repassword {
		c.Redirect(http.StatusFound, "/admin")
	}
	secretPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if utils.IsNotNil(err) {
		c.Redirect(http.StatusFound, "/admin")
	}
	user := model.User{Email: email, Password: string(secretPwd)}
	_, err = dao.AddUser(user)
	if utils.IsNotNil(err) {
		c.Redirect(http.StatusFound, "/admin")
	}
	c.Redirect(http.StatusFound, "/")
}

func Signin(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		c.Status(http.StatusInternalServerError)
		return
	}
	user := dao.GetUser(email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		data := make(map[string]interface{})
		data["msg"] = "密码不正确，请重试"
		pkg.AdminRender(data, c, "401")
		return
	}
	c.SetCookie("email", email, -1, "/", "", true, true)
	c.Redirect(http.StatusFound, "/admin")
	return
}
