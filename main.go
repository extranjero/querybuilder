package main

import (
	"database/sql"
	"strings"
	"strconv"
)

type Query struct {
	from string
	columns []string
	wheres []string
	params []interface{}
	limit int
	offset int
	orderBy string
	sort string
	db *sql.DB
}

func (q *Query) From(from string) *Query {
	q.from = from
	return q
}

func (q *Query) Select(columns []string) *Query {
	q.columns = columns
	return q
}

func (q *Query) Where(condition string, args ...interface{}) *Query {
	q.where("", condition, args)
	return q
}

func (q *Query) AndWhere(condition string, args ...interface{}) *Query {
	q.where("AND ", condition, args)
	return q
}

func (q *Query) OrWhere(condition string, args ...interface{}) *Query {
	q.where("OR ", condition, args)
	return q
}

func (q *Query) where(sep, condition string, args ...interface{}) {
	var count int = 1

	if len( q.params ) > 0 {
		count = len( q.params ) + 1
	}
	i := 0
	for strings.Index(condition, "?") != -1 {
		condition = strings.Replace(condition,"?", "$"+strconv.Itoa(count), 1)
		count++
		if len(args) > i {
			q.params = append(q.params, args[i] )
			i++
		} else {
			panic("Few Arguments passed")
		}
	}
	q.wheres = append(q.wheres, sep + condition)
}

func (q *Query) Limit(limit int) *Query {
	q.limit = limit
	return q
}

func (q *Query) Offset(offset int) *Query {
	q.offset = offset
	return q
}

func (q *Query) OrderBy(orderBy string) *Query {
	q.orderBy = orderBy
	return q
}

func (q *Query) Sort(sort string) *Query {
	q.sort = sort
	return q
}

func (q *Query) Sql() string {
	cols := strings.Join( q.columns, ",")
	// Adding Columns & table name
	sql := "SELECT "+cols+" FROM "+q.from

	where := strings.Join( q.wheres, " " )
	// Adding where conditions
	sql += " WHERE "+where
	// If limit set then add Limit
	if q.limit != 0 {
		sql += " LIMIT "+strconv.Itoa(q.limit)
	}
	// Adding Offset
	sql += " OFFSET "+strconv.Itoa(q.offset)
	// if OrderBy exists then add it
	if q.orderBy != "" {
		sql += " ORDER BY "+q.orderBy
	}
	// Adding sort
	sql += " "+q.sort
	//returning sql
	return sql
}

func (q *Query) Query() (*sql.Rows, error) {
	return q.db.Query(q.Sql(), q.params...);
}

func QB(db *sql.DB) *Query {
	return &Query{ offset:0, sort: "ASC", db: db }
}
