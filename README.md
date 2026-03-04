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
        - On linux/mac: export VARIABLE=value
        - On Windows:   set VARIABLE=value
        (do not put any spaces between the equals sign.)

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
    * implement add
    * implement list
    * implement summary
    * implement update
    * implement delete
    * add writing to CSV
    * add filtering 

# IDEAS
    * find a better way to abstract away argument parsing
    * multiple tables/databases
    * manage user credentials/authentication 
    * softlock floats to round to 2 decimals, add an option to change precision

# NOTES
    * mysql authentication can be a little finnicky at first. Make sure you can login in the command prompt,
    check the creds with status, then copy them to the program / as an env variable.
