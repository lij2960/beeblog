package controllers

import (
	"beeblog/models"
	"github.com/astaxie/beego"
)

type ReplyController struct {
	beego.Controller
}

func (c *ReplyController) Add() {
	tid := c.Input().Get("tid")
	err := models.AddReply(tid,
		c.Input().Get("nickname"),
		c.Input().Get("content"))
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/"+tid, 302)
}

func (c *ReplyController) Delete() {
	if !checkAccount(c.Ctx) {
		return
	}
	rid := c.Input().Get("rid")
	tid := c.Input().Get("tid")
	err := models.DeleteComment(rid)
	if err != nil {
		beego.Error(err)
	}
	c.Redirect("/topic/view/" + tid, 302)
}