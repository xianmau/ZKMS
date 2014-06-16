package extractor

import (
	"ZKMS/models/datastructure"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)

const (
	EXTRACTSTATUS_RUNNING = 1
	EXTRACTSTATUS_STOPPED = 0
)

var (
	EXTRACTORSTATUS = EXTRACTSTATUS_STOPPED
)
var (
	zones map[string][]string
)

func RunExtractor() error {
	var lock sync.Mutex

	lock.Lock()
	if EXTRACTORSTATUS == EXTRACTSTATUS_RUNNING {
		lock.Unlock()
		return errors.New("already running")
	}
	lock.Unlock()

	// 设置状态为运行状态
	lock.Lock()
	EXTRACTORSTATUS = EXTRACTSTATUS_RUNNING
	lock.Unlock()

	go func() {
		// 定时抓取数据
		timer := time.Tick(60 * time.Second)
		for _ = range timer {
			// 只有当前状态为运行时才执行抓取
			lock.Lock()
			if EXTRACTORSTATUS == EXTRACTSTATUS_STOPPED {
				lock.Unlock()
				break
			}
			lock.Unlock()

			// 先从数据库读取zones
			zones = make(map[string][]string)
			if err := getZonesMap(); err != nil {
				log.Println("N0", err)
				continue
			}
			// 设置同步
			var wg sync.WaitGroup
			for zoneid, ips := range zones {
				wg.Add(1)
				go func(zoneid string, ips []string) {
					defer wg.Done()
					if err := extracting(ips); err != nil {
						log.Println("N1", err)
					}
					log.Printf("zone [%s] done\n", zoneid)
				}(zoneid, ips)
			}
			wg.Wait()
			log.Println("all zones done")
		}
	}()

	return nil
}

func StopExtractor() error {
	var lock sync.Mutex
	lock.Lock()
	EXTRACTORSTATUS = EXTRACTSTATUS_STOPPED
	lock.Unlock()

	return nil
}

// 从数据库中获取所有的zones
func getZonesMap() error {
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		return err
	}
	rows, err := db.Query("select * from `tb_zone`")
	defer rows.Close()
	if err != nil {
		return err
	}
	for rows.Next() {
		var Id string
		var Ip string
		if err := rows.Scan(&Id, &Ip); err != nil {
			return err
		}
		var ips []string
		if err := json.Unmarshal([]byte(Ip), &ips); err != nil {
			return err
		}
		zones[Id] = ips
	}
	return nil
}

func extracting(host []string) error {
	conn, _, err := zk.Connect(host, time.Second)
	defer conn.Close()
	if err != nil {
		return err
	}

	// 实例化一个ymb，把从zookeeper上的数据先放到里面
	//t1 := time.Now()
	ymb := datastructure.YMB{
		[]datastructure.BrokerInfo{},
		[]datastructure.LoggerInfo{},
		[]datastructure.AppInfo{},
		[]datastructure.TopicInfo{},
		string(""),
		[]string{},
		[]datastructure.BrokerTopic{},
	}
	extractingZoneId(conn, &ymb)
	extractingRemoteZones(conn, &ymb)
	extractingApps(conn, &ymb)
	extractingBrokers(conn, &ymb)
	extractingLoggers(conn, &ymb)
	extractingTopics(conn, &ymb)
	//t2 := time.Now()
	//log.Printf("Extract data from [%s] using %v\n", ymb.ZoneId, t2.Sub(t1))

	// 再把ymb持久化到数据库
	//t1 = time.Now()
	persistToLocalStorageUsingTx(&ymb)
	//t2 = time.Now()
	//log.Printf("Store data in database using %v\n", t2.Sub(t1))

	return nil
}

// extract zoneid znodes' data
func extractingZoneId(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/zoneid"
	if flag, _, err := conn.Exists(path); err == nil && flag {
		if data, _, err := conn.Get(path); err == nil {
			ymb.ZoneId = string(data)
			log.Println(ymb.ZoneId)
		} else {
			log.Println("N2", err)
		}
	} else if err != nil {
		log.Println("N3", err)
	}
}

// extract remote_zones znodes' data
func extractingRemoteZones(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/remote_zones"
	if flag, _, err := conn.Exists(path); err == nil && flag {
		if data, _, err := conn.Get(path); err == nil {
			var item []string
			if err := json.Unmarshal([]byte(data), &item); err != nil {
				log.Println("N4", err)
			}
			ymb.RemoteZones = item
		} else {
			log.Println("N5", err)
		}
	} else if err != nil {
		log.Println("N6", err)
	}
}

// extract app znodes' data
func extractingApps(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/appid"
	children, _, err := conn.Children(path)
	if err != nil {
		log.Println("N7", err)
	}

	for _, child := range children {
		cur_znode := path + "/" + child
		if data, _, err := conn.Get(cur_znode); err == nil {
			var item datastructure.AppInfo
			if err := json.Unmarshal([]byte(data), &item); err != nil {
				log.Println("N8", err)
			}
			ymb.Apps = append(ymb.Apps, item)
		}
	}
}

// extract logger znodes' data
func extractingLoggers(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/loggers"
	children, _, err := conn.Children(path)
	if err != nil {
		log.Println("N9", err)
	}

	for _, znode := range children {
		cur_znode := path + "/" + znode
		if data, _, err := conn.Get(cur_znode); err == nil {
			var item datastructure.LoggerInfo
			if err := json.Unmarshal([]byte(data), &item); err != nil {
				log.Println("N10", err)
			}
			ymb.Loggers = append(ymb.Loggers, item)
		}
	}
}

// extract broker znodes' data
func extractingBrokers(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/brokers"
	children, _, err := conn.Children(path)
	if err != nil {
		log.Println("N11", err)
	}

	for _, znode := range children {
		cur_znode := path + "/" + znode
		if data, _, err := conn.Get(cur_znode); err == nil {
			var item datastructure.BrokerInfo
			if err := json.Unmarshal([]byte(data), &item); err != nil {
				log.Println("N12", err)
			}
			ymb.Brokers = append(ymb.Brokers, item)
		}
	}
}

// extract topic znodes' data
func extractingTopics(conn *zk.Conn, ymb *datastructure.YMB) {
	path := "/ymb/topics"
	children, _, err := conn.Children(path)
	if err != nil {
		log.Println("N13", err)
	}

	for _, znode := range children {
		sub_path := path + "/" + znode
		children2, _, err := conn.Children(sub_path)
		if err != nil {
			log.Println("N14", err)
		}

		for _, sub_znode := range children2 {
			sub_sub_path := sub_path + "/" + sub_znode
			if data, _, err := conn.Get(sub_sub_path); err == nil {
				var item datastructure.TopicInfo
				if err := json.Unmarshal([]byte(data), &item); err != nil {
					log.Println("N15", err)
				}
				ymb.Topics = append(ymb.Topics, item)

				children3, _, err := conn.Children(sub_sub_path)
				if err != nil {
					log.Println("N16", err)
				}
				for _, sub_sub_znode := range children3 {
					var item2 = datastructure.BrokerTopic{}
					item2.AppId = znode
					item2.TopicName = sub_znode
					item2.BrokerId = sub_sub_znode
					ymb.BrokersTopics = append(ymb.BrokersTopics, item2)
				}
			}
		}
	}
}

// persist object to local storage
func persistToLocalStorageUsingTx(ymb *datastructure.YMB) error {
	var (
		db   *sql.DB
		stmt *sql.Stmt
		tx   *sql.Tx
	)
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		return err
	}

	// 执行事务
	tx, err = db.Begin()
	if err != nil {
		log.Println("N17", err)
	}

	// store app
	_, err = tx.Exec("update `tb_app` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	if err != nil {
		log.Println("N18", err)
	}
	for _, e := range ymb.Apps {
		stmt, err = tx.Prepare("insert into `tb_app`(`Id`,`ZoneId`,`Key`,`Status`) values (?,?,?,?) on duplicate key update `Key`=values(`Key`),`Status`=true")
		if err != nil {
			log.Println("N19", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Key, true)
		if err != nil {
			log.Println("N20", err)
		}
	}

	// store broker
	_, err = tx.Exec("update `tb_broker` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	if err != nil {
		log.Println("N21", err)
	}
	for _, e := range ymb.Brokers {
		stmt, err = tx.Prepare("insert into `tb_broker`(`Id`,`ZoneId`,`Addrs`,`Status`) values (?,?,?,?) on duplicate key update `Addrs`=values(`Addrs`),`Status`=true")
		if err != nil {
			log.Println("N22", err)
		}
		addrs_json, _ := json.Marshal(e.Addrs)
		_, err = stmt.Exec(e.Id, ymb.ZoneId, string(addrs_json), true)
		if err != nil {
			log.Println("N23", err)
		}
	}

	// store broker_stat
	for _, e := range ymb.Brokers {
		stmt, err = tx.Prepare("insert into `tb_broker_stat`(`BrokerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		if err != nil {
			log.Println("N24", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		if err != nil {
			log.Println("N25", err)
		}
	}

	// store latest broker_stat
	stmt, err = tx.Prepare("delete from `tb_latest_broker_stat` where `ZoneId`=?")
	if err != nil {
		log.Println("N26", err)
	}
	_, err = stmt.Exec(ymb.ZoneId)
	for _, e := range ymb.Brokers {
		stmt, err = tx.Prepare("insert into `tb_latest_broker_stat`(`BrokerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		if err != nil {
			log.Println("N27", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		if err != nil {
			log.Println("N28", err)
		}
	}

	// store logger
	_, err = tx.Exec("update `tb_logger` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	if err != nil {
		log.Println("N29", err)
	}
	for _, e := range ymb.Loggers {
		stmt, err = tx.Prepare("insert into `tb_logger`(`Id`,`ZoneId`,`Addr`,`BlkDev`,`Status`) values (?,?,?,?,?) on duplicate key update `Addr`=values(`Addr`),`BlkDev`=values(`BlkDev`),`Status`=true")
		if err != nil {
			log.Println("N30", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Addr, e.BlkDev, true)
		if err != nil {
			log.Println("N31", err)
		}
	}

	// store logger_stat
	for _, e := range ymb.Loggers {
		stmt, err = tx.Prepare("insert into `tb_logger_stat`(`LoggerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		if err != nil {
			log.Println("N32", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		if err != nil {
			log.Println("N33", err)
		}
	}

	// store latest logger_stat
	stmt, err = tx.Prepare("delete from `tb_latest_logger_stat` where `ZoneId`=?")
	if err != nil {
		log.Println("N34", err)
	}
	_, err = stmt.Exec(ymb.ZoneId)
	for _, e := range ymb.Loggers {
		stmt, err = tx.Prepare("insert into `tb_latest_logger_stat`(`LoggerId`,`ZoneId`,`CPU`,`Net`,`Disk`) values (?,?,?,?,?)")
		if err != nil {
			log.Println("N35", err)
		}
		_, err = stmt.Exec(e.Id, ymb.ZoneId, e.Stat.Cpu, e.Stat.Net, e.Stat.Disk)
		if err != nil {
			log.Println("N36", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("N37", err)
	}

	// store topic
	tx, err = db.Begin()
	_, err = tx.Exec("update `tb_topic` set `Status`=false where `zoneid`=?", ymb.ZoneId)
	if err != nil {
		log.Println("N38", err)
	}
	for _, e := range ymb.Topics {
		stmt, err = tx.Prepare("insert into `tb_topic`(`Name`,`AppId`,`ZoneId`,`BrokerId`,`ReplicaNum`,`Retention`,`Segments`,`Status`) values (?,?,?,?,?,?,?,?) on duplicate key update `BrokerId`=values(`BrokerId`),`ReplicaNum`=values(`ReplicaNum`),`Retention`=values(`Retention`),`Segments`=values(`Segments`),`Status`=true")
		if err != nil {
			log.Println("N39", err)
		}
		segments_json, _ := json.Marshal(e.Segments)
		_, err = stmt.Exec(e.Name, e.AppId, ymb.ZoneId, GetBrokerId(ymb, e.Name, e.AppId), e.ReplicaNum, e.ReplicaNum, segments_json, true)
		if err != nil {
			log.Println("N40", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Println("N41", err)
	}

	return nil
}

func GetBrokerId(ymb *datastructure.YMB, Name string, AppId string) string {
	for _, v := range ymb.BrokersTopics {
		if v.TopicName == Name && v.AppId == AppId {
			return v.BrokerId
		}
	}
	return ""
}

// traverse all znodes under the specified path
func traverse(conn *zk.Conn, path string) {
	children, _, err := conn.Children(path)
	if err != nil {
		log.Println(err)
	}

	if len(children) <= 0 {
		data, _, err := conn.Get(path)
		if err == nil {
			fmt.Println("#Leaf ZNode Found:")
			fmt.Println("#PATH: ", path)
			fmt.Println("#DATA: ", string(data))
		}
	}
	for _, znode := range children {
		if path == "/" {
			fmt.Printf("Searching ZNode: /%s\n", znode)
			traverse(conn, "/"+znode)
		} else {
			fmt.Printf("Searching ZNode: %s/%s\n", path, znode)
			traverse(conn, path+"/"+znode)
		}
	}
}
