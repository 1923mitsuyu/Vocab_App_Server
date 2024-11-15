package db

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql" // MySQLドライバーをインポート
)

var db *sql.DB

func Init() {
    var err error
    // MySQL接続設定
    db, err = sql.Open("mysql", "ikufumi:1923@tcp(localhost:3306)/Vocab_App")
    if err != nil {
        fmt.Println("Error opening database:", err)
		return
    }

    // Check if the database is reachable
    err = db.Ping()
    if err != nil {
        fmt.Println("Error connecting to database:", err)
        return
    }

    fmt.Println("Database connected successfully!")
}

func GetDB() *sql.DB {
    return db
}
