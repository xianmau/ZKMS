package webmaster

import (
	"ZKMS/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"log"
	"strings"
)

type ZonesController struct {
	beego.Controller
}

func (this *ZonesController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/zones.html"

	// 获取zones列表
	zones := []models.ZoneModel{}
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
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
		var zone models.ZoneModel
		var Id string
		var Ip string
		if err := rows.Scan(&Id, &Ip); err != nil {
			log.Println(err)
			return
		}
		zone.Id = Id
		if err := json.Unmarshal([]byte(Ip), &zone.Ip); err != nil {
			log.Println(err)
			return
		}
		n := len(zone.Ip)
		for i := 0; i < n; i++ {
			if !strings.ContainsRune(zone.Ip[i], rune(':')) {
				zone.Ip[i] = fmt.Sprintf("%s:%d", zone.Ip[i], 2181)
			}
		}
		if n < 5 {
			for ; n < 5; n++ {
				zone.Ip = append(zone.Ip, "")
			}
		}
		zones = append(zones, zone)

	}
	rows.Close()

	this.Data["zones"] = zones
}
