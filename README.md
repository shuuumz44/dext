# NAME
    dext: Expense Tracker

# SYNOPSIS
    **dext** [OPTION] ... 

# DESCRIPTION
    **dext** manipulates (CRUD) a database of expenses.
    **dext** allows setting a monthly budget, and will optionally warn you if
    that budget is exceeded after an expense. This warning is automatically on,
    but can be toggled off.
    It also allows exporting a file into a CSV.

# OPTIONS
    **File Input/Output**
        **-o open**
            opens a database for operating on.

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
    * make GetConfig take optional database argument 

# NOTES
