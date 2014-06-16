package dbmaster

import (
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func ExecuteSql(sqlstr string) (int64, error) {
	db, err := sql.Open("mysql", beego.AppConfig.String("mysql_conn_str"))
	defer db.Close()
	if err != nil {
		return 0, err
	}
	res, err := db.Exec(sqlstr)
	log.Println(sqlstr)
	if err != nil {
		return 0, err
	}
	if n, err := res.RowsAffected(); err != nil {
		return 0, err
	} else {
		return n, err
	}
	return 0, nil
}
