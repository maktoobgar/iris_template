package repositories

import (
	"context"
	"database/sql"
	rawErrors "errors"
	"fmt"
	"reflect"
	"service/pkg/errors"
	"strings"
	"time"

	"github.com/georgysavva/scany/v2/sqlscan"
)

// Query generator structure
type Query struct {
	tableName string
	row       any
	query     string
	dbType    string
}

type QueryGenerator interface {
	// Sets the current row
	SetRowData(row any)
	// Sets the table name of the current row
	SetTableName(tableName string)
	// Sets current database type to generate right query
	SetDbType(dbType string)
	// Returns current db type
	GetDbType() string

	// Returns select fields for select operation
	GetSelectFields(prefix ...string) string
	// Returns insert fields for insert into operation
	GetInsertFields() (string, string)
	// Returns update fields for update operations
	GetUpdateFields() string
	// Formats all passed wheres in a string with `and` operator between them and `=` operator for key values
	GetWheres(where map[string]any) string
	// Formats all passed wheres in a string with `or` operator between them and `Like` operator for key values
	GetLikeWheres(where map[string]string) string

	// Generates a insert statement based on the row into the query builder
	InsertInto() QueryGenerator
	// Generates a insert statement of a slice
	InsertIntoMulti(data []QueryGenerator) QueryGenerator
	// Generates an update statement which updates current row
	UpdateMe() QueryGenerator
	// Generates a select statement and generates where with `GetWheres` function
	Select(optionalWhere ...map[string]any) QueryGenerator
	// Generates a select statement
	SelectWhere(where string) QueryGenerator
	// Adds order into the select query
	OrderBy(orderBy string, ascOrDesc string) QueryGenerator
	// Adds pagination into the select query
	Paginate(limit, whichPage int) QueryGenerator
	// Deletes current data from database
	DeleteMe() QueryGenerator
	// Generates an sql statement which will delete specific data with desired specifications
	Delete(optionalWhere ...map[string]any) QueryGenerator
	// Generates an sql statement which counts all the data with the same specifications
	SelectCount(optionalWhere ...map[string]any) QueryGenerator
	// Generates an update statement based on changed information in the row
	Update(optionalWhere ...map[string]any) QueryGenerator
	// Generates an update statement with desired specifications
	UpdateSpecific(set map[string]any, optionalWhere ...map[string]any) QueryGenerator
	// Generates a sql query which will get all information of the current row based on id of the row
	GetMe() QueryGenerator
	// An alias for GetMe
	SelectMe() QueryGenerator
	// Wraps around the query with BEGIN and END; to be atomic in database
	AtomicTransaction() QueryGenerator
	// Call this method from the table which has a foreign key field
	//
	// Like: User and Permission
	//
	// Every Permission has a field named user_id, so we call OneToMany
	// method from permission side like:
	//
	// # permission.OneToMany("users", 1)
	//
	// and then we get all permissions which user with id 1 has
	OneToMany(destinationTable string, destinationId int64) QueryGenerator
	// Call this method from every side of a ManyToMany relationship you like
	//
	// Like: User and Group
	//
	// There is a table at the middle of User and Group table named
	// whatever, just know the name and then if from User side ran the
	// method, you get which Users has the same Group with id 1.
	//
	// # user.ManyToMany("users_groups", "groups", 1)
	//
	// If you run it from Group side, you get which Groups a user
	// with id 1 has.
	//
	// # group.ManyToMany("users_groups", "users", 1)
	ManyToMany(middleTable, destinationTable string, destinationId int64) QueryGenerator
	// Call this method from every side of a ManyToMany relationship you like
	//
	// Like: User and Group
	//
	// There is a table at the middle of User and Group table named
	// whatever, just know the name and then if from User side ran the
	// method and you will remove all groups which that user has.
	//
	// # user.DeleteManyToMany("users_groups", "groups", []int{1, 2, 3, 4})
	//
	// In reverse It does exactly the same but in reverse which means that
	// it will remove all users attached to that group
	//
	// # group.DeleteManyToMany("users_groups", "users", []int{1, 2, 3, 4})
	DeleteManyToMany(middleTable, destinationTable string, destinationTableIdRange []int64) QueryGenerator
	// Call this method from every side of a ManyToMany relationship you like
	//
	// Like: User and Group
	//
	// There is a table at the middle of User and Group table named
	// whatever, just know the name and then if from User side ran the
	// method and you will append groups to current user
	//
	// # user.InsertManyToMany("users_groups", "groups", []int{1, 2, 3, 4})
	//
	// In reverse it does the same but reverse which it will append
	// Users to the current group
	//
	// # group.InsertManyToMany("users_groups", "users", []int{1, 2, 3, 4})
	InsertManyToMany(middleTable, destinationTable string, destinationTableIdRange []int64) QueryGenerator
	// If you want to handle the query your self, here you go
	RawQuery(input string) QueryGenerator
	// Retrusn the generated query and reset it
	Query() string

	// ExecQuery executes a query without returning any rows.
	//
	// # Used for insert/delete/update operations
	//
	// # Returns the last inserted/deleted/updated id
	//
	// # Returns 0 if couldn't give that id
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQuery(ctx context.Context, db *sql.DB) int64
	// ExecQueryRow executes a query that is expected to return one row.
	//
	// # Used for SelectOneRow Operations
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQueryRow(ctx context.Context, db *sql.DB)
	// ExecQueryRowErr executes a query that is expected to return one row.
	//
	// # Used for SelectOneRow Operations
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQueryRowErr(ctx context.Context, db *sql.DB) error
	// ExecQueryCount executes a query that is expected to return one row.
	//
	// # Used for SelectOneRow Operations
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQueryCount(ctx context.Context, db *sql.DB) int64
	// ExecQueryMulti executes a query that is expected to return multiple.
	//
	// # Used for SelectMultipleRows Operations
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQueryMulti(ctx context.Context, db *sql.DB, scanInto any)
	// ExecQueryMultiErr executes a query that is expected to return multiple.
	//
	// # Used for SelectMultipleRows Operations
	//
	// Query which is recorded inside will get removed after execution of this method.
	ExecQueryMultiErr(ctx context.Context, db *sql.DB, scanInto any) error
}

// Checks if passed input is a struct
func (q *Query) structCheck(data any) (reflect.Type, reflect.Value) {
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

// Check if passed input is a slice
func (q *Query) sliceCheck(data any) reflect.Value {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	for dataType.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
		dataType = dataValue.Type()
	}
	if dataType.Kind() != reflect.Slice {
		panic(rawErrors.New("repositories: data is not slice"))
	}
	return dataValue
}

// Formates values to be understandable by database
func (q *Query) formatValue(input any, nilIfEmpty ...bool) string {
	nilOnEmpty := false
	if len(nilIfEmpty) > 0 {
		nilOnEmpty = nilIfEmpty[0]
	}
	switch input := input.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		if fmt.Sprint(input) == "0" && nilOnEmpty {
			return "NULL"
		}
		return fmt.Sprint(input)
	case time.Time:
		return fmt.Sprintf("'%s'", input.UTC().Format(time.RFC3339Nano))
	case bool:
		if input {
			return "TRUE"
		} else {
			return "FALSE"
		}
	case string:
		if input == "" && nilOnEmpty {
			return "NULL"
		}
		if q.dbType == "postgres" || q.dbType == "sqlite3" {
			return fmt.Sprintf("'%s'", safePostgresSqlite(input))
		}
		return fmt.Sprintf("'%s'", input)
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("'%s'", input)
	}
}

func safePostgresSqlite(data string) string {
	var builder strings.Builder
	for _, char := range data {
		if char == '\'' {
			builder.WriteString("''")
		} else {
			builder.WriteRune(char)
		}
	}
	return builder.String()
}

func (q *Query) InsertInto() QueryGenerator {
	keys, values := q.GetInsertFields()
	q.query = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", q.tableName, keys, values)
	return q
}

func (q *Query) InsertIntoMulti(data []QueryGenerator) QueryGenerator {
	q.sliceCheck(data)
	keys, _ := q.GetInsertFields()
	values := ""
	for _, generator := range data {
		_, elementValues := generator.GetInsertFields()
		if values == "" {
			values = "(" + elementValues + ")"
			continue
		}
		values += ", (" + elementValues + ")"
	}
	q.query = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s", q.tableName, keys, values)
	return q
}

func (q *Query) UpdateMe() QueryGenerator {
	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()
	sets := q.GetUpdateFields()
	wheres := q.GetWheres(map[string]any{"id": id})
	q.query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", q.tableName, sets, wheres)
	return q
}

func (q *Query) Select(optionalWhere ...map[string]any) QueryGenerator {
	where := map[string]any{}
	if len(optionalWhere) != 0 {
		where = optionalWhere[0]
	}
	keys := q.GetSelectFields()
	wheres := q.GetWheres(where)
	if wheres != "" {
		q.query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", keys, q.tableName, wheres)
	} else {
		q.query = fmt.Sprintf("SELECT %s FROM %s", keys, q.tableName)
	}
	return q
}

func (q *Query) SelectWhere(wheres string) QueryGenerator {
	keys := q.GetSelectFields()
	if wheres != "" {
		q.query = fmt.Sprintf("SELECT %s FROM %s WHERE %s", keys, q.tableName, wheres)
	} else {
		q.query = fmt.Sprintf("SELECT %s FROM %s", keys, q.tableName)
	}
	return q
}

func (q *Query) OrderBy(orderBy string, ascOrDesc string) QueryGenerator {
	if q.query != "" {
		if ascOrDesc == "desc" {
			q.query += fmt.Sprintf(" ORDER BY %s DESC", orderBy)
		} else if ascOrDesc == "asc" {
			q.query += fmt.Sprintf(" ORDER BY %s ASC", orderBy)
		} else {
			panic(errors.New(errors.UnexpectedStatus, "InternalServerError", fmt.Sprintf("invalid ascOrDesc parameter: %s", ascOrDesc)))
		}
	} else {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", "no query to add order by into it"))
	}

	return q
}

func (q *Query) Paginate(limit, whichPage int) QueryGenerator {
	if q.query != "" {
		q.query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (whichPage-1)*limit)
	} else {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", "no query to paginate"))
	}

	return q
}

func (q *Query) DeleteMe() QueryGenerator {
	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()
	q.query = fmt.Sprintf("DELETE FROM %s WHERE id = %v", q.tableName, id)
	return q
}

func (q *Query) Delete(optionalWhere ...map[string]any) QueryGenerator {
	where := map[string]any{}
	if len(optionalWhere) != 0 {
		where = optionalWhere[0]
	}
	wheres := q.GetWheres(where)
	if wheres != "" {
		q.query = fmt.Sprintf("DELETE FROM %s WHERE %s", q.tableName, wheres)
	} else {
		q.query = fmt.Sprintf("DELETE FROM %s", q.tableName)
	}
	return q
}

func (q *Query) SelectCount(optionalWhere ...map[string]any) QueryGenerator {
	where := map[string]any{}
	if len(optionalWhere) != 0 {
		where = optionalWhere[0]
	}
	wheres := q.GetWheres(where)
	if wheres != "" {
		q.query = fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE %s", q.tableName, wheres)
	} else {
		q.query = fmt.Sprintf("SELECT COUNT(*) as count FROM %s", q.tableName)
	}
	return q
}

func (q *Query) Update(optionalWhere ...map[string]any) QueryGenerator {
	where := map[string]any{}
	if len(optionalWhere) != 0 {
		where = optionalWhere[0]
	}
	sets := q.GetUpdateFields()
	wheres := q.GetWheres(where)
	if wheres != "" {
		q.query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", q.tableName, sets, wheres)
	} else {
		q.query = fmt.Sprintf("UPDATE %s SET %s", q.tableName, sets)
	}
	return q
}

func (q *Query) UpdateSpecific(set map[string]any, optionalWhere ...map[string]any) QueryGenerator {
	where := map[string]any{}
	if len(optionalWhere) != 0 {
		where = optionalWhere[0]
	}
	sets := ""
	for key, value := range set {
		if sets == "" {
			sets = fmt.Sprintf("%s = %s", key, q.formatValue(value))
			continue
		}
		sets += fmt.Sprintf(", %s = %s", key, q.formatValue(value))
	}
	wheres := q.GetWheres(where)
	if wheres != "" {
		q.query = fmt.Sprintf("UPDATE %s SET %s WHERE %s", q.tableName, sets, wheres)
	} else {
		q.query = fmt.Sprintf("UPDATE %s SET %s", q.tableName, sets)
	}
	return q
}

func (q *Query) GetMe() QueryGenerator {
	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()
	q.query = fmt.Sprintf("SELECT * FROM %s WHERE id = %v", q.tableName, id)
	return q
}

func (q *Query) SelectMe() QueryGenerator {
	return q.GetMe()
}

func (q *Query) AtomicTransaction() QueryGenerator {
	q.query = fmt.Sprintf("BEGIN\n%s\nEND;", q.query)
	return q
}

func (q *Query) OneToMany(destinationTable string, destinationId int64) QueryGenerator {
	dataType, _ := q.structCheck(q.row)
	keys := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			if name == "-" || name == "" {
				continue
			}
			if keys == "" {
				keys = "main." + name
				continue
			}
			keys += ", main." + name
		}
	}
	q.query = fmt.Sprintf("SELECT DISTINCT %s FROM %s main JOIN %s destination ON main.%s_id = destination.id WHERE destination.id = %d", keys, q.tableName, destinationTable, destinationTable[:len(destinationTable)-1], destinationId)
	return q
}

func (q *Query) ManyToMany(middleTable, destinationTable string, destinationId int64) QueryGenerator {
	dataType, _ := q.structCheck(q.row)
	keys := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			if name == "-" || name == "" {
				continue
			}
			if keys == "" {
				keys = "main." + name
				continue
			}
			keys += ", main." + name
		}
	}
	q.query = fmt.Sprintf("SELECT DISTINCT %s FROM %s main JOIN %s middle ON main.id = middle.%s_id JOIN %s destination ON destination.id = middle.%s_id WHERE destination.id = %d", keys, q.tableName, middleTable, q.tableName[:len(q.tableName)-1], destinationTable, destinationTable[:len(destinationTable)-1], destinationId)
	return q
}

func (q *Query) DeleteManyToMany(middleTable, destinationTable string, destinationTableIdRange []int64) QueryGenerator {
	ins := ""
	for _, destinationId := range destinationTableIdRange {
		if ins == "" {
			ins = fmt.Sprint(destinationId)
			continue
		}
		ins += ", " + fmt.Sprint(destinationId)
	}

	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()
	q.query = fmt.Sprintf("DELETE FROM %s WHERE %s_id = %v AND %s_id IN (%s)", middleTable, q.tableName[:len(q.tableName)-1], id, destinationTable[:len(destinationTable)-1], ins)
	return q
}

func (q *Query) InsertManyToMany(middleTable, destinationTable string, destinationTableIdRange []int64) QueryGenerator {
	value := reflect.ValueOf(q.row)
	id := value.Elem().FieldByName("Id").Interface()

	values := ""
	for _, destinationId := range destinationTableIdRange {
		if values == "" {
			values = fmt.Sprintf("(%v, %d)", id, destinationId)
			continue
		}
		values += fmt.Sprintf(" (%v, %d)", id, destinationId)
	}

	q.query = fmt.Sprintf("INSERT INTO %s (%s_id, %s_id) VALUES %s", middleTable, q.tableName[:len(q.tableName)-1], destinationTable[:len(destinationTable)-1], values)
	return q
}

func (q *Query) RawQuery(input string) QueryGenerator {
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

func (q *Query) SetRowData(row any) {
	q.row = row
}

func (q *Query) SetTableName(tableName string) {
	q.tableName = tableName
}

func (q *Query) SetDbType(dbType string) {
	q.dbType = dbType
}

func (q *Query) GetDbType() string {
	return q.dbType
}

func (q *Query) GetSelectFields(prefix ...string) string {
	pre := ""
	if len(prefix) > 0 {
		pre = prefix[0]
	}
	dataType, _ := q.structCheck(q.row)
	keys := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			if name == "-" || name == "" {
				continue
			}
			if keys == "" {
				if pre != "" {
					keys = pre + "." + name
				} else {
					keys = name
				}
				continue
			}
			if pre != "" {
				keys += ", " + pre + "." + name
			} else {
				keys += ", " + name
			}
		}
	}
	return keys
}

func (q *Query) GetInsertFields() (string, string) {
	dataType, dataValue := q.structCheck(q.row)
	keys := ""
	values := ""
	for _, f := range reflect.VisibleFields(dataType) {
		if f.IsExported() {
			name := f.Tag.Get("db")
			fieldName := f.Name
			if name == "-" || name == "" || f.Tag.Get("skipInsert") == "+" {
				continue
			}
			fieldValue := dataValue.FieldByName(fieldName)
			for fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			var value any = nil
			if fieldValue.IsValid() {
				value = fieldValue.Interface()
			}
			if keys == "" {
				keys = name
				values = q.formatValue(value, f.Tag.Get("nilOnEmpty") == "+")
				continue
			}
			keys += ", " + name
			values += ", " + q.formatValue(value, f.Tag.Get("nilOnEmpty") == "+")
		}
	}

	return keys, values
}

func (q *Query) GetUpdateFields() string {
	dataType, dataValue := q.structCheck(q.row)
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
					sets = fmt.Sprintf("%s = %s", name, q.formatValue(value, f.Tag.Get("nilOnEmpty") == "+"))
					continue
				}
				sets += fmt.Sprintf(", %s = %s", name, q.formatValue(value, f.Tag.Get("nilOnEmpty") == "+"))
			}
		}
	}
	return sets
}

func (q *Query) GetWheres(where map[string]any) string {
	wheres := ""
	for key, value := range where {
		if wheres == "" {
			value := q.formatValue(value)
			if value == "NULL" {
				wheres = fmt.Sprintf("%s IS %s", key, value)
			} else {
				wheres = fmt.Sprintf("%s = %s", key, value)
			}
			continue
		}
		value := q.formatValue(value)
		if value == "NULL" {
			wheres += fmt.Sprintf(" AND %s IS %s", key, value)
		} else {
			wheres += fmt.Sprintf(" AND %s = %s", key, value)
		}
	}

	return wheres
}

func (q *Query) GetLikeWheres(where map[string]string) string {
	wheres := ""
	like := ""
	if q.dbType == "sqlite3" {
		like = "LIKE"
	} else if q.dbType == "postgres" {
		like = "ILIKE"
	}
	for key, value := range where {
		if wheres == "" {
			wheres = fmt.Sprintf("%s %s %s", key, like, q.formatValue("%%"+value+"%%"))
			continue
		}
		wheres += fmt.Sprintf(" OR %s %s %s", key, like, q.formatValue("%%"+value+"%%"))
	}

	return wheres
}

func (q *Query) ExecQuery(ctx context.Context, db *sql.DB) int64 {
	query := q.Query()
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", err.Error()+" Query: "+query))
	}

	if rows, err := result.RowsAffected(); err == nil && rows > 0 {
		if lastId, err := result.LastInsertId(); err == nil && lastId > 0 {
			reflect.ValueOf(q.row).Elem().FieldByName("Id").Set(reflect.ValueOf(lastId))
			return lastId
		}
	}
	return 0
}

func (q *Query) ExecQueryRow(ctx context.Context, db *sql.DB) {
	query := q.Query()
	err := sqlscan.Get(ctx, db, q.row, query)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", err.Error()+" Query: "+query))
	}
}

func (q *Query) ExecQueryRowErr(ctx context.Context, db *sql.DB) error {
	query := q.Query()
	err := sqlscan.Get(ctx, db, q.row, query)
	return err
}

func (q *Query) ExecQueryCount(ctx context.Context, db *sql.DB) int64 {
	query := q.Query()
	count := int64(-1)
	err := db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", err.Error()+" Query: "+query))
	}
	return count
}

func (q *Query) ExecQueryMulti(ctx context.Context, db *sql.DB, scanInto any) {
	query := q.Query()
	err := sqlscan.Select(ctx, db, scanInto, query)
	if err != nil {
		panic(errors.New(errors.UnexpectedStatus, "InternalServerError", err.Error()+" Query: "+query))
	}
}

func (q *Query) ExecQueryMultiErr(ctx context.Context, db *sql.DB, scanInto any) error {
	query := q.Query()
	err := sqlscan.Select(ctx, db, scanInto, query)
	return err
}

// Returns a new QueryGenerator
func NewQueryGenerator(tableName string) QueryGenerator {
	return &Query{
		tableName: tableName,
	}
}
