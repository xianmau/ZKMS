package webmaster

import (
	"ZKMS/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
	"time"
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

func (this *ZonesController) ZkTree() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/zktree.html"

	zone := this.GetString("zone")

	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select * from `tb_zone` where `Id`=?", zone)
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
		this.Data["cur_zone"] = Ip
		return
	}
	this.Data["cur_zone"] = "[]"
}

func (this *ZonesController) GetChildren() {
	currentzone := this.GetString("zoneid")
	curzone := []string{}
	err := json.Unmarshal([]byte(currentzone), &curzone)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	currentnode := "/" + this.GetString("znode")
	data := []models.ZkTreeNode{}
	conn, _, err := zk.Connect(curzone, time.Second)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	if flag, _, err := conn.Exists(currentnode); err == nil && flag {
		children, _, err := conn.Children(currentnode)
		if err != nil {
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
		for _, znode := range children {
			d := models.ZkTreeNode{znode, true}
			data = append(data, d)
		}
	}
	datajson, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	this.Ctx.WriteString(string(datajson))
}

func (this *ZonesController) GetData() {
	currentzone := this.GetString("zoneid")
	curzone := []string{}
	err := json.Unmarshal([]byte(currentzone), &curzone)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	currentnode := this.GetString("znode")
	conn, _, err := zk.Connect(curzone, time.Second)
	defer conn.Close()
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	data, stat, err := conn.Get(currentnode)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	if data == nil || len(data) <= 0 || data[0] == 0 {
		this.Ctx.WriteString("[]")
		return
	}
	stat_json, err := json.Marshal(stat)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	//fmt.Printf("Stat -> %-v\n", stat)
	this.Ctx.WriteString(`{"data":` + string(data) + `,"stat":` + string(stat_json) + `}`)
}
