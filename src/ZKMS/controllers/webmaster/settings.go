package webmaster

import (
	"database/sql"
	"github.com/astaxie/beego"
	"log"
)

type SettingsController struct {
	beego.Controller
}

func (this *SettingsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/settings.html"

	// 获取zones列表
	zones := make(map[string]string)
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select * from `tb_zone`")
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var Id string
		var Ip string
		if err := rows.Scan(&Id, &Ip); err != nil {
			log.Println(err)
			return
		}
		zones[Id] = Ip
	}
	rows.Close()

	this.Data["zones"] = zones
}

func (this *SettingsController) CreateZone() {
	zoneid := this.Input().Get("createzoneid")
	hosts := this.Input().Get("createhosts")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	_, err = db.Exec("insert into `tb_zone`(`Id`,`Ip`) values(?,?)", zoneid, hosts)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}

	this.Ctx.WriteString("")
}

func (this *SettingsController) EditZone() {
	zoneid := this.Input().Get("editzoneid")
	hosts := this.Input().Get("edithosts")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	_, err = db.Exec("update `tb_zone` set `Ip`=? where `Id`=?", hosts, zoneid)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
	}

	this.Ctx.WriteString("")
}

func (this *SettingsController) DeleteZone() {
	zoneid := this.Input().Get("deletezoneid")
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	_, err = db.Exec("delete from `tb_zone` where `Id`=?", zoneid)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}

	this.Ctx.WriteString("")
}
