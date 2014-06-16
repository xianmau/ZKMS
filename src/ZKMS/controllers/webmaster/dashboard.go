package webmaster

import (
	"ZKMS/controllers/dbmaster"
	"ZKMS/controllers/extractor"
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

	if extractor.EXTRACTORSTATUS == extractor.EXTRACTSTATUS_RUNNING {
		this.Data["ExtractorStatus"] = "Running"
	} else if extractor.EXTRACTORSTATUS == extractor.EXTRACTSTATUS_STOPPED {
		this.Data["ExtractorStatus"] = "Stopped"
	} else {
		this.Data["ExtractorStatus"] = "Unknown"
	}
}

func (this *DashboardController) SetExtractor() {
	act := this.Input().Get("act")
	if act == "run" {
		if err := extractor.RunExtractor(); err != nil {
			this.Ctx.WriteString(err.Error())
			return
		}
	} else if act == "stop" {
		if err := extractor.StopExtractor(); err != nil {
			this.Ctx.WriteString(err.Error())
			return
		}
	}
	this.Ctx.WriteString("")
}

func (this *DashboardController) ExecuteSql() {
	sqlstr := this.Input().Get("sqlstr")
	if n, err := dbmaster.ExecuteSql(sqlstr); err != nil {
		this.Ctx.WriteString(err.Error())
		return
	} else {
		this.Ctx.WriteString(string(n))
		return
	}
	this.Ctx.WriteString("")
}
