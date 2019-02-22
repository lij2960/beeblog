package routers

import (
	"beeblog/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.HomeController{})
    beego.Router("/login", &controllers.LoginController{})
    beego.Router("/category", &controllers.CategoryController{})
}
