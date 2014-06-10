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
	beego.Router("/webmaster/zktree", &webmaster.ZonesController{}, "get:ZkTree")
	beego.Router("/webmaster/zktree/getdata", &webmaster.ZonesController{}, "get:GetData")
	beego.Router("/webmaster/zktree/getchildren", &webmaster.ZonesController{}, "get:GetChildren")

	// 设置帐户
	beego.Router("/webmaster/accounts", &webmaster.AccountsController{})
	beego.Router("/webmaster/accounts/Create", &webmaster.AccountsController{}, "post:Create")
	beego.Router("/webmaster/accounts/Edit", &webmaster.AccountsController{}, "post:Edit")
	beego.Router("/webmaster/accounts/Delete", &webmaster.AccountsController{}, "post:Delete")

	//
	//beego.Router("/webmaster/accounts/create", &webmaster.AccountsController{}, "get:CreateAccount")

	// 设置管理，包括
	beego.Router("/webmaster/settings", &webmaster.SettingsController{})
	beego.Router("/webmaster/settings/CreateZone", &webmaster.SettingsController{}, "post:CreateZone")
	beego.Router("/webmaster/settings/EditZone", &webmaster.SettingsController{}, "post:EditZone")
	beego.Router("/webmaster/settings/DeleteZone", &webmaster.SettingsController{}, "post:DeleteZone")
}
