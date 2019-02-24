package routers

import (
	"beeblog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.HomeController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/category", &controllers.CategoryController{})
    beego.Router("/topic", &controllers.TopicController{})
    beego.Router("/reply", &controllers.ReplyController{})
    beego.Router("/reply/add", &controllers.ReplyController{}, "post:Add")
    beego.Router("/reply/delete", &controllers.ReplyController{}, "get:Delete")
    //自动路由
    beego.AutoRouter(&controllers.TopicController{})
}
