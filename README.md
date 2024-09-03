# Go SqliteUtils

A small lib to help with keeping track of an SQLite-Databases version and upgrades.
It is inspired by how IndexedDB in the Browser handles it's structural upgrades.

# install
```bash
go get github.com/rocco-gossmann/go_sqliteutils
```

# usage

## The `updateHandler`
Let's say for example, we have a database that was rolled out with a required version 1

Now we wan't to change the structure, but we can't just ignore all the People,
who may have used version 1 until now. Just deleting all their data and starting new
would be an option, but extremely user unfriendly.

The best idea here is a Switch-Slide.
Each Possible DB Version then becomes a `case` in that `switch`.
By ending each `case` with a `fallthrough`,
the case that follows after is treated as part of the case that has the fallthrough.

That way the updater can jump into any case that applies to the users current DB-Version
and then slide down to the end. Resulting in every DB ending, at the current version,
no matter on what version it started.

## Code Example:

```go
package main

import "github.com/rocco-gossmann/go_sqliteutils"

// It is always good, to wrap your initialization into an extra function.
// these can become big.
func OpenMyDB() go_sqliteutils.DBRessource {
 
    const databaseFile = "mysqlitefile.db";
    const requiredDBVersion = 2;

    db, err := go_sqliteutils.Open(
        databaseFile,

        requiredDBVersion,

        // This is the Update-Handler
        func(tx *sql.Tx, isVersion, shouldBeVersion uint) error {
            
            switch(isVersion) {
            // Update starts here if the DB was freshly created
            // Creates the initial State of the DB
            case 0:
                if _, err := tx.Exec("CREATE TABLE tab_hello ( ... )"); err != nil {
                    return err
                }

                //we wan't slide through to next version until we reached the end
                fallthrough                 

            // Update starts here, if the User has already data for version 1 of the DB
            // Updates the State to Version 2
            case 1:
                if _, err := tx.Exec("CREATE TABLE tab_world ( ... )"); err != nil {
                    return err
                }


            }

            // return nil if all Updates ran without issues
            return nil
        }
    );

    if err != nil {
        panic(err)
    }

    return db
}


func main() {

    myDB := OpenMyDB();
    // Allways close your SQLite DB, when you are done
    // to prevent a lockup
    defer myDB.Close();
    
    // Check the next section to see, what you can do from here.

}
```

# Database Interaction

For any more complex Database Interaction, than what the Key=>Value Store offers,
There must always be a Transaction present. The Transaction then provides access
to the `mysql.Tx`s - Methods. 

```go
var tx *sql.Tx
var err error

tx, err = myDB.Begin()
if err != nil {
    panic(err)
}
defer tx.Rollback() // No matter what happens next,
// if no Commit is made, then Rollback everything, when the
// function ends

_, err = tx.Exec("INSERT INTO tab_hello(...) VALUES (?, ?)", "hallo", "hi");

if err != nil {
    panic(err)
    return
}

// After all the changes in this scope are done without issue 
// Commit TX to save your Changes to the DB
if err = tx.Commit(); err != nil {
    panic(err)
}
```


# The Key=>Value Store

this Lib comes with a Key => Value Store build in.
after you opened your DB you can use the following functions.

To save / change values:
```go
var err error = myDB.SetValue("myKey", "value for myKey")
```

To remove values:
```go
var err error = myDB.DropValue("myKey")
```

And finally to read a value:
```go
var output string
var err error

output, err = myDB.GetValue("myKey")
```

To Check if a value exists, you can use the Error-Message returned by `GetValue`

```go
if err == go_sqliteutils.NoResultsForKey {
    fmt.Println("no value registered for `myKey`")

} else if err != nil {
    fmt.Println(err);
    panic("PAAAANIC !!! * runs in circles *")
}

```

