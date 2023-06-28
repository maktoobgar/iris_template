package repositories

import (
	"context"
	"database/sql"
	rawErrors "errors"
	"fmt"
	"reflect"
	"service/pkg/errors"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
)

type Query struct {
	tableName string
	row       any
	query     string
}

func structCheck(data any) (reflect.Type, reflect.Value) {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	for dataType.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
		dataType = dataValue.Type()
	}
	if dataType.Kind() != reflect.Struct {
		panic(rawErrors.New("repositories: data is not struct"))
	}
	return dataType, dataValue
}

func formatValue(input any) string {
	switch input := input.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		return fmt.Sprint(input)
	case time.Time:
		return fmt.Sprintf("'%s'", input.Format(time.RFC3339Nano))
	case bool:
		if input {
			return "TRUE"
		} else {
			return "FALSE"
		}
	default:
		return fmt.Sprintf("'%s'", input)
	}
}

func (q *Query) InsertInto() *Query {
	dataType, dataValue := structCheck(q.row)
	query := fmt.Sprintf("INSERT INTO %s ", q.tableName)
	keys := ""
	values := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			fieldName := f.Name
			if name == "-" || name == "" || f.Tag.Get("skipInsert") == "+" {
				continue
			}
			value := dataValue.FieldByName(fieldName).Interface()
			if value != nil {
				if keys == "" {
					keys = name
					values = formatValue(value)
					continue
				}
				keys += ", " + name
				values += ", " + formatValue(value)
			}
		}
	}
	q.query = query + fmt.Sprintf("(%s) VALUES(%s)", keys, values)
	return q
}

func (q *Query) Select(keyValues map[string]any) *Query {
	dataType, _ := structCheck(q.row)
	keys := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			if name == "-" || name == "" {
				continue
			}
			if keys == "" {
				keys = name
				continue
			}
			keys += ", " + name
		}
	}
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	q.query = fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s", keys, q.tableName, wheres)
	return q
}

func (q *Query) Delete(keyValues map[string]any) *Query {
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	q.query = fmt.Sprintf("DELETE FROM %s WHERE %s", q.tableName, wheres)
	return q
}

func (q *Query) SelectCount(keyValues map[string]any) *Query {
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE", q.tableName)
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	q.query = fmt.Sprintf("%s %s", query, wheres)
	return q
}

func (q *Query) Update(keyValues map[string]any) *Query {
	query := fmt.Sprintf("UPDATE %s SET ", q.tableName)
	dataType, dataValue := structCheck(q.row)
	sets := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			fieldName := f.Name
			if name == "-" || name == "" || f.Tag.Get("skipUpdate") == "+" {
				continue
			}
			value := dataValue.FieldByName(fieldName).Interface()
			if value != nil {
				if sets == "" {
					sets = fmt.Sprintf("%s = %s", name, formatValue(value))
					continue
				}
				sets += fmt.Sprintf(", %s = %s", name, formatValue(value))
			}
		}
	}
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	query += sets + " WHERE " + wheres
	q.query = query
	return q
}

func (q *Query) UpdateSpecific(updates map[string]any, keyValues map[string]any) *Query {
	query := fmt.Sprintf("UPDATE %s SET ", q.tableName)
	sets := ""
	for key, value := range updates {
		if sets == "" {
			sets = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		sets += fmt.Sprintf(", %s = %s", key, formatValue(value))
	}
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	query += sets + " WHERE " + wheres
	q.query = query
	return q
}

func (q *Query) GetMe() *Query {
	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()
	q.query = fmt.Sprintf("SELECT * FROM %s WHERE id = %v", q.tableName, id)
	return q
}

func (q *Query) RawQuery(input string) *Query {
	q.query = input
	return q
}

func (q *Query) Query() string {
	output := ""
	if len(q.query) == 0 {
		output = ";"
	} else if q.query[len(q.query)-1] != ';' {
		output = q.query + ";"
	} else {
		output = q.query
	}
	q.query = ""
	return output
}

func (q *Query) SetTableName(tableName string) {
	q.tableName = tableName
}

func (q *Query) SetRowData(row any) {
	q.row = row
}

// ExecContext executes a query without returning any rows.
//
// # Used for insert/delete/update operations
//
// # Returns the last inserted/deleted/updated id
//
// # Returns 0 if couldn't give that id
//
// Query which is recorded inside will get removed after execution of this method.
func (q *Query) ExecContext(ctx context.Context, db *sql.DB) int64 {
	query := q.Query()
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, errors.Resend, "InternalServerError", err.Error(), nil))
	}

	if rows, err := result.RowsAffected(); err == nil && rows > 0 {
		if lastId, err := result.LastInsertId(); err == nil && lastId > 0 {
			reflect.ValueOf(q.row).Elem().FieldByName("Id").Set(reflect.ValueOf(lastId))
			return lastId
		}
	}
	return 0
}

// QueryRowContext executes a query that is expected to return at most one row.
//
// # Used for SelectOneRow Operations
//
// Query which is recorded inside will get removed after execution of this method.
func (q *Query) QueryRowContext(ctx context.Context, db *sql.DB) {
	query := q.Query()
	err := sqlscan.Get(ctx, db, q.row, query)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, errors.Resend, "InternalServerError", err.Error(), nil))
	}
}

// QueryRowContext executes a query that is expected to return at most one row.
//
// # Used for SelectOneRow Operations
//
// Query which is recorded inside will get removed after execution of this method.
func (q *Query) QueryRowContextError(ctx context.Context, db *sql.DB) error {
	query := q.Query()
	err := sqlscan.Get(ctx, db, q.row, query)
	return err
}

func NewQuery(tableName string) Query {
	return Query{
		tableName: tableName,
	}
}
