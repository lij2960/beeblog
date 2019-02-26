package main

import (
	"beeblog/models"
	_ "beeblog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
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
	i18n.SetMessage("en-US", "conf/locale_en-US.ini")
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")
	beego.AddFuncMap("i18n", i18n.Tr)
	beego.Run()
}

