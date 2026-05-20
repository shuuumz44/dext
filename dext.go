package main

import (
	"encoding/csv"
	"errors"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-sql-driver/mysql"
)

type Expense struct {
	ID 			int
	Amount		float64
	Name		string
	Date		string
	Created 	string
	Category	*string
}


type Column struct {
	columnName	string
	columnValue	any
}

type execAction int

// func (c category) String() string {}

const (
	NULL	execAction = iota
	INSERT
	SELECT
	UPDATE
	DELETE
)

var DEFAULTAMOUNT float64 = -1.1111
var help string =
`Usage: dext [COMMAND] [OPTION...]
COMMANDS:
	CRUD
		add
		update
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

	-b, --budget
		set the amount of the budget

	-c, --category
		manage categories
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

	switch os.Args[1] {
	case "add":
		_, opErr = AddExp(db, args)

	case "list":
		opErr = ListExp(db, args)

	case "summary":
		opErr = SumExp(db, args)

	case "update":
		_, opErr = UpdateExp(db, args)

	case "delete":
		_, opErr = DeleteExp(db, args)
	
	case "budget":
		_, opErr = Budget(db, args)

	case "category":
		_, opErr = Category(db, args)

	case "export":
		opErr = Export(db, args)

	default:
		fmt.Println(help)
		return
	}

	if opErr != nil {
		fmt.Printf("error: %s\n", os.Args[1])
		log.Fatal(opErr)
	}

}

// connect to database
func GetConfig() (*sql.DB, error) {
	cfg := mysql.NewConfig()
	cfg.User = "user"
	cfg.Passwd = os.Getenv("DBPASS") 
	cfg.Net = "tcp"
	cfg.Addr = "127.0.0.1:3306"
	cfg.DBName = os.Getenv("DBNAME") 

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

// determine inputted string is not malicious
func Sanitize(str string) error {
	invalid := errors.New("invalid string")

	if len(str) >= 100 {
		return errors.New("string too long")
	}

	for _, a := range str {
		if !unicode.IsLetter(a) && !unicode.IsNumber(a) && !unicode.IsSpace(a) {
			return invalid
		} 
	}

	return nil
}

// determine inputted string is a valid month, whether or not it maps to an integer.
func GetMonth(m string) (int, error) {
	invalid := errors.New("invalid month")

	n, notNumber := strconv.Atoi(m)
	if notNumber == nil {
		if n < 1 || n > 12 {
			return 0, invalid
		}
		return n, nil
	}

	month := make([]rune, len(m))
	for i, a := range m {
		month[i] = unicode.ToLower(a)
	}

	switch string(month) {
	case "january", "jan":
		return 1, nil
	case "february", "feb":
		return 2, nil
	case "march", "mar":
		return 3, nil
	case "april", "apr":
		return 4, nil
	case "may":
		return 5, nil
	case "june", "jun":
		return 6, nil
	case "july", "jul":
		return 7, nil
	case "august", "aug":
		return 8, nil
	case "september", "sep":
		return 9, nil
	case "october", "oct":
		return 10, nil
	case "november", "nov":
		return 11, nil
	case "december", "dec":
		return 12, nil
	default:
		return 0, invalid
	}
}

// dynamically create/execute an sql executable, using a Column's values
func ParseExec(db *sql.DB, values []Column, action execAction) (*sql.Result, error) {
	var r, v strings.Builder

	switch action {
	case NULL:
		return nil, nil

	case INSERT:
		var list []any
		for i, c := range values {
			if i == 0 {
				fmt.Fprint(&r, c.columnName)
				fmt.Fprint(&v, "?")
			} else {
				fmt.Fprintf(&r, ", %s", c.columnName)
				fmt.Fprintf(&v, ", ?")
			}
			list = append(list, c.columnValue)
		}

		rr := r.String()
		vv := v.String()
		q := fmt.Sprintf("INSERT INTO expenses (%s) VALUES (%s)", rr, vv)

		// for error checking
		// fmt.Printf("r: %s\nv: %s\n", rr, vv)
		// fmt.Printf("query: %s\n", q)

		res, exeErr := db.Exec(q, list...)
		return &res, exeErr

	case UPDATE:
	case DELETE:
	}

	return nil, errors.New("execute case unrecognized")
}

// dynamically create/execute an sql query, using a []Column's values
func ParseQuery(db *sql.DB, values []Column) (*sql.Rows, error) {
	var w strings.Builder
	var list []any

	// filter by year (make helper function)
	// year := GetYear()
	// row := db.QueryRow("SELECT purchased FROM expenses ORDER BY purchased DESC LIMIT 1")
	// row.Scan(&year)
	// y := year[:4]
	//q = fmt.Sprintf(`%s WHERE purchased LIKE '%s-%s-%%'`, q, y, monthNumber)

	for i, c := range values {
		if c.columnName == "month" {
			var monthNumber string
			m := c.columnValue.(int)
			if m < 10 {
				monthNumber = fmt.Sprintf("0%d", m)
			} else {
				monthNumber = fmt.Sprintf("%d", m)
			}
			c.columnValue = monthNumber
		}

		if i == 0 {
			fmt.Fprintf(&w, "WHERE %s=%s", c.columnName, c.columnValue)
		} else {
			fmt.Fprintf(&w, "AND %s=%s", c.columnName, c.columnValue)
		}
		list = append(list, c.columnValue)
	}

	ww := w.String()
	q := fmt.Sprintf("SELECT * FROM expenses %s", ww)

	// for error checking
	// fmt.Printf("w: %q\n", ww)
	// fmt.Printf("query: %q\n", q)

	rows, err := db.Query(q, list...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// add flag value to collect if it was used
func CollectFlag(a any, name string, val *[]Column) {
	switch t := a.(type) {
	case string:
		if t != "" {
			unsafeString := Sanitize(t)
			if unsafeString != nil {
				panic(unsafeString)
			}

			if name == "month" {
				_, invalidMonth := GetMonth(t)
				if invalidMonth != nil {
					panic(invalidMonth)
				}
			}

			*val = append((*val), Column {name, t})
			fmt.Printf("%v appended:\t%v\n", t, a)
		}

	case int:
		if t != 0 {
			*val = append((*val), Column {name, t})
			fmt.Printf("%v appended:\t%v\n", t, a)
		}

	case float64:
		if t != DEFAULTAMOUNT {
			*val = append((*val), Column {name, t})
			fmt.Printf("%v appended:\t%v\n", t, a)
		}

	default:
		err := fmt.Errorf("type unrecognized.")
		panic(err)
	}
}

// add a flag value to collect whether or not it was used
func ForceCollectFlag(a any, name string, val *[]Column) {
	switch t := a.(type) {
	case string:
		unsafeString := Sanitize(t)
		if unsafeString != nil {
			panic(unsafeString)
		}
	}
	*val = append((*val), Column {name, a})
}

// add an expense 
func AddExp(db *sql.DB, args []string) (*sql.Result, error) {
	var name, date, category	string
	var amount, limit, total 	float64
	var values []Column

	nameUsage 		:= "the name of the expense"
	dateUsage 		:= "the date of the transaction"
	amountUsage 	:= "the cost of the expense"
	categoryUsage 	:= "the category of the expense"

	fs := flag.NewFlagSet("add", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("add usage:")
		fmt.Printf("-n, --name\n\t%s\n", nameUsage)
		fmt.Printf("-d, --date\n\t%s\n", dateUsage)
		fmt.Printf("-a, --amount\n\t%s\n", amountUsage)
		fmt.Printf("-c, --category\n\t%s\n", categoryUsage)
	}

	fs.StringVar(&name, "name", "", "")
	fs.StringVar(&name, "n", "", "")

	fs.StringVar(&date, "date", "", "")
	fs.StringVar(&date, "d", "", "")

	fs.StringVar(&category, "category", "", "")
	fs.StringVar(&category, "c", "", "")

	fs.Float64Var(&amount, "amount", 0, "")
	fs.Float64Var(&amount, "a", 0, "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: %w", err)
		return nil, err
	}

	if name == "" && amount == 0 {
		fmt.Println("No value specified.")
		return nil, nil
	}

	CollectFlag(name, "name", &values)
	CollectFlag(amount, "amount", &values)
	CollectFlag(date, "purchased", &values)
	CollectFlag(category, "category", &values)

	res, exeErr := ParseExec(db, values, INSERT)
	if exeErr != nil {
		return nil, exeErr
	}

	// compare to budget, warn accordingly
	b := db.QueryRow("SELECT threshold FROM budget LIMIT 1")
	b.Scan(&limit)
	t := db.QueryRow("SELECT SUM(amount) FROM expenses")
	t.Scan(&total)

	if (limit > 0 && total > limit) {
		fmt.Println("WARNING: budget exceeded")
	}
	fmt.Printf("added: %s - $%.2f\n", name, amount)
	return res, nil
}

// list all expenses
func ListExp(db *sql.DB, args []string) error {
	var month, category	string
	var total, limit 	float64
	var values []Column

	monthUsage 		:= "month to filter by"
	categoryUsage 	:= "category to filter by"

	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("list usage:")
		fmt.Printf("-m, --month\n\t%s\n", monthUsage)
		fmt.Printf("-c, --category\n\t%s\n", categoryUsage)
	}

	fs.StringVar(&month, "month", "", "")
	fs.StringVar(&month, "m", "", "")

	fs.StringVar(&category, "category", "", "")
	fs.StringVar(&category, "c", "", " ")

	if err := fs.Parse(args); err != nil {
		return err
	}

	b := db.QueryRow("SELECT threshold FROM budget LIMIT 1")
	b.Scan(&limit)

	CollectFlag(month, "month",  &values) 
	CollectFlag(category, "category", &values)

	rows, exeErr := ParseQuery(db, values)
	if exeErr != nil {
		return exeErr
	}
	defer rows.Close()

	// print column names
	col, colErr := rows.Columns()
	if colErr != nil {
		return colErr
	}
	fmt.Printf("%s\t%s\t\t%s\t\t%s\t\t%s\n", col[0], col[1], col[2], col[3], col[4])

	for rows.Next() {
		var exp Expense

		scanErr := rows.Scan(&exp.ID, &exp.Category, &exp.Name, &exp.Amount, &exp.Date, &exp.Created) 
		if scanErr != nil {
			return scanErr
		}
		fmt.Printf("%d\t%v\t\t%s\t\t$%.2f\t\t%s\n", exp.ID, exp.Category, exp.Name, exp.Amount, exp.Date)
		total += exp.Amount
	}

	if (limit > 0 && total > limit) {
		fmt.Println("WARNING: budget exceeded")
	}

	return nil
}

// sum all expenses in latest month
func SumExp(db *sql.DB, args []string) error {
	var month 			string
	var total, limit 	float64

	monthUsage := "month to filter by"

	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("summary usage:")
		fmt.Printf("-m, --month\n\t%s\n", monthUsage)
	}

	fs.StringVar(&month, "month", "", "")
	fs.StringVar(&month, "m", "", "")

	if err := fs.Parse(args); err != nil {
		return err
	}

	q := "SELECT SUM(amount) FROM expenses"

	b := db.QueryRow("SELECT threshold FROM budget LIMIT 1")
	b.Scan(&limit)

	if month != "" {
		unsafeMonth := Sanitize(month)
		if unsafeMonth != nil {
			return unsafeMonth
		}

		m, monthErr := GetMonth(month)
		if monthErr != nil {
			return monthErr
		}

		var year string
		var monthNumber string
		row := db.QueryRow("SELECT purchased FROM expenses ORDER BY purchased DESC LIMIT 1")
		row.Scan(&year)
		y := year[:4]

		if m < 10 {
			monthNumber = fmt.Sprintf("0%d", m)
		} else {
			monthNumber = fmt.Sprintf("%d", m)
		}

		q = fmt.Sprintf(`%s WHERE purchased LIKE '%s-%s-%%'`, q, y, monthNumber)
	}

	row := db.QueryRow(q)
	scanErr := row.Scan(&total) 
	if scanErr != nil {
		return scanErr
	}

	if (limit > 0 && total > limit) {
		fmt.Println("WARNING: budget exceeded")
	}

	fmt.Printf("Total: %.2f\n", total)
	fmt.Printf("Budget: %.2f\n", limit)
	return nil
}

// change parameters of an expense(s)
func UpdateExp(db *sql.DB, args []string) (*sql.Result, error) {
	var id 						int
	var amount, total, limit	float64
	var name, date				string

	idUsage 		:= "the ID of the expense"
	nameUsage 		:= "the name of the expense"
	dateUsage		:= "the date of the transaction"
	amountUsage		:= "the cost of the expense"

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
		fmt.Errorf("parse: %w", err)
		return nil, err
	}

	updateName := (name != "")
	updateAmount := (amount != 0) 
	if id <= 0 || (!updateName && !updateAmount) {
		fmt.Println("No changes specified.")
		return nil, nil
	}

	// make a function to validate flags
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

	t := db.QueryRow("SELECT SUM(amount) FROM expenses")
	t.Scan(&total)

	b := db.QueryRow("SELECT threshold FROM budget LIMIT 1")
	b.Scan(&limit)

	if (limit > 0 && total > limit) {
		fmt.Println("WARNING: budget exceeded")
	}

	fmt.Println("updated: " + name + " - $" + strconv.FormatFloat(amount, 'f', 2, 64))
	return &res, nil
} 

// delete an expense
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
		fmt.Errorf("parse: %w", err)
		return nil, err
	}

	if id == "" {
		fmt.Println("no id specified.")
		return nil, nil
	}

	execution := fmt.Sprintf("DELETE FROM expenses WHERE id IN (%s)", id)
	res, exeErr := db.Exec(execution)
	if exeErr != nil {
		return nil, exeErr
	}

	return &res, nil
}

// manage the budget
func Budget(db *sql.DB, args []string) (*sql.Result, error) {
	var amount 	float64 

	amountUsage		:= "the amount to set the budget"

	fs := flag.NewFlagSet("budget", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("budget usage:")
		fmt.Printf("-a, --amount\n\t%s\n", amountUsage)
	}

	fs.Float64Var(&amount, "amount", DEFAULTAMOUNT, "")
	fs.Float64Var(&amount, "a", DEFAULTAMOUNT, "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: %w", err)
		return nil, err
	}

	if (amount < 0) {
		if (amount == DEFAULTAMOUNT) {
			var b float64
			row := db.QueryRow("SELECT threshold FROM budget WHERE id=0")
			row.Scan(&b)
			fmt.Printf("budget: %.2f\n", b)
		}
		return nil, nil
	}

	res, exeErr := db.Exec("UPDATE budget SET threshold=? WHERE id=0", amount)
	if exeErr != nil {
		return nil, exeErr
	}

	fmt.Printf("set budget to: %.2f\n", amount)
	return &res, nil
}

// manage categories
func Category(db *sql.DB, args []string) (*sql.Result, error) {
	var name 	string 

	nameUsage := "the name of the new category"

	fs := flag.NewFlagSet("category", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Println("category usage:")
		fmt.Printf("-n, --name\n\t%s\n", nameUsage)
	}

	fs.StringVar(&name, "name", "", "")
	fs.StringVar(&name, "n", "", "")

	if err := fs.Parse(args); err != nil {
		fmt.Errorf("parse: %w", err)
		return nil, err
	}

	unsafeName := Sanitize(name)
	if unsafeName != nil {
		return nil, errors.New("unsafe category name")
	}

	return nil, nil
}

// export the table to a CSV
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

		scanErr := rows.Scan(&exp.ID, &exp.Date, &exp.Name, &exp.Amount, &exp.Created) 
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
