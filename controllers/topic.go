package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
	"path"
	"strings"
)

type TopicController struct {
	beego.Controller
}

func (c *TopicController) Get() {
	var err error
	c.Data["IsTopic"] = true
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Topics"], err = models.GetAllTopics("", "", false)
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
	category := c.Input().Get("category")
	label := c.Input().Get("label")

	//获取文件
	_, fh, err := c.GetFile("attachment")
	if err != nil {
		beego.Error(err)
	}
	var attachment string
	if fh != nil {
		attachment = fh.Filename
		beego.Info(attachment)
		err = c.SaveToFile("attachment", path.Join("attachment", attachment))
		if err != nil {
			beego.Error(err)
		}
	}

	if len(id) == 0 {
		err = models.AddTopic(title, category, label, content, attachment)
	} else {
		err = models.ModifyTopic(id, title, category, label, content, attachment)
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
	tid := c.Ctx.Input.Param("0")
	comments, err := models.GetAllComments(tid)
	if err != nil {
		beego.Error(err)
		return
	}
	c.Data["Topic"] = topic
	c.Data["Id"] = tid
	c.Data["Comments"] = comments
	c.Data["IsLogin"] = checkAccount(c.Ctx)
	c.Data["Labels"] = strings.Split(topic.Labels, " ")
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