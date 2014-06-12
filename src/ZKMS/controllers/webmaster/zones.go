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

	zonelist := []models.ZoneModel{}

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
		zonelist = append(zonelist, zone)
	}
	rows.Close()
	this.Data["zonelist"] = zonelist

	zone := this.GetString("zone")
	if zone == "" {
		zone = zonelist[0].Id
	}

	rows, err = db.Query("select * from `tb_zone` where `Id`=?", zone)
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
		this.Data["cur_zoneid"] = Id
		this.Data["cur_zone"] = Ip
		return
	}
	this.Data["cur_zoneid"] = "Unknown Zone"
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

	retjson := `{"data":` + string(data) + `,"stat":` + string(stat_json) + `}`
	this.Ctx.WriteString(retjson)
	//log.Println(retjson)
}

func (this *ZonesController) CreateNode() {
	zone := this.GetString("zone")
	hosts := []string{}
	err := json.Unmarshal([]byte(zone), &hosts)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	path := this.GetString("path")
	data := this.GetString("data")

	conn, _, err := zk.Connect(hosts, time.Second)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	if flag, _, err := conn.Exists(path); err == nil && !flag {
		conn.Create(path, []byte(data), 0, zk.WorldACL(zk.PermAll))
	} else if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	this.Ctx.WriteString("ok")
}

func (this *ZonesController) UpdateNode() {
	zone := this.GetString("zone")
	hosts := []string{}
	err := json.Unmarshal([]byte(zone), &hosts)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	path := this.GetString("path")
	data := this.GetString("data")

	conn, _, err := zk.Connect(hosts, time.Second)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	if flag, _, err := conn.Exists(path); err == nil && flag {
		conn.Set(path, []byte(data), -1)
	} else if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	this.Ctx.WriteString("ok")
}

func (this *ZonesController) DeleteNode() {
	zone := this.GetString("zone")
	hosts := []string{}
	err := json.Unmarshal([]byte(zone), &hosts)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	path := this.GetString("path")

	conn, _, err := zk.Connect(hosts, time.Second)
	if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	if flag, _, err := conn.Exists(path); err == nil && flag {
		if err := conn.Delete(path, -1); err != nil {
			log.Println(err)
			this.Ctx.WriteString(err.Error())
			return
		}
	} else if err != nil {
		log.Println(err)
		this.Ctx.WriteString(err.Error())
		return
	}
	this.Ctx.WriteString("ok ")
}
