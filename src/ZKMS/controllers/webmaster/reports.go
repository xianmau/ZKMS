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

// 显示某一个logger的详细性能数据
func (this *ReportsController) LoggerDetail() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/loggerdetail.html"

	// 获取参数
	zoneid := this.GetString("zoneid")
	loggerid := this.GetString("loggerid")

	this.Data["loggerid"] = loggerid

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

	// 获取logger的历史状态
	statdata := []models.LoggerStatModel{}
	rows, err = db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_logger_stat` where `LoggerId`=? and `ZoneId`=? order by `Timestamp` asc", loggerid, zoneid)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		tmp := models.LoggerStatModel{}
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

// 获取指定zone和logger的最新数据，用于实时数据展示
func (this *ReportsController) GetLatestLoggerData() {

	zoneid := this.GetString("zoneid")
	loggerid := this.GetString("loggerid")

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	rows, _ := db.Query("select `Timestamp`,`Cpu`,`Net`,`Disk` from `tb_logger_stat` where `LoggerId`=? and `ZoneId`=? and (`Timestamp`>=now() - interval 1 minute)", loggerid, zoneid)
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

func (this *ReportsController) Brokers() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/brokers.html"

	brokerlist := []models.BrokerModel{}
	zonelist := []models.ZoneModel{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))

	rows, _ := db.Query("select * from `tb_zone`")
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

	zoneid := this.GetString("zoneid")
	if zoneid == "" {
		zoneid = zonelist[0].Id
	}

	rows, _ = db.Query("select A.`Id`,A.`Addrs`,A.`Status`,B.`Cpu`,B.`Net`,B.`Disk` from `tb_broker` as A left join `tb_latest_broker_stat` as B on A.`Id`=B.`BrokerId` and A.`ZoneId`=B.`ZoneId` and (B.`Timestamp`>=now() - interval 1 minute) where A.`ZoneId`=?", zoneid)
	for rows.Next() {
		tmp := models.BrokerModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		if err := rows.Scan(&tmp.Id, &tmp.Addrs, &tmp.Status, &Cpu, &Net, &Disk); err != nil {
			log.Println(err)
			return
		}

		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		brokerlist = append(brokerlist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["zoneid"] = zoneid
	this.Data["zonelist"] = zonelist
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["brokerlist"] = brokerlist

}

func (this *ReportsController) Loggers() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/loggers.html"

	loggerlist := []models.LoggerModel{}
	zonelist := []models.ZoneModel{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))

	rows, _ := db.Query("select * from `tb_zone`")
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

	zoneid := this.GetString("zoneid")
	if zoneid == "" {
		zoneid = zonelist[0].Id
	}

	t1 := time.Now()
	rows, _ = db.Query("select A.`Id`,A.`Addr`,A.`BlkDev`,A.`Status`,B.`Cpu`,B.`Net`,B.`Disk` from `tb_logger` as A left join `tb_latest_logger_stat` as B on A.`Id`=B.`LoggerId` and A.`ZoneId`=B.`ZoneId` and (B.`Timestamp`>=now() - interval 1 minute) where A.`ZoneId`=?", zoneid)
	t2 := time.Now()
	log.Printf("Load data from [%s] using %v\n", "broker", t2.Sub(t1))
	for rows.Next() {
		tmp := models.LoggerModel{}
		var Cpu sql.NullString
		var Net sql.NullString
		var Disk sql.NullString
		if err := rows.Scan(&tmp.Id, &tmp.Addr, &tmp.BlkDev, &tmp.Status, &Cpu, &Net, &Disk); err != nil {
			log.Println(err)
			return
		}

		tmp.Cpu, _ = strconv.ParseFloat(Cpu.String, 64)
		tmp.Net, _ = strconv.ParseFloat(Net.String, 64)
		tmp.Disk, _ = strconv.ParseFloat(Disk.String, 64)

		loggerlist = append(loggerlist, tmp)

	}
	rows.Close()
	db.Close()

	this.Data["zoneid"] = zoneid
	this.Data["zonelist"] = zonelist
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["loggerlist"] = loggerlist
}

func (this *ReportsController) Apps() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/apps.html"

	applist := []models.AppModel{}
	zonelist := []models.ZoneModel{}

	db, _ := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))

	rows, _ := db.Query("select * from `tb_zone`")
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

	zoneid := this.GetString("zoneid")
	if zoneid == "" {
		zoneid = zonelist[0].Id
	}

	rows, _ = db.Query("select `Id`,`Key`,`Status` from `tb_app` where `ZoneId`=?", zoneid)
	for rows.Next() {
		tmp := models.AppModel{}
		if err := rows.Scan(&tmp.Id, &tmp.Key, &tmp.Status); err != nil {
			log.Println(err)
			return
		}
		applist = append(applist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["zoneid"] = zoneid
	this.Data["zonelist"] = zonelist
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")
	this.Data["applist"] = applist

}

func (this *ReportsController) Topics() {
	this.Data["IsLogin"] = true
	this.Data["LoginName"] = "xianmau"
	this.Layout = "webmaster/layout.tpl"
	this.TplNames = "webmaster/topics.html"

	topiclist := []models.TopicModel{}
	zonelist := []models.ZoneModel{}
	applist := []string{}
	brokerlist := []string{}

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

	zoneid := this.GetString("zoneid")
	if zoneid == "" {
		zoneid = zonelist[0].Id
	}

	rows, err = db.Query("select `Id` from `tb_app` where `ZoneId`=?", zoneid)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var Id string
		if err := rows.Scan(&Id); err != nil {
			log.Println(err)
			return
		}
		applist = append(applist, Id)
	}
	rows.Close()

	rows, err = db.Query("select `Id` from `tb_broker` where `ZoneId`=?", zoneid)
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		var Id string
		if err := rows.Scan(&Id); err != nil {
			log.Println(err)
			return
		}
		brokerlist = append(brokerlist, Id)
	}
	rows.Close()

	appid := this.GetString("appid")
	brokerid := this.GetString("brokerid")

	if appid != "" {
		rows, err = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `ZoneId`=? and `AppId`=?", zoneid, appid)
		if err != nil {
			log.Println(err)
			return
		}
	} else if brokerid != "" {
		rows, err = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `ZoneId`=? and `BrokerId`=?", zoneid, brokerid)
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		//rows, err = db.Query("select `Name`,`AppId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status` from `tb_topic` where `ZoneId`=?", zoneid)
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
	}
	for rows.Next() {
		tmp := models.TopicModel{}
		if err := rows.Scan(&tmp.Name, &tmp.AppId, &tmp.BrokerId, &tmp.ReplicaNum, &tmp.Retention, &tmp.Segments, &tmp.Status); err != nil {
			log.Println(err)
			return
		}
		topiclist = append(topiclist, tmp)
	}
	rows.Close()
	db.Close()

	this.Data["zoneid"] = zoneid
	this.Data["zonelist"] = zonelist
	this.Data["appid"] = appid
	this.Data["applist"] = applist
	this.Data["brokerid"] = brokerid
	this.Data["brokerlist"] = brokerlist
	this.Data["topiclist"] = topiclist
	this.Data["lastsync"] = time.Now().Format("2006-01-02 15:04")

}
