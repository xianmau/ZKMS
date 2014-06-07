package webmaster

import (
	"github.com/astaxie/beego"
)

type ReportsController struct {
	beego.Controller
}

func (this *ReportsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/reports.html"
}
