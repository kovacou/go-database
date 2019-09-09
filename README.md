# go-database

`go-database` is a library made for `mysql` which provides a set of extensions on top of [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) such as a querybuilder, profiler, context & transactions for performances. **This is not an ORM**.

## ➡ features

- Opening connection with environment variables :
    - Global variables `DATABASE_*`
    - Aliases variables for multiple connection `DATABASE_ALIAS_*`
- Opening connection with Environ
- A pool to manage all connections opened
- Transactions
- Profiling & Context
    - Log out queries as string (formated) ordered by execution time grouped by context to detect slow queries.  
    **Example :** You can profile an application that uses goroutines such as webserver.
- Query builder for complex query at the SQL layer
- **This is not an ORM (yet)**

## ➡ install

`go get github.com/kovacou/go-database`

## ➡ usage

Below is an example which shows some common use cases for go-database. 

```ini
# .env file configuration (default settings)

# You can use DSN
DATABASE_DSN=

# Or use specfic variables
DATABASE_DRIVER=mysql
DATABASE_HOST=172.18.0.1
DATABASE_USER=test
DATABASE_PASS=test
DATABASE_SCHEMA=dbtest
DATABASE_PORT=3306
```

```go
package main

import (
    "log"

    "github.com/kovacou/go-database"
    b "github.com/kovacou/go-database/builder"
)

type User struct {
    ID   int64
    Name string
}

func main() {
    // Opening a new connection based on environment variables.
    // By default, the connection is postponed until there is an action. 
    db, err := database.Open()
    if err != nil {
        log.Fatal(err)
    }

    // Example of builder
    s := b.Select{
        Table: "users",
        Where: b.ParseWhere("id < ?", 10)
    }

    println(s.String()) // SELECT * FROM users WHERE id < ?

    // Mapping through a map
    // (faster than structscan if you look for performance)
    out := []User{}
    db.SelectMap(s, func(values map[string]interface{}){
        out = append(out, User{
            ID:   values["id"].(int64),
            Name: string(values["name"]),
        })
    })

    // More faster than previous one.
    out2 := []User{}
    db.SelectSlice(s, func(values []interface{}) {
        out = append(out, User{
            ID:   values[0].(int64),
            Name: string(values[1]),
        })
    })
}
```

## ➡ opening a new connection
### **With environment variables**
### **With environ**

## ➡ closing all connections

```go
func main() {
    // Defering the close in your main ensure closing 
    // the connection before exiting the program.
    defer database.Close()

    // your code...
}
```

## ➡ transactions

## ➡ profiling & context
