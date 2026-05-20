# NAME
    dext: Expense Tracker

# SYNOPSIS
    **dext** \[COMMAND\] \[OPTION...\] 

# DESCRIPTION
    **dext** manipulates (CRUD) a database of expenses.

# DEPENDENCIES
    * mysql

# SETUP
    * git clone/etc.
    * run go get github.com/go-sql-driver/mysql
    * using the mysql CLI:
        * CREATE DATABASE [database name]
        * USE [database name]
        * source ./init.sql
    * mysql has now created a table named "expenses", inside of the database you created. 
    The user is "root".
    The password is read from an environment variable, DBPASS.
    The database name is read from DBNAME.
    * To export a value into an environment variable:
        - Linux/Mac: export VARIABLE=value
        - Windows:   set VARIABLE=value

# OPTIONS
    **CRUD Operations**
        **add**
            add an expense into the database.
        
        **list** 
            List all expenses. 

        **summary**
            Print the total cost of expenses.
        
        **delete** ID
            Delete the expense specified by the id(s).

    **Database**
        **category**
            Manipulate the existing categories.

# TODO
    * make so many tests for this God Willing
should've went with TDD...
    * ParseExec(): UPDATE/DELETE
    * add back in filter latest year (when filtering by month)
    * make config flag
        ^ toggle do not add expenses that exceed budget
        ^ toggle if budget is present in listing expenses
    * DeleteExp()
        ^ select by amount/date/name
    * Querying
        ^ choose by year
        ^ choose by ranges (dates or amount)
        ^ choose column to order rows by, and increasing/decreasing
        ^ abstract case switching to a helper function
        ^ specifying by month chooses the most recent year. make a toggle for this
        ^ date parsing (different formats, etc. with slashes: YYYY/MM/DD)
    * Output
        ^ space output rows cleanly
        ^ add export to JSON

# INTUITION
    * db.Exec() and db.Query() return different types in the sql package. Therefore there should be two
    separate functions to parse them.

# TEST
    * AddExp()
        - handle any missing values
    * Summary()
    * ParseExec()
    * CollectFlag()
    * Sanitize()

# ERRORS
    * Sanitize() mismatched arguments in test

# IDEAS
    * streamline setup
    * manage user credentials/authentication 
    * output formatting options
    * make testing utilize a temporary database/directory
    * multiple tables/databases

# NOTES
    * mysql authentication can be a little finnicky at first. Make sure you can login at the command prompt, check the creds with status, then copy them to the program / as an env variable.
    * argument parsing is basically one and the same for building out the specific function. It doesn't have to be, it can and should probably be abstracted, but it's okay for now
    * go's sql package recognizes mysql DATETIMEs as time.Time types. When calling row(s).Scan, this can be stored into pointers to time.Time, interface{}, string, or []byte. Warning though, time.Time types (and sql.NullTime by extension) have to be used with the sql.Scanner interface. Despite all this, it is easiest to store DATETIMEs as strings.
    * sql accepts DATETIME values as strings or integers. Strings are accepted as 'YYYY-MM-DD hh:mm:ss (with an optional fraction part). Time is in military time. If the time is omitted, it is set to 0:00:00.
    * when os.GetEnv() fails, be sure that your variable is exported (in the same shell that you build the program, if using a multiplexer).
    * Rows can be selected by the row number with the ROW_NUMBER clause. You only have to specify 
    the order of the table in the same statement. Unfortunately this is useless as the program exits
    after every execution.
    * if you have multiple .go files and/or your main file isn't named 'main.go', go build might
    tell you it found packages and not actually create the output file. This began after I added
    'dext_test.go". If this happens, just replace . with your program name (output file).
    * errors should be found and handled in testing, not while compiling/running.
