package webmaster

import (
	"github.com/astaxie/beego"
)

type SettingsController struct {
	beego.Controller
}

func (this *SettingsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/settings.html"
}
