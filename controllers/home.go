package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type baseController struct {
	beego.Controller
	i18n.Locale
}

func (c *baseController) Prepare() {
	lang := c.GetString("lang")
	if lang == "zh-CN" {
		c.Lang = lang
	} else {
		c.Lang = "en-US"
	}
	c.Data["Lang"] = c.Lang
}

type HomeController struct {
	baseController
}

func (c *HomeController) Get() {
	var err error
	c.Data["IsHome"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Category"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"], err = models.GetAllTopics(c.Input().Get("cate"), c.Input().Get("label"), true)
	if err != nil {
		beego.Error(err)
	}

	c.Data["Hi"] = c.Tr("hi")
	c.Data["Bye"] = c.Tr("bye")
	c.Data["About"] = "about"
	c.TplName = "home.html"
}
