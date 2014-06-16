package routers

import (
	//"ZKMS/controllers"
	"ZKMS/controllers/webmaster"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &webmaster.DashboardController{})

	beego.Router("/webmaster", &webmaster.DashboardController{})
	beego.Router("/webmaster/dashboard/SetExtractor", &webmaster.DashboardController{}, "post:SetExtractor")
	beego.Router("/webmaster/dashboard/ExecuteSql", &webmaster.DashboardController{}, "post:ExecuteSql")

	// 管理报表
	beego.Router("/webmaster/reports", &webmaster.ReportsController{})
	beego.Router("/webmaster/reports/brokerdetail", &webmaster.ReportsController{}, "get:BrokerDetail")
	beego.Router("/webmaster/reports/brokerdetail/getlatestdata", &webmaster.ReportsController{}, "get:GetLatestBrokerData")
	beego.Router("/webmaster/reports/loggerdetail", &webmaster.ReportsController{}, "get:LoggerDetail")
	beego.Router("/webmaster/reports/loggerdetail/getlatestdata", &webmaster.ReportsController{}, "get:GetLatestLoggerData")

	beego.Router("/webmaster/reports/brokers", &webmaster.ReportsController{}, "get:Brokers")
	beego.Router("/webmaster/reports/loggers", &webmaster.ReportsController{}, "get:Loggers")
	beego.Router("/webmaster/reports/apps", &webmaster.ReportsController{}, "get:Apps")
	beego.Router("/webmaster/reports/topics", &webmaster.ReportsController{}, "get:Topics")

	// 管理zones
	beego.Router("/webmaster/zones", &webmaster.ZonesController{})
	beego.Router("/webmaster/zones/CreateNode", &webmaster.ZonesController{}, "post:CreateNode")
	beego.Router("/webmaster/zones/UpdateNode", &webmaster.ZonesController{}, "post:UpdateNode")
	beego.Router("/webmaster/zones/DeleteNode", &webmaster.ZonesController{}, "post:DeleteNode")
	beego.Router("/webmaster/zktree", &webmaster.ZonesController{}, "get:ZkTree")
	beego.Router("/webmaster/zktree/getdata", &webmaster.ZonesController{}, "get:GetData")
	beego.Router("/webmaster/zktree/getchildren", &webmaster.ZonesController{}, "get:GetChildren")

	// 帐户管理
	beego.Router("/webmaster/accounts", &webmaster.AccountsController{})
	beego.Router("/webmaster/accounts/Create", &webmaster.AccountsController{}, "post:Create")
	beego.Router("/webmaster/accounts/Edit", &webmaster.AccountsController{}, "post:Edit")
	beego.Router("/webmaster/accounts/Delete", &webmaster.AccountsController{}, "post:Delete")

	// 设置管理，包括
	beego.Router("/webmaster/settings", &webmaster.SettingsController{})
	beego.Router("/webmaster/settings/CreateZone", &webmaster.SettingsController{}, "post:CreateZone")
	beego.Router("/webmaster/settings/EditZone", &webmaster.SettingsController{}, "post:EditZone")
	beego.Router("/webmaster/settings/DeleteZone", &webmaster.SettingsController{}, "post:DeleteZone")
}
