package webmaster

import (
	"github.com/astaxie/beego"
)

type ZonesController struct {
	beego.Controller
}

func (this *ZonesController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/zones.html"
}
