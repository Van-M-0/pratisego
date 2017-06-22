package sqlpx

import (
	"database/sql"
	"fmt"
	"net"
	"time"
	"github.com/go-sql-driver/mysql"
)

type SqlProxyCfg struct {
	DbName 				string
	DbUser 				string
	DbPass 				string
	DbAddr 				string
	DbAlive 			time.Duration
	DbOption 			[]string
}

type SqlProxy struct {
	db 			*sql.DB
	cfg 		*SqlProxyCfg
	dsn 		string
	closeCh 	chan <- int
}

func NewSqlProxy() *SqlProxy {
	return &SqlProxy{
		closeCh: make(chan <- int),
	}
}

func (sp *SqlProxy) Init(cfg *SqlProxyCfg) error {
	sp.cfg = cfg
	netAddr := fmt.Sprintf("%s(%s)", "tcp", cfg.DbAddr)
	sp.dsn = fmt.Sprintf("%s:%s@%s/%s?timeout=%ds&strict=true&multiStatements=true", cfg.DbUser, cfg.DbPass, netAddr, cfg.DbName, cfg.DbAlive.Seconds())
	c, err := net.Dial("tcp", netAddr)
	if err == nil {
		c.Close()
		fmt.Println("mysql init success -> ", sp.dsn)
	} else {
		panic("mysql init error")
	}
	return  nil
}

func (sp *SqlProxy) Start() error {
	if _, err := mysql.ParseDSN(sp.dsn); err != nil {
		return err
	}
	db, err := sql.Open("mysql", sp.dsn)
	if err != nil {
		return err
	}
	sp.db = db
	return nil
}

func (sp *SqlProxy) Stop() error {
	return nil
}

func (sp *SqlProxy) Exec(query string, args ...interface{}) (sql.Result, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[EXEC RECOVER] error ", err)
		}
	}()

	res, err := sp.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[EXEC] %s -> %s : ", query, err.Error())
	}
	return res, nil
}

func (sp *SqlProxy) Query(query string, args ...interface{}) (*sql.Rows, error) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("[Query RECOVER] error ", err)
		}
	}()

	row, err := sp.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("[QUERY] %s -> %s : ", query, err.Error())
	}
	return row, nil
}



