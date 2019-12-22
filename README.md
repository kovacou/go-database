# go-database

Personal project.  

`go-database` is a library made for `mysql` which provides a set of extensions on top of [`jmoiron/sqlx`](https://github.com/jmoiron/sqlx) such as a querybuilder, profiler, context & transactions for performances. **This is not an ORM**.  
I plan to support `pgsql` in the future.

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
    ID   uint64
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
    db.SelectMap(&s, func(values map[string]interface{}){
        out = append(out, User{
            ID:   values["id"].(int64),
            Name: string(values["name"].([]byte)),
        })
    })

    // More faster than previous one.
    out2 := []User{}
    db.SelectSlice(&s, func(values []interface{}) {
        out2 = append(out2, User{
            ID:   values[0].(int64),
            Name: string(values[1].([]byte)),
        })
    })

    // Using raw query
    out3 := []User{}
    db.QuerySlice(builder.NewQuery("SELECT * FROM users WHERE id < ?", 10), func(values []interface{}) {
        out3 = append(out3, User{
            // ...
        })
    })

    // Using upsert
    i := builder.Insert{
        Table: "users",
        Values: builder.H{
            "id": 1,
            "name": "John",
        },
        OnUpdateKeys: builder.Keys{"name"},
    }
    db.Exec(&i)
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

Support of transactions.

```go
// tx, err := db.Tx() 
tx, err := db.Tx(sql.LevelSerializable)
if err != nil {
    panic(err)
}

// use tx to run some requests.

tx.Commit()
tx.Rollback()

// You can't use tx anymore, else an error will occur.
```

## ➡ profiling & context

## ➡ statements

### **Select**

```go
s := builder.Select{
    Table: "users",
    Columns: builder.ParseColumns("id", "name"),
    Where: builder.ParseWhere("id > ?", 1),
    OrderBy: builder.ParseOrderBy("name ASC"),
}
```

#### Map
The columns can be read by key name.

```go
// Parse 1 row only. 
{
    out := User{}
    n, err := db.SelectMapRow(&s, func(v map[string]interface{}) {   
        // If there is more than 1 row, an error occur.
        // `n` return 0 or 1.
        outRow.ID = v["id"].(int64)
        outRow.Name = string(v["name"].([]byte))
    })
}

// Parse multiple rows.
{
    out := []User{}
    n, err := db.SelectMap(&s, func(v map[string]interface{}) {
        // `n` contains the number of rows returned.
        out = append(out, User{
            ID:   v["id"].(int64),
            Name: string(v["name"].([]byte)),
        })
    })
}
```

#### Slice

The columns can be read by indexes from the Column clause (same order).  
**Note:** Slice is faster than Map. Prefer use Slice when the columns have always the same order.
```go
// Parse 1 row only
{
    out := User{}
    n, err := db.SelectSliceRow(&s, func(v []interface{}){
        // If there is more than1 row, an error occur.
        // `n` return 0 or 1
        out.ID = v[0].(int64)
        out.Name = string(v[1].([]byte))
    })
}

// Parse multiple rows
{
    out := []User{}
    n, err := db.SelectSlice(&s, func(v []interface{}){
        out = append(out, User{
            ID:   v[0].(int64),
            Name: string(v[1].([]byte)),
        })
    })
}
```

### **Exec**

#### Insert
```go
// Example of Insert
i := builder.Insert{
    Table:      "users",
    IgnoreMode: false, // False by default
    Values:     builder.H{
        "name": "John",
    },
}

println(i.String()) // INSERT INTO users(name) VALUES(?) 

r, err := db.Exec(&i)
```

#### Upsert

```go
// Example of Insert with upsert mode.
i := builder.Insert{
    Table:        "users",
    Values:       builder.H{
        "id": 15,
        "name": "John",
    },
    OnUpdateKeys: builder.Keys{"name"},
}

println(i.String()) // INSERT INTO users(id, name) VALUES(?, ?) 
                    // ON DUPLICATE KEY UPDATE name = VALUES(name) 

r, err := db.Exec(&i)
```

#### Update

```go
// Example of Update
u := builder.Update{
    Table:  "users",
    Values: builder.H{
        "name": "John",
    },
    Where: builder.ParseWhere("id = ?", 1) // Used to initiate the value (if needed)
    // To initiate empty value :
    // builder.NewWhere()
}

// You can also use where like following :
// u.Where.And("id = ?", 1)

println(u.String()) // UPDATE users SET name = ? WHERE id = ?

r, err := db.Exec(&u)
```

#### Delete

```go
// Example of Delete
d := builder.Delete{
    Table: "users",
    Where: builder.ParseWhere("id = ?", 1) // Used to initiate the value (if needed)
    // To initiate empty value :
    // builder.NewWhere()
}

// You can also use where like following :
// d.Where.And("id = ?", 1)

println(d.String()) // DELETE FROM users WHERE id = ?

r, err := db.Exec(&d)
```