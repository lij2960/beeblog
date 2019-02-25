package main

import (
	"beeblog/models"
	_ "beeblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"os"
)

func init()  {
	models.RegisterDB()
}

func main() {
	orm.Debug = true
	orm.RunSyncdb("default", false, true)
	/**
	此处与视频不一致，routers放入routers目录
	 */
	//创建附件目录
	os.Mkdir("attachment", os.ModePerm)
	beego.Run()
}

