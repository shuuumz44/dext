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
		opErr = ListExp(db, args)

	case "summary":
		opErr = SumExp(db, args)

	case "update":
		opErr = UpdateExp(db, args)

	case "delete":
		opErr = DeleteExp(db, args)

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
	cfg.Passwd = "sincere-aquarium"
	// os.Getenv() stopped working. >:(
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "account"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return db, nil
}

func AddExp(db *sql.DB, args []string) (error) {
	var description string
	var amount 		float64
	//var date		string

	descriptionUsage 	:= "the name of the expense"
	amountUsage 		:= "the cost of the expense"
	//dateUsage			:= "the date of the transaction"

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	fs.StringVar(&description, "description", "", descriptionUsage)
	fs.StringVar(&description, "d", "", descriptionUsage)
	//fs.StringVar(&date, "date", "", dateUsage)

	fs.Float64Var(&amount, "amount", 0, amountUsage)
	fs.Float64Var(&amount, "a", 0, amountUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	_, exeErr := db.Exec("INSERT INTO expenses (description, amount) VALUES (?, ?)", description, amount)
	if exeErr != nil {
		return exeErr
	}

	fmt.Println("added: " + description + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return nil
}

func ListExp(db *sql.DB, args []string) (error) {
	// var filter string
	// filterUsage := "tag(s) to filter by"

	fs := flag.NewFlagSet("list", flag.ExitOnError)
	// fs.StringVar(&filter, "filter", NULL, filterUsage)
	// fs.StringVar(&filter, "f", NULL, filterUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	// eventually alter query to filter for tags.
	rows, exeErr := db.Query("SELECT * FROM expenses")
	if exeErr != nil {
		return exeErr
	}
	defer rows.Close()

	for rows.Next() {
		var exp Expense

		scanErr := rows.Scan(&exp.ID, &exp.Date, &exp.Desc, &exp.Amount) 
		if scanErr != nil {
			return scanErr
		}
		fmt.Printf("%d %s:\t\t$%.2f\t\t%s\n", exp.ID, exp.Desc, exp.Amount, exp.Date)
		// **instead of exp.Date.GoString()
	}

	return nil
}

func SumExp(db *sql.DB, args []string) (error) {
	var total float64 = 0
	// var filter string
	// filterUsage := "tag(s) to filter by"

	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	// fs.StringVar(&filter, "filter", NULL, filterUsage)
	// fs.StringVar(&filter, "f", NULL, filterUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	rows, exeErr := db.Query("SELECT * FROM expenses")
	if exeErr != nil {
		return exeErr
	}
	defer rows.Close()

	for rows.Next() {
		var exp Expense
		scanErr := rows.Scan(&exp.ID, &exp.Date, &exp.Desc, &exp.Amount) 
		if scanErr != nil {
			return scanErr
		}
		total += exp.Amount
	}
	fmt.Println("Total: ", total)

	return nil
}

func UpdateExp(db *sql.DB, args []string) (error) {
	var id 		int
	var amount 	float64
	//var date	string

	idUsage 			:= "the ID of the expense"
	amountUsage 		:= "the cost of the expense"
	//dateUsage			:= "the date of the transaction"

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	fs.IntVar(&id, "id", 0, idUsage)
	//fs.StringVar(&date, "date", "", dateUsage)

	fs.Float64Var(&amount, "amount", 0, amountUsage)
	fs.Float64Var(&amount, "a", 0, amountUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	op := "UPDATE expenses SET amount=? WHERE id=?"
	_, exeErr := db.Exec(op, amount, id)
	if exeErr != nil {
		return exeErr
	}

	//fmt.Println("updated: " + description + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return nil
} 

func DeleteExp(db *sql.DB, args []string) (error) {
	var id 		int
	//var amount 	float64
	//var date	string

	idUsage 			:= "the ID of the expense"
	//amountUsage 		:= "the cost of the expense"
	//dateUsage			:= "the date of the transaction"

	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	fs.IntVar(&id, "id", 0, idUsage)
	//fs.StringVar(&date, "date", "", dateUsage)

	if err := fs.Parse(args); err != nil {
		return err
	}

	op := "DELETE FROM expenses WHERE id=?"
	_, exeErr := db.Exec(op, id)
	if exeErr != nil {
		return exeErr
	}

	//fmt.Println("deleted: " + description + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return nil
}
