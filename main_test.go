package main

import (
	"testing"
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Where struct {
	condType int // 1 = Where 2 = OR 3 = AND
	cond string
	param interface{}
}

type Request struct {
	columns []string
	from string
	wheres []Where
	limit int
	offset int
	orderBy string
	sort string
}

func init(){

	var err error
	db, err = sql.Open("postgres", "postgres://hrkb:hrkb@/hrkb")

	if err != nil {
		panic("Unable to connect to remote database")
	}
}

func (r Request) prepare() *Query {
	qb := QB(db)

	qb.From(r.from).Select(r.columns)
	if len(r.wheres) > 0 {
		for _, where := range r.wheres {
			switch t := where.condType; t {
			case 1:
				qb.Where(where.cond, where.param)
			case 2:
				qb.OrWhere(where.cond, where.param)
			case 3:
				qb.AndWhere(where.cond, where.param)
			}
		}
	}

	if r.limit != 0 {
		qb.Limit(r.limit)
	}

	if r.offset != 0 {
		qb.Limit(r.offset)
	}

	if r.orderBy != "" {
		qb.OrderBy(r.orderBy)
	}

	if r.sort != "" {
		qb.Sort(r.sort)
	}

	//returning QB object
	return qb
}

func (r Request) Sql() string {
	qb := r.prepare()
	return qb.Sql()
}

func (r Request) Query() (*sql.Rows, error) {
	qb := r.prepare()
	return qb.Query()
}

func TestMain(t *testing.T) {
	var cases = []struct {
		in Request
		out string
	}{
		{Request{ columns: []string{"id","login"}, from: "users", wheres: []Where{ Where{1,"id=?",3} } }, "SELECT id,login FROM users WHERE id=$1 OFFSET 0 ASC"},
	}

	for _, c := range cases {
		sql := c.in.Sql()
		query, err := c.in.Query()

		if sql != c.out {
			t.Errorf("Sql forming test: DB => %q, want %q", sql, c.out)
		}

		if err == nil {
			t.Errorf("Query to DB => %q, result %q", sql, query )
		}
	}
}
