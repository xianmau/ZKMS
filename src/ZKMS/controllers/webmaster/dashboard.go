package webmaster

import (
	"github.com/astaxie/beego"
)

type DashboardController struct {
	beego.Controller
}

func (this *DashboardController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/dashboard.html"
}
