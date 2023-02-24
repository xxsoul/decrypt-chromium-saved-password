package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type dbData struct {
	url   string
	uname string
	psw   []byte
}

// fetchChromiumPswDataFromDb Chromium浏览器获取保存的密码操作sqlite
func fetchChromiumPswDataFromDb(dbPath string, count int) ([]dbData, int) {
	logger.Println("开始读取用户数据库...")
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		logger.Println("读取用户数据库（2）错误:" + dbErr.Error())
		return nil, 0
	}
	defer db.Close()

	rows, dbErr := db.Query(fmt.Sprintf("SELECT origin_url, username_value, password_value FROM logins WHERE blacklisted_by_user <1 ORDER BY times_used DESC,date_created DESC LIMIT %d", count))
	if dbErr != nil {
		logger.Println("查询用户数据错误:" + dbErr.Error())
		return nil, 0
	}
	defer rows.Close()

	dbRes := make([]dbData, count)
	i := 0
	for rows.Next() && i < count {
		// item := dbData{}
		rows.Scan(&dbRes[i].url, &dbRes[i].uname, &dbRes[i].psw)
		i++
	}
	total := 0
	rowCount := db.QueryRow("SELECT COUNT(*) FROM logins WHERE blacklisted_by_user <1")
	rowCount.Scan(&total)
	return dbRes, total
}

// fetchChromiumHistoryDataFromDb Chromium浏览器获取浏览记录操作sqlite
func fetchChromiumHistoryDataFromDb(dbPath string, count int) ([]dbData, int) {
	logger.Println("开始读取用户数据库...")
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		logger.Println("读取用户数据库（2）错误:" + dbErr.Error())
		return nil, 0
	}
	defer db.Close()

	rows, dbErr := db.Query(fmt.Sprintf("SELECT url, title, visit_count, last_visit_time FROM urls WHERE 1=1 ORDER BY last_visit_time DESC LIMIT %d", count))
	if dbErr != nil {
		logger.Println("查询用户数据错误:" + dbErr.Error())
		return nil, 0
	}
	defer rows.Close()

	dbRes := make([]dbData, count)
	i := 0
	for rows.Next() && i < count {
		rows.Scan(&dbRes[i].url, &dbRes[i].uname, new(sql.RawBytes), new(sql.RawBytes))
		i++
	}
	total := 0
	rowCount := db.QueryRow("SELECT COUNT(*) FROM urls")
	rowCount.Scan(&total)
	return dbRes, total
}
