package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	var err error
	c.Data["IsHome"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Category"], err = models.GetAllCategories()
	if err != nil {
		beego.Error(err)
	}
	c.Data["Topics"], err = models.GetAllTopics(c.Input().Get("cate"), true)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "home.html"
}
