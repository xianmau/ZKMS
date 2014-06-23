package webmaster

import (
	"ZKMS/models"
	"database/sql"
	"github.com/astaxie/beego"
	"log"
)

type LoginController struct {
	beego.Controller
}

func (this *LoginController) Get() {
	this.TplNames = "webmaster/login.html"
}

func (this *LoginController) Post() {
	name := this.Input().Get("Name")
	password := this.Input().Get("Password")
	//remenber := this.Input.Get("Remenber")

	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select * from `tb_account` where Name=? and Password=?", name, password)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return
	}
	if rows.Next() {
		var account models.AccountModel
		var Name string
		var Password string
		var Remark string
		if err := rows.Scan(&Name, &Password, &Remark); err != nil {
			log.Println(err)
			return
		}
		account.Name = Name
		account.Password = Password
		account.Remark = Remark

		this.SetSession("admin", account)
		this.Ctx.Redirect(302, "/webmaster")
	} else {
		this.Ctx.WriteString("failed")
	}
}

func (this *LoginController) Logout() {
	this.DelSession("admin")
	this.Ctx.Redirect(302, "/webmaster/login")
}
