package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

//登录
func Login(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	r.ParseForm() //解析参数，默认是不会解析的
	uname := r.Form.Get("name")
	upassword := r.Form.Get("password")
	sql := "select id from `user` where username=? and password=?"
	defer db.Close()
	userinfo := make(map[string]interface{})
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Query(uname, upassword)
	for rows.Next() {
		var uid int
		err = rows.Scan(&uid)
		checkErr(err)
		userinfo["id"] = uid
		b, err := json.Marshal(userinfo)
		checkErr(err)
		fmt.Fprintf(w, string(b))
	}
}

// 预定
func reserved(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	uname := r.Form.Get("name")
	seatid := r.Form.Get("seatid")
	adtime := r.Form.Get("adTime")
	var sql, seatsql string
	if adtime == "1" {
		sql = "update `user` set seatid=? , adTime=? ,state = 1 where username=?"
		seatsql = "update `seatinfo` set used = 1 ,Time1=1 where id=?"
	} else if adtime == "2" {
		sql = "update `user` set seatid=? , adTime=? ,state = 1 where username=?"
		seatsql = "update `seatinfo` set used = 1 ,Time2=1 where id=?"
	} else if adtime == "3" {
		sql = "update `user` set seatid=? , adTime=? ,state = 1 where username=?"
		seatsql = "update `seatinfo` set used = 1 ,Time3=1 where id=?"
	}
	defer db.Close()
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Exec(seatid, adtime, uname)
	affect, err := rows.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	stmt, err := db.Prepare(seatsql)
	row, err := stmt.Exec(seatid)
	affects, err := row.RowsAffected()
	checkErr(err)
	fmt.Println(affects)
}

//个人Info
func userinfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	uname := r.Form.Get("name")
	sql := "select * from `user` where username=?"
	defer db.Close()
	userinfo := make(map[string]interface{})
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Query(uname)
	for rows.Next() {
		var username, upassword, uid, seatid, state string
		var adtime string
		err = rows.Scan(&uid, &username, &upassword, &seatid, &adtime, &state)
		checkErr(err)
		userinfo["id"] = uid
		userinfo["name"] = username
		userinfo["password"] = upassword
		userinfo["seatid"] = seatid
		userinfo["adtime"] = adtime
		userinfo["state"] = state
		b, err := json.Marshal(userinfo)
		checkErr(err)
		fmt.Fprintf(w, string(b))
	}
}

//签到
func SignSeat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	uname := r.Form.Get("name")
	seatid := r.Form.Get("seatid")
	sql := "update `user` set state = 2 where username=?"
	defer db.Close()
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Exec(uname)
	affect, err := rows.RowsAffected()
	fmt.Println(affect)
	checkErr(err)
	seatsql := "update `seatinfo` set used = 2 where id=?"
	stmt, err := db.Prepare(seatsql)
	row, err := stmt.Exec(seatid)
	affects, err := row.RowsAffected()
	fmt.Println(affects)
	checkErr(err)
}

//座位状态
func seatinfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	time := r.Form.Get("time")
	number := r.Form.Get("number")
	var sql string
	if number == "1" {
		sql = "select * from `seatinfo` where Time1=?"
	} else if number == "2" {
		sql = "select * from `seatinfo` where Time2=?"
	} else if number == "3" {
		sql = "select * from `seatinfo` where Time3=?"
	}
	defer db.Close()
	seatinfo := make(map[string]interface{})
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Query(time)
	jsonString := "["
	for rows.Next() {
		var id, used, time1, time2, time3 string
		err = rows.Scan(&id, &used, &time1, &time2, &time3)
		checkErr(err)
		seatinfo["id"] = id
		seatinfo["used"] = used
		seatinfo["time"] = time
		b, err := json.Marshal(seatinfo)
		checkErr(err)
		jsonString = jsonString + string(b) + ","
		//fmt.Fprintf(w, string(b))
	}
	jsonString = jsonString + "{}]"
	fmt.Fprintf(w, string(jsonString))
}

//离开座位
func leave(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	uname := r.Form.Get("name")
	id := r.Form.Get("seatid")
	seatime := r.Form.Get("time")
	sql := "update `user` set seatid=0 ,adTime=0 ,state=0 where username=?"
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Exec(uname)
	affects, err := rows.RowsAffected()
	checkErr(err)
	fmt.Println(affects)
	var seatsql string
	if seatime == "1" {
		seatsql = "update `seatinfo` set used = 0 ,Time1 = 0 where id=?"
	} else if seatime == "2" {
		seatsql = "update `seatinfo` set used = 0 ,Time2 = 0 where id=?"
	} else if seatime == "3" {
		seatsql = "update `seatinfo` set used = 0 ,Time3 = 0 where id=?"
	}
	stmt, err := db.Prepare(seatsql)
	row, err := stmt.Exec(id)
	affect, err := row.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	defer db.Close()
}

//检查座位
func checkseat(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析参数，默认是不会解析的
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/appdb?charset=utf8")
	checkErr(err)
	id := r.Form.Get("id")
	time := r.Form.Get("time")
	var sql string
	if time == "1" {
		sql = "select id from `seatinfo` where Time1 = 1 and id = ?"
	} else if time == "2" {
		sql = "select id from `seatinfo` where Time2 = 1 and id = ?"
	} else if time == "3" {
		sql = "select id from `seatinfo` where Time3 = 1 and id = ?"
	}
	defer db.Close()
	seatinfo := make(map[string]interface{})
	stmtOut, err := db.Prepare(sql)
	rows, err := stmtOut.Query(id)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		checkErr(err)
		seatinfo["id"] = id
		b, err := json.Marshal(seatinfo)
		checkErr(err)
		fmt.Fprintf(w, string(b))
	}
}

func main() {
	http.HandleFunc("/Login", Login)
	http.HandleFunc("/leave", leave)
	http.HandleFunc("/SignSeat", SignSeat)
	http.HandleFunc("/userinfo", userinfo)
	http.HandleFunc("/reserved", reserved)
	http.HandleFunc("/seatinfo", seatinfo)
	http.HandleFunc("/checkseat", checkseat)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
