# NAME
    dext: Expense Tracker

# SYNOPSIS
    **dext** \[COMMAND\] \[OPTION...\] 

# DESCRIPTION
    **dext** manipulates (CRUD) a database of expenses.
    **dext** allows setting a monthly budget, and will optionally warn you if
    that budget is exceeded after an expense. This warning is automatically on,
    but can be toggled off.
    It also allows exporting a file into a CSV.

# DEPENDENCIES
    * mysql

# SETUP
    * git clone/etc.
    * run go get github.com/go-sql-driver/mysql
    * using the mysql CLI:
        * CREATE DATABASE [database_name]
        * USE [database_name]
        * source ./init.sql
    * mysql has now created a table named "expenses", inside of the database you created. 
    The user is "root".
    The password is read from an environment variable, DBPASS.
    * To export a value into an environment variable:
        - Linux/Mac: export VARIABLE=value
        - Windows:   set VARIABLE=value

# OPTIONS
    **File Input/Output**
        **-o open**
            opens a table for operating on.

    **CRUD Operations**
        **-a add**
            add an expense into the database.
        
        **-l list** 
            List all expenses. Can be combined with the **-f** (**filter**) category.

        **-s summary**
            Print the total cost of expenses.
        
        **-d delete** ID
            Delete the expense specified by the id(s).

    **Output Control**
        **-f filter** CATEGORY
            Filter the displayed expenses by their category.

# TODO
    * make dates added manually (if specified)
    * properly format list output
    * auto update ids
    * add writing to CSV
    * add filtering 

# ERRORS
    * GetConfig() does not accept os.Getenv() for the password.
    * error message does not combine the usage message of long and shorthand flag names

# IDEAS
    * keep table info in a struct to efficiently read stuff like amount of rows
    * find a better way to abstract away argument parsing
    * multiple tables/databases
    * manage user credentials/authentication 
    * option to format output
    * elaborate date parsing

# NOTES
    * mysql authentication can be a little finnicky at first. Make sure you can login at the command prompt, check the creds with status, then copy them to the program / as an env variable.
    * argument parsing is basically one and the same for building out the specific function. It doesn't have to be, it can and should probably be abstracted, but it's okay for now
    * go's sql package recognizes mysql DATETIMEs as time.Time types. When calling row(s).Scan, this can be stored into pointers to time.Time, interface{}, string, or []byte. Warning though, time.Time types (and sql.NullTime by extension) have to be used with the sql.Scanner interface. Despite all this, it is easiest to store DATETIMEs as strings.
    * sql accepts DATETIME values as strings or integers. Strings are accepted as 'YYYY-MM-DD hh:mm:ss (with an optional fraction part). Time is in military time. If the time is omitted, it is set to 0:00:00.
