package webmaster

import (
	"ZKMS/models"
	"database/sql"
	"encoding/json"
	"github.com/astaxie/beego"
	"log"
	"strconv"
	"time"
)

type ReportsController struct {
	beego.Controller
}

func (this *ReportsController) Get() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Data["Email"] = "xianmaulin@gmail.com"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/reports.html"
}

// 显示某一个broker的详细性能数据
func (this *ReportsController) BrokerDetail() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/brokerdetail.html"

	// 获取参数
	zoneid := this.GetString("zoneid")
	brokerid := this.GetString("brokerid")

	this.Data["brokerid"] = brokerid

	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select * from `tb_zone` where `Id`=?", zoneid)
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
		// 将当前zoneid及其对应的地址集传到视图
		this.Data["zoneid"] = Id
		this.Data["hosts"] = Ip
	}

	// 获取broker的历史状态
	statdata := []models.BrokerStatModel{}
	rows, err = db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_broker_stat` where `BrokerId`=? and `ZoneId`=? order by `Timestamp` asc", brokerid, zoneid)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		tmp := models.BrokerStatModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		if err := rows.Scan(&tmp.Timestamp, &Cpu, &Net, &Disk); err != nil {
			log.Println(err)
			return
		}
		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		statdata = append(statdata, tmp)
	}

	// 将数据处理成json格式可用于直接渲染曲线图
	cpuData := [][]int64{}
	netData := [][]int64{}
	diskData := [][]int64{}
	for _, v := range statdata {
		timestamp, _ := time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy := []int64{timestamp.Unix() * 1000, int64(v.Cpu * 100)}
		cpuData = append(cpuData, xy)

		timestamp, _ = time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy = []int64{timestamp.Unix() * 1000, int64(v.Net * 100)}
		netData = append(netData, xy)

		timestamp, _ = time.Parse("2006-01-02 15:04:05", v.Timestamp)
		xy = []int64{timestamp.Unix() * 1000, int64(v.Disk * 100)}
		diskData = append(diskData, xy)
	}

	tmp, _ := json.Marshal(cpuData)
	this.Data["cpuData"] = string(tmp)

	tmp, _ = json.Marshal(netData)
	this.Data["netData"] = string(tmp)

	tmp, _ = json.Marshal(diskData)
	this.Data["diskData"] = string(tmp)
}

// 获取指定zone和broker的最新数据，用于实时数据展示
func (this *ReportsController) GetLatestBrokerData() {

	zoneid := this.GetString("zoneid")
	brokerid := this.GetString("brokerid")

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_broker_stat` where `BrokerId`=? and `ZoneId`=? and (`Timestamp`>=now() - interval 1 minute)", brokerid, zoneid)
	for rows.Next() {
		var Timestamp string
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		if err := rows.Scan(&Timestamp, &Cpu, &Net, &Disk); err != nil {
			log.Println(err)
			this.Ctx.WriteString("[]")
			return
		}
		timestamp, _ := time.Parse("2006-01-02 15:04:05", Timestamp)
		cpu, _ := strconv.ParseFloat(Cpu.String, 64)
		net, _ := strconv.ParseFloat(Net.String, 64)
		disk, _ := strconv.ParseFloat(Disk.String, 64)
		ret := "["
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64(cpu*100), 10) + "],"
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64(net*100), 10) + "],"
		ret += "[" + strconv.FormatInt(timestamp.Unix()*1000, 10) + "," + strconv.FormatInt(int64(disk*100), 10) + "]"
		ret += "]"
		this.Ctx.WriteString(ret)
	}
	this.Ctx.WriteString("[]")
	return
}
