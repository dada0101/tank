package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Register(name,pwd string) (int,bool){
	db, err := sql.Open("mysql", "root:root@/Tank3D")
	if err != nil {
		log.Println(err)
		return -1,false
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return -1,false
	}
	rows, err := db.Query("select * from user where name = ?",name)
	if err != nil {
		log.Println(err)
		return -1,false
	}
	if rows.Next() {
		log.Println("用户已存在，请尝试新的用户名！")
		return -1,false
	}
	_, err = db.Exec("insert into user(name, pwd) value(?, ?)", name, pwd)
	if err != nil {
		log.Println(err)
		return -1,false
	}
	rows, _ = db.Query("select user_id from user where name = ?",name)
	var id int
	for rows.Next() {
		rows.Scan(&id)
	}
	return id,true
}

func CheckPwd(name, pwd string) (uid int, check bool, login_cnt int) {
	uid = -1
	check = false
	db, err := sql.Open("mysql", "root:root@/Tank3D")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	rows, err := db.Query("select user_id, pwd, login_cnt from user where name = ?",name)
	if err != nil {
		log.Println(err)
		return
	}

	for rows.Next() {
		var pw string
		rows.Scan(&uid, &pw, &login_cnt)
		log.Println("db check pwd: [uid:", uid, ", pwd: ", pw, "]")
		if pw == pwd {
			check = true
			return
		}
		log.Println("用户不存在！or 密码不正确！")
		return
	}
	log.Println("用户不存在！or 密码不正确！")
	return
}

func CreateUserData(uid int) bool {
	db, err := sql.Open("mysql", "root:root@/Tank3D")
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = db.Exec("insert into user_data(user_id, score, win, fail) value(?, 0, 0, 0)", uid)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func GetUserData(uid int) (score, win, fail int, err error) {
	var db *sql.DB
	score, win, fail = 0, 0, 0
	db, err = sql.Open("mysql", "root:root@/Tank3D")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return
	}
	var rows * sql.Rows
	rows, err = db.Query("select score, win, fail from user_data where user_id = ?",uid)
	if err != nil {
		log.Println(err)
		return
	}
	if rows.Next() {
		rows.Scan(&score, &win, &fail)
		err = nil
		return
	}
	return
}

func SetUserData(uid, score, win, fail, loginCnt int) bool {
	db, err := sql.Open("mysql", "root:root@/Tank3D")
	if err != nil {
		log.Println(err)
		return false
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = db.Exec("update user set login_cnt = ? where user_id = ?",  loginCnt, uid)
	if err != nil {
		log.Println(err)
		return false
	}
	_, err = db.Exec("update user_data set score = ?, win = ?, fail = ? where user_id = ?", score, win, fail, uid)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

