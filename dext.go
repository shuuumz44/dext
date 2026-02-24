package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
)

type Expense struct {
	ID 		int64
	Date	string
	Desc	string
	Amount	int
}

var help string =
`Usage: dext [OPTION]...
CRUD operations:
	-a add
	-l list
	-s summarize
	-d delete
`

func main() {
	args := os.Args
	amnt := len(args)
	if (amnt > 6) {
		return
	}
	
	db, err := GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to database.")

	// ping in main to avoid annoying unused var warning
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func GetConfig() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = "user"
	cfg.Passwd = os.Getenv("DBPASS")
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "expenses"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}
