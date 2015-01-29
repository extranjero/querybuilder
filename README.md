## querybuilder
Go SQL Select builder

```go
package main
import (
  "github.com/extranjero/querybuilder"
  _ "github.com/lib/pq"
)

func main(){
  db, err := sql.Open("postgres", "postgres://hrkb:hrkb@/hrkb")

  if err != nil {
    panic("Unable to connect to database")
  }
  
  qb := QB(db)
  
  qb.From("table").Select([]string{"id","name"}).Where("id=?",1).Limit(10).Offset(5).OrderBy("name DESC").Sort("ASC")
  
  qb.Sql()
  // returns
  // "Select id,name FROM table WHERE id=? LIMIT 10 OFFSET 5 ORDER BY name DESC ASC"
  
  qb.Query()
  //return rows from db
  
}

```
