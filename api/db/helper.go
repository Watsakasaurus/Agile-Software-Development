package db

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"xorm.io/builder"

	"medicare-api/types"
)

// SelectQuery is just a shortcut to create a xorm/builder
func (c *Client) SelectQuery(cols string) *builder.Builder {
	return builder.Select(cols)
}

// Builder used for building queries
func (c *Client) Builder() *builder.Builder {
	return builder.Dialect(builder.POSTGRES)
}

// Select executes the query and scans the results i to the passed interface
func (c *Client) Select(
	result interface{},
	query *builder.Builder) error {

	sql, args, err := query.ToSQL()
	if err != nil {
		log.Debugf("ERROR Getting Rows: %s", err)
		return err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return err
	}

	return sqlx.Select(c.ex, result, sql, args...)
}

// Get helper for getting one record
func (c *Client) Get(result interface{},
	query *builder.Builder) error {

	sql, args, err := prepareQuery(query)
	if err != nil {
		log.Debugf("ERROR Preparing query: %s", err)
		return err
	}
	return sqlx.Get(c.ex, result, sql, args...)
}

func prepareQuery(query *builder.Builder) (string, []interface{}, error) {
	sql, args, err := query.ToSQL()
	if err != nil {
		return "", nil, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return "", nil, err
	}

	return sql, args, nil
}

// SelectWithCount executes the query and scans the results into the passed interface
// Additionally it fetches the rowcount and adds Pagination
func (c *Client) SelectWithCount(
	result interface{},
	query *builder.Builder,
	pageNumber, perPage int) (int, error) {

	// Check that the result slice sent in is an empty slice
	//		created via ```name = New(type)``` or ``` name := type{}```
	//		has a pointer to an empty array
	// rather than a nil slice
	//		created via ```var name type```
	//		has no pointer
	// as a nil array is json encoded as null rather than []
	if reflect.Indirect(reflect.ValueOf(result)).IsNil() {
		return 0, errors.New("Result list must be empty and not nil")
	}

	// copy the struct so we can replace the select fields
	countQuery := &builder.Builder{}
	*countQuery = *query
	countQuery.Select("count(*)").OrderBy("")

	sql, args, err := countQuery.ToSQL()
	if err != nil {
		return 0, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return 0, err
	}
	log.Debugf("SQL: %s : %+v", sql, args)

	total := 0
	err = sqlx.Get(c.ex, &total, sql, args...)
	if err != nil {
		log.Debugf("ERROR Getting Totals: %s", err)
		return 0, err
	}

	sql, args, err = query.ToSQL()
	if err != nil {
		log.Debugf("ERROR Getting Rows: %s", err)
		return 0, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return 0, err
	}

	sql = c.paginate(sql, pageNumber, perPage)
	return total, sqlx.Select(c.ex, result, sql, args...)
}

// SelectWithCountSQL executes the query (raw sql)  and scans the results into the passed interface
// Additionally it fetches the rowcount and adds Pagination
func (c *Client) SelectWithCountSQL(result interface{}, query string, params []interface{},
	orderBy string, pageNumber, perPage int) (int, error) {

	var total int
	err := sqlx.Get(c.ex, &total, "SELECT count(*) FROM ("+query+") AS orig", params...) // nolint: gosec
	if err != nil {
		return 0, err
	}

	if orderBy != "" {
		query += " ORDER BY " + orderBy
	}
	query = c.paginate(query, pageNumber, perPage)
	err = sqlx.Select(c.ex, result, query, params...)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (c *Client) paginate(query string, pageNumber int, perPage int) string {
	if pageNumber > 0 && perPage > 0 {
		// TODO Change to use placeholders, go-xorm/builder should support it soon
		query += fmt.Sprintf(
			" OFFSET %d LIMIT %d", perPage*(pageNumber-1), perPage)
	}
	return query
}

// Exec executes a query against the database
func (c *Client) Exec(query *builder.Builder) (sql.Result, error) {
	sql, args, err := query.ToSQL()
	if err != nil {
		return nil, err
	}
	sql, err = builder.ConvertPlaceholder(sql, "$")
	if err != nil {
		return nil, err
	}

	log.Debugf("EXEC: %s %#v", sql, args)
	return c.ex.Exec(sql, args...)
}

// ILike defines ilike condition
type ILike [2]string

var _ builder.Cond = ILike{"", ""}

// WriteTo write SQL to Writer
func (ilike ILike) WriteTo(w builder.Writer) error {
	if _, err := fmt.Fprintf(w, "%s ILIKE ?", ilike[0]); err != nil {
		return err
	}
	// FIXME: if use other regular express, this will be failed. but for compatible, keep this
	if ilike[1][0] == '%' || ilike[1][len(ilike[1])-1] == '%' {
		w.Append(ilike[1])
	} else {
		w.Append("%" + ilike[1] + "%")
	}
	return nil
}

// And implements And with other conditions
func (ilike ILike) And(conds ...builder.Cond) builder.Cond {
	return builder.And(ilike, builder.And(conds...))
}

// Or implements Or with other conditions
func (ilike ILike) Or(conds ...builder.Cond) builder.Cond {
	return builder.Or(ilike, builder.Or(conds...))
}

// IsValid tests if this condition is valid
func (ilike ILike) IsValid() bool {
	return len(ilike[0]) > 0 && len(ilike[1]) > 0
}

const (
	foreignKeyViolation = "23503"
	uniqueViolation     = "23505"
)

func (c *Client) transformError(err error) *types.Error {
	if err == nil {
		return nil
	}
	// Check key violation error
	if pqErr, ok := err.(*pq.Error); ok {
		if pqErr.Code == foreignKeyViolation {
			return types.ValidationError(fmt.Sprintf("%s: %s", pqErr.Message, pqErr.Detail))
		} else if pqErr.Code == uniqueViolation {
			return types.DuplicateError(fmt.Sprintf("%s: %s", pqErr.Message, pqErr.Detail))
		}
	}

	return types.DatabaseError(err)
}
