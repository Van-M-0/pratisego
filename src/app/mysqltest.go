package main

import (
	"fmt"
	"net"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
var (
	user = "root"
	pass = "1"
	dbname = "gotest"
	addr = "localhost:3306"
	netAddr = ""
	dsn = ""
)

func mysql_test() {
	mysql_init()

	//testEmptyQuery()
	testcurd()
}

func mysql_init() {
	netAddr = fmt.Sprintf("%s(%s)", "tcp", addr)
	dsn = fmt.Sprintf("%s:%s@%s/%s?timeout=30s&strict=true", user, pass, netAddr, dbname)
	c, err := net.Dial("tcp", addr)
	if err == nil {
		c.Close()
		fmt.Println("mysql init success -> ", dsn)
	} else {
		panic("mysql init error")
	}
}

type dbtest struct {
	db 		*sql.DB
}

func (db *dbtest) fail(method, query string, err error) {
	fmt.Println("error on %s %s : %s", method, query, err.Error())
}

func (db *dbtest) mustQuery(query string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(query)
	if err != nil {
		db.fail("query", query, err)
	}
	return rows
}

func (db *dbtest) mustExec(query string, args ...interface{}) (ret sql.Result) {
	ret, err := db.db.Exec(query, args...)
	if err != nil {
		db.fail("exec", query, err)
	}
	return ret
}

func runtest(dsn string, test ...func(db *dbtest)) {

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("open mysql database error ", err.Error());
		return
	}

	defer func() {
		db.Close()
		if err := recover(); err != nil {
			fmt.Println("recover error ", err)
		}
	}()

	dbt := &dbtest{db: db}

	for _, f := range test {
		dbt.db.Exec("drop table if exists test")
		f(dbt)
	}
}

func testEmptyQuery() {
	runtest(dsn, func(dbt *dbtest) {
		rows := dbt.mustQuery("---")
		if rows.Next() {
			fmt.Println("rows on next must be false")
		}
	})
}

func testcurd() {
	runtest(dsn, func(dbt *dbtest) {

		dbt.mustExec("drop table if exists curdtest;")
		dbt.mustExec("create table curdtest (value bool);")

		raws := dbt.mustQuery("select * from curdtest")
		if raws.Next() {
			fmt.Println("unexpected data int mysql")
		}

		res := dbt.mustExec("insert into curdtest values(false)")
		res = dbt.mustExec("insert into curdtest values(true)")
		res = dbt.mustExec("insert into curdtest values(false)")
		count, err := res.RowsAffected()
		fmt.Println("insert curdtest values res : ", count, err)


		linsert, err := res.LastInsertId()
		fmt.Println("insert crudtest insert id res ", linsert, err)

		rows := dbt.mustQuery("select * from curdtest")
		rowcount := 0
		for {
			if !rows.Next() {
				break
			}

			rowcount++

			var out bool
			rows.Scan(&out)
			fmt.Println("row count ", rowcount, out)
		}

		//res = dbt.mustExec("update curdtest set value = ? where value = ?", true , false)

		res = dbt.mustExec("delete from curdtest where value = true")

		res = dbt.mustExec("delete from curdtest")
		fmt.Println(res.RowsAffected())
	})
}


