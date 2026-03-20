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

	--id
		the identifying number of an expense
	
	-f, --filter
		the tag(s) to filter by
`


func main() {
	args := os.Args[2:]
	//amnt := len(args)
	
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
		fmt.Errorf("parse: ", err)
		return err
	}

	if description=="" && amount==0 {
		fmt.Println("No value specified.")
		return nil
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
		fmt.Printf("%d\t%s:\t\t$%.2f\t\t%s\n", exp.ID, exp.Desc, exp.Amount, exp.Date)
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
		fmt.Errorf("parse: ", err)
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
	var id 			int
	var amount 		float64
	var description	string
	//var date		string

	idUsage 			:= "the ID of the expense"
	amountUsage 		:= "the cost of the expense"
	descriptionUsage 	:= "the name of the expense"
	//dateUsage			:= "the date of the transaction"

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	fs.IntVar(&id, "id", 0, idUsage)

	fs.StringVar(&description, "description", "", descriptionUsage)
	fs.StringVar(&description, "d", "", descriptionUsage)
	//fs.StringVar(&date, "date", "", dateUsage)

	fs.Float64Var(&amount, "amount", 0, amountUsage)
	fs.Float64Var(&amount, "a", 0, amountUsage)

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: ", err)
		return err
	}

	updateDescription := (description != "")
	updateAmount := (amount != 0) 
	if id <= 0 || (!updateDescription && !updateAmount) {
		fmt.Println("No changes specified.")
		return nil
	}

	var exeErr error
	switch {
	case updateDescription && updateAmount:
		_, exeErr = db.Exec("UPDATE expenses SET description=?, amount=? WHERE id=?", description, amount, id)
	case updateDescription:
		_, exeErr = db.Exec("UPDATE expenses SET description=? WHERE id=?", description, id)
	case updateAmount:
		_, exeErr = db.Exec("UPDATE expenses SET amount=? WHERE id=?", amount, id)
	}

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
		fmt.Errorf("parse: ", err)
		return err
	}

	if id <= 0 {
		fmt.Println("id out of bounds.")
		return nil
	}

	op := "DELETE FROM expenses WHERE id=?"
	_, exeErr := db.Exec(op, id)
	if exeErr != nil {
		return exeErr
	}

	//fmt.Println("deleted: " + description + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return nil
}
