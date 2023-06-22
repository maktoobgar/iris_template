package repositories

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func structCheck(data any) (reflect.Type, reflect.Value) {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	for dataType.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
		dataType = dataValue.Type()
	}
	if dataType.Kind() != reflect.Struct {
		panic(errors.New("repositories: data is not struct"))
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

func InsertInto(table string, data any) string {
	dataType, dataValue := structCheck(data)
	query := fmt.Sprintf("INSERT INTO %s ", table)
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
	return query + fmt.Sprintf("(%s) VALUES(%s);", keys, values)
}

func Select(table string, data any, keyValues map[string]any) string {
	dataType, _ := structCheck(data)
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
	return fmt.Sprintf("SELECT DISTINCT %s FROM %s WHERE %s;", keys, table, wheres)
}

func Delete(table string, keyValues map[string]any) string {
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	return fmt.Sprintf("DELETE FROM %s WHERE %s;", table, wheres)
}

func SelectCount(table string, keyValues map[string]any) string {
	query := fmt.Sprintf("SELECT COUNT(*) as count FROM %s WHERE", table)
	wheres := ""
	for key, value := range keyValues {
		if wheres == "" {
			wheres = fmt.Sprintf("%s = %s", key, formatValue(value))
			continue
		}
		wheres += fmt.Sprintf(" AND %s = %s", key, formatValue(value))
	}
	return fmt.Sprintf("%s %s;", query, wheres)
}

func Update(table string, data any, keyValues map[string]any) string {
	query := fmt.Sprintf("UPDATE %s SET ", table)
	dataType, dataValue := structCheck(data)
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
	query += sets + " WHERE " + wheres + ";"
	return query
}

func UpdateSpecific(table string, updates map[string]any, keyValues map[string]any) string {
	query := fmt.Sprintf("UPDATE %s SET ", table)
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
	query += sets + " WHERE " + wheres + ";"
	return query
}
