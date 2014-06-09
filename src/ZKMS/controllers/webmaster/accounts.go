package webmaster

import (
	"ZKMS/models"
	"database/sql"
	"github.com/astaxie/beego"
	"log"
)

type AccountsController struct {
	beego.Controller
}

func (this *AccountsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/accounts.html"

	// 获取accounts列表
	accounts := []models.AccountModel{}
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select * from `tb_account`")
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
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
		accounts = append(accounts, account)
	}
	rows.Close()

	this.Data["accounts"] = accounts
}

func (this *AccountsController) Create() {
	createname := this.Input().Get("createname")
	createpassword := this.Input().Get("createpassword")
	createremark := this.Input().Get("createremark")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}
	_, err = db.Exec("insert into `tb_account`(`Name`,`Password`,`Remark`) values(?,?,?)", createname, createpassword, createremark)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}

	this.Ctx.WriteString("")
}

func (this *AccountsController) Edit() {
	name := this.Input().Get("editname")
	password := this.Input().Get("editpassword")
	remark := this.Input().Get("editremark")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}
	_, err = db.Exec("update `tb_account` set `Password`=?,`Remark`=? where `Name`=?", password, remark, name)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}

	this.Ctx.WriteString("")
}

func (this *AccountsController) Delete() {
	name := this.Input().Get("deletename")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}
	_, err = db.Exec("delete from `tb_account` where `Name`=?", name)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}

	this.Ctx.WriteString("")
}
