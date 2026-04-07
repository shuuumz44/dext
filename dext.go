package main

import (
	"encoding/csv"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

type Expense struct {
	ID 		int
	Amount	float64
	Name	string
	Date	string
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
	-n, --name
		the name of an expense

	-d, --date
		the date an expense was made
	
	-a, --amount
		the cost of an expense

	--id
		the identifying number of an expense
	
	-f, --filter
		the tag(s) to filter by
`


func main() {
	if len(os.Args) <= 1 {
		return
	}

	args := os.Args[2:]
	
	db, confErr := GetConfig()
	if confErr != nil {
		log.Fatal("database connection error: ", confErr)
	}

	var opErr error
	// var res *sql.Result

	switch os.Args[1] {
	case "add":
		_, opErr = AddExp(db, args)

		/*
		lastID, IDerr := (*res).LastInsertId()
		affected, affErr := (*res).RowsAffected()
		if (IDerr != nil) {
			log.Fatal(IDerr)
		}
		if (affErr != nil) {
			log.Fatal(affErr)
		}

		fmt.Printf("last added ID: %d\n", lastID)
		fmt.Printf("altered rows: %d\n", affected)
		*/

	case "list":
		opErr = ListExp(db, args)

	case "summary":
		opErr = SumExp(db, args)

	case "update":
		_, opErr = UpdateExp(db, args)

	case "delete":
		_, opErr = DeleteExp(db, args)
		
	case "export":
		opErr = Export(db, args)

	default:
		fmt.Println(help)
		return
	}

	if opErr != nil {
		log.Fatal(opErr)
	}

}

func GetConfig() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = "user"
	cfg.Passwd = os.Getenv("DBPASS") 
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = "account"

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return db, nil
}

func AddExp(db *sql.DB, args []string) (*sql.Result, error) {
	var name	string
	var amount 	float64
	var date	string

	nameUsage 	:= "the name of the expense"
	amountUsage := "the cost of the expense"
	dateUsage 	:= "the date of the transaction"

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("add usage:")
		fmt.Printf("-n, --name\n\t%s\n", nameUsage)
		fmt.Printf("-a, --amount\n\t%s\n", amountUsage)
		fmt.Printf("-d, --date\n\t%s\n", dateUsage)
	}

	fs.StringVar(&name, "name", "", "")
	fs.StringVar(&name, "n", "", "")

	fs.StringVar(&date, "date", "", "")
	fs.StringVar(&date, "d", "", "")

	fs.Float64Var(&amount, "amount", 0, "")
	fs.Float64Var(&amount, "a", 0, "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: ", err)
		return nil, err
	}

	if name == "" && amount == 0 {
		fmt.Println("No value specified.")
		return nil, nil
	}

	var exeErr error
	var res sql.Result
	if date == "" {
		res, exeErr = db.Exec("INSERT INTO expenses (name, amount) VALUES (?, ?)", name, amount)
	} else { 
		res, exeErr = db.Exec("INSERT INTO expenses (name, amount, purchased) VALUES (?, ?, ?)", name, amount, date)
	}

	if exeErr != nil {
		return nil, exeErr
	}

	fmt.Println("added: " + name + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return &res, nil
}

func ListExp(db *sql.DB, args []string) (error) {
	// var filter string
	// filterUsage := "tag(s) to filter by"

	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("list usage:")
		//fmt.Printf("-f, --filter\n\t%s\n", filterUsage)
	}
	// fs.StringVar(&filter, "filter", NULL, "")
	// fs.StringVar(&filter, "f", NULL, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// eventually alter query to filter for tags.
	rows, exeErr := db.Query("SELECT * FROM expenses")
	if exeErr != nil {
		return exeErr
	}
	defer rows.Close()

	// print column names
	col, colErr := rows.Columns()
	if colErr != nil {
		return colErr
	}
	fmt.Printf("%s %s\t\t%s\t\t%s\n", col[0], col[1], col[2], col[3])

	for rows.Next() {
		var exp Expense

		scanErr := rows.Scan(&exp.ID, &exp.Name, &exp.Amount, &exp.Date) 
		if scanErr != nil {
			return scanErr
		}
		fmt.Printf("%d  %s\t\t$%.2f\t\t%s\n", exp.ID, exp.Name, exp.Amount, exp.Date)
	}

	return nil
}

func SumExp(db *sql.DB, args []string) (error) {
	var total float64 = 0
	// var filter string

	// filterUsage := "tag(s) to filter by"

	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("summary usage:")
		//fmt.Printf("-f, --filter\n\t%s\n", filterUsage)
	}

	// fs.StringVar(&filter, "filter", NULL, "")
	// fs.StringVar(&filter, "f", NULL, "")

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

		scanErr := rows.Scan(&exp.ID, &exp.Name, &exp.Amount, &exp.Date) 
		if scanErr != nil {
			return scanErr
		}
		total += exp.Amount
	}
	fmt.Println("Total: ", total)

	return nil
}

func UpdateExp(db *sql.DB, args []string) (*sql.Result, error) {
	var id 		int
	var amount 	float64
	var name	string
	var date	string

	idUsage 		:= "the ID of the expense"
	amountUsage		:= "the cost of the expense"
	nameUsage 		:= "the name of the expense"
	dateUsage		:= "the date of the transaction"

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("update usage:")
		fmt.Printf("--id\n\t%s\n", idUsage)
		fmt.Printf("-n, --name\n\t%s\n", nameUsage)
		fmt.Printf("-d, --date\n\t%s\n", dateUsage)
		fmt.Printf("-a, --amount\n\t%s\n", amountUsage)
	}

	fs.IntVar(&id, "id", 0, idUsage)

	fs.StringVar(&name, "name", "", "")
	fs.StringVar(&name, "n", "", "")

	fs.StringVar(&date, "d", "", "")
	fs.StringVar(&date, "date", "", "")

	fs.Float64Var(&amount, "amount", 0, "")
	fs.Float64Var(&amount, "a", 0, "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: ", err)
		return nil, err
	}

	updateName := (name != "")
	updateAmount := (amount != 0) 
	if id <= 0 || (!updateName && !updateAmount) {
		fmt.Println("No changes specified.")
		return nil, nil
	}

	var exeErr error
	var res sql.Result
	switch {
	case updateName && updateAmount:
		res, exeErr = db.Exec("UPDATE expenses SET name=?, amount=? WHERE id=?", name, amount, id)
	case updateName:
		res, exeErr = db.Exec("UPDATE expenses SET name=? WHERE id=?", name, id)
	case updateAmount:
		res, exeErr = db.Exec("UPDATE expenses SET amount=? WHERE id=?", amount, id)
	}

	if exeErr != nil {
		return nil, exeErr
	}

	fmt.Println("updated: " + name + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return &res, nil
} 

func DeleteExp(db *sql.DB, args []string) (*sql.Result, error) {
	var id 		string
	// var amount 	float64
	// var date	string

	idUsage 			:= "the ID of the expense(s)"
	// amountUsage 		:= "the cost of the expense"
	// dateUsage		:= "the date of the transaction"

	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("delete usage:")
		fmt.Printf("--id\n\t%s\n", idUsage)
		// fmt.Printf("-a, --amount\n\t%s\n", amountUsage)
		// fmt.Printf("-d, --date\n\t%s\n", dateUsage)
	}

	fs.StringVar(&id, "id", "", idUsage)

	// fs.float64Var(&amount, "amount", "", "")
	// fs.float64Var(&amount, "a", "", "")

	// fs.StringVar(&date, "date", "", "")
	// fs.StringVar(&date, "d", "", "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: ", err)
		return nil, err
	}

	if id == "" {
		fmt.Println("no id specified.")
		return nil, nil
	}

	execution := fmt.Sprintf("DELETE FROM expenses WHERE id IN (%s)", id)
	// Sanitize(execution)
	res, exeErr := db.Exec(execution)
	if exeErr != nil {
		return nil, exeErr
	}

	// fmt.Println("deleted: " + name + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return &res, nil
}

func Export(db *sql.DB, args []string) (error) {
	var filename string
	// var filter string

	filenameUsage	:= "name of .csv file"
	// filterUsage	:= "tag(s) to filter by"

	fs := flag.NewFlagSet("export", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("export usage:")
		fmt.Printf("-n, --filename\n\t%s\n", filenameUsage)
		// fmt.Printf("-f, --filter\n\t%s\n", filterUsage)
	}

	fs.StringVar(&filename, "filename", "", "")
	fs.StringVar(&filename, "n", "", "")

	// fs.StringVar(&filter, "filter", NULL, "")
	// fs.StringVar(&filter, "f", NULL, "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if filename == "" {
		fmt.Println("filename cannot be null.")
		return nil
	}

	// create connection to new .csv file
	fn := fmt.Sprintf("%s.csv", filename)
	f, createErr := os.Create(fn)
	if createErr != nil {
		return createErr
	}
	w := csv.NewWriter(f)

	rows, exeErr := db.Query("SELECT * FROM expenses")
	if exeErr != nil {
		return exeErr
	}
	defer rows.Close()

	out := make([]string, 4)
	cols, colErr := rows.Columns()
	if colErr != nil {
		return colErr
	}

	out[0] = cols[0]
	out[1] = cols[1]
	out[2] = cols[2]
	out[3] = cols[3]
	w.Write(out)

	for rows.Next() {
		var exp Expense

		scanErr := rows.Scan(&exp.ID, &exp.Date, &exp.Name, &exp.Amount) 
		if scanErr != nil {
			return scanErr
		}

		// output into CSV
		out[0] = strconv.Itoa(exp.ID)
		out[1] = exp.Name
		out[2] = fmt.Sprintf("%f", exp.Amount)
		out[3] = exp.Date
		w.Write(out)
	}

	w.Flush()

	if wErr := w.Error(); wErr != nil {
		return wErr
	}

	return nil
}
