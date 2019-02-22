package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	var err error
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Topics"], err = models.GetAllTopics(false)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "topic.html"
}

func (c *TopicController) Add() {
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.TplName = "topic_add.html"
}

func (c *TopicController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	title := c.Input().Get("title")
	content := c.Input().Get("content")
	id := c.Input().Get("id")

	var err error
	if len(id) == 0 {
		err = models.AddTopic(title, content)
	} else {
		err = models.ModifyTopic(id, title, content)
	}
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
}

func (c *TopicController) View() {
	topic, err := models.GetTopic(c.Ctx.Input.Param("0"))
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Id"] = c.Ctx.Input.Param("0")
	c.TplName = "topic_view.html"
}

func (c *TopicController) Modify() {
	id := c.Input().Get("id")
	topic, err := models.GetTopic(id)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Id"] = id
	c.TplName = "topic_modify.html"
}

func (c *TopicController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/", 302)
		return
	}
	err := models.DeleteTopic(c.Input().Get("id"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic", 302)
}