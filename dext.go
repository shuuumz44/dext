package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type Expense struct {
	ID 		int64
	Date	string
	Desc	string
	Amount	float64
}

var help string =
`Usage: dext [COMMAND] [OPTION...]
COMMANDS:
	CRUD
		add
		list
		summary
		delete

OPTIONS:
	-d, --description
		the name of an expense
	
	-a, --amount
		the cost of an expense
`


func main() {
	args := os.Args[2:]
	amnt := len(args)
	if (amnt > 4) {
		return
	}
	

	db, confErr := GetConfig()
	if confErr != nil {
		log.Fatal("database connection error: ", confErr)
	}

	// ping in main to avoid annoying unused var warning
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal("database ping error: ", pingErr)
	}

	var opErr error
	switch os.Args[1] {
	case "add":
		opErr = AddExp(db, args)

	case "list":
		opErr = ListExp(db)

	case "summary":
		opErr = SumExp(db)

	case "update":
		opErr = UpdateExp(db)

	case "delete":
		opErr = DeleteExp(db)

	default:
		fmt.Println(help)
		return
	}

	if opErr != nil {
		log.Fatal(opErr)
		return
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

func AddExp(db *sql.DB, args []string) (error) {
	// variables
	var description string
	var amount 		int

	descriptionUsage 	:= "the name of the expense"
	amountUsage 		:= "the cost of the expense"

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	fs.StringVar(&description, "description", "", descriptionUsage)
	fs.StringVar(&description, "d", "", descriptionUsage)

	fs.IntVar(&amount, "amount", 0, amountUsage)
	fs.IntVar(&amount, "a", 0, amountUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// add query

	fmt.Println("added: " + description + " - " + strconv.Itoa(amount))
	return nil
}

func ListExp(db *sql.DB) (error) {
	//
	return nil
}

func SumExp(db *sql.DB) (error) {
	//
	return nil
}

func UpdateExp(db *sql.DB) (error) {
	//
	return nil
}

func DeleteExp(db *sql.DB) (error) {
	//
	return nil
}
