package webmaster

import (
	"github.com/astaxie/beego"
)

type WebmasterController struct {
	beego.Controller
}

func (this *WebmasterController) Get() {
	this.Data["Website"] = "beego.me"
	this.Data["Email"] = "astaxie@gmail.com"
	this.TplNames = "webmaster/dashboard.tpl"
}
