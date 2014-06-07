package webmaster

import (
	"github.com/astaxie/beego"
)

type AccountsController struct {
	beego.Controller
}

func (this *AccountsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/accounts.html"
}
