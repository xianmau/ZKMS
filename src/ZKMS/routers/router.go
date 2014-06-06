package routers

import (
	"ZKMS/controllers"
	"ZKMS/controllers/webmaster"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("/webmaster", &webmaster.WebmasterController{})

}
