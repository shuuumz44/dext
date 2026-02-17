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

func main() {
	help :=
`help message`
	
	args := os.Args
	amnt := len(args)
	if (amnt > 6) {
		return
	}
	var user, password, table string

	cfg := GetConfig(user, password, table)
}

func GetConfig(user, password, table string) ([]Expense, error) {
	//
}
