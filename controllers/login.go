package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	isExits := c.Input().Get("exit") == "true"
	if isExits {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("pwd", "", -1, "/")
		//视频讲解是301，不过测试有问题，需要设置为302
		c.Redirect("/", 302)
		return
	}
	c.TplName = "login.html"
}

func (c *LoginController) Post()  {
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLogin") == "on"

	if beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd{
		maxAge := 0
		if autoLogin {
			maxAge = 1 << 31 -1
		}
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")
	}

	c.Redirect("/", 302)
	return
}

/**
此处与视频不一致，原因参见博文：
https://blog.php2.cc/beego-context%E6%96%B0%E8%80%81%E7%89%88%E6%9C%AC%E4%BD%BF%E7%94%A8-89.html
 */
func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value
	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}
	pwd := ck.Value
	return beego.AppConfig.String("uname") == uname &&
		beego.AppConfig.String("pwd") == pwd
}