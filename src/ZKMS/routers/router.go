package routers

import (
	//"ZKMS/controllers"
	"ZKMS/controllers/webmaster"
	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})

	beego.Router("/webmaster", &webmaster.DashboardController{})

	beego.Router("/webmaster/reports", &webmaster.ReportsController{})

	beego.Router("/webmaster/zones", &webmaster.ZonesController{})

	// 设置帐户
	beego.Router("/webmaster/accounts", &webmaster.AccountsController{})

	//
	//beego.Router("/webmaster/accounts/create", &webmaster.AccountsController{}, "get:CreateAccount")

	// 设置管理，包括
	beego.Router("/webmaster/settings", &webmaster.SettingsController{})
	// 设置主机，也就是管理数据库里的zones
	//beego.Router("/webmaster/settings/hosts", &webmaster.SettingsController{}, "get:SetHost")
}
