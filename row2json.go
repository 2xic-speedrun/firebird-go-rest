// Modified from https://stackoverflow.com/questions/42774467/how-to-convert-sql-rows-to-typed-json-in-golang
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
)

type jsonNullInt64 struct {
	sql.NullInt64
}

func (v jsonNullInt64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Int64)
}

type jsonNullFloat64 struct {
	sql.NullFloat64
}

func (v jsonNullFloat64) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Float64)
}

type jsonNullString struct {
	sql.NullString
}

func (v jsonNullString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(v.String)
}

var jsonNullInt64Type = reflect.TypeOf(jsonNullInt64{})
var jsonNullFloat64Type = reflect.TypeOf(jsonNullFloat64{})

var jsonNullStringType = reflect.TypeOf(jsonNullString{})

var jsonNullTimeType = reflect.TypeOf(sql.NullTime{})
var nullInt64Type = reflect.TypeOf(sql.NullInt64{})
var nullFloat64Type = reflect.TypeOf(sql.NullFloat64{})
var nullStringType = reflect.TypeOf(sql.NullString{})

func UNUSED(x ...interface{}) {}

// SQLToJSON takes an SQL result and converts it to a nice JSON form. It also
// handles possibly-null values nicely. See https://stackoverflow.com/a/52572145/265521
func SQLToJSON(rows *sql.Rows) (map[string][]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Column error: %v", err)
	}

	values, err := mapColumnTypes(rows)
	if err != nil {
		panic(err)
	}

	data := make(map[string][]interface{})

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			panic("Failed to scan values")
		}
		for i, v := range values {
			data[columns[i]] = append(data[columns[i]], v)
		}
	}

	return data, nil
}

func mapColumnTypes(rows *sql.Rows) ([]interface{}, error) {
	tt, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("Column type error: %v", err)
	}

	types := make([]reflect.Type, len(tt))
	for i, tp := range tt {
		st := tp.ScanType()
		if st == nil {
			return nil, fmt.Errorf("Scantype is null for column: %v", err)
		}

		// TODO: Don't do it like this sir, just for testing
		if st.Name() == "string" {
			types[i] = jsonNullStringType
		} else if st.Name() == "float64" {
			types[i] = jsonNullStringType
		} else if st.Name() == "int16" {
			types[i] = jsonNullStringType
		} else if st.Name() == "int32" {
			types[i] = jsonNullStringType
		} else if st.Name() == "Decimal" {
			types[i] = jsonNullFloat64Type
		} else if st.Name() == "Time" {
			types[i] = jsonNullTimeType
		} else {
			switch st {
			case nullStringType:
				types[i] = jsonNullStringType
			case nullInt64Type:
				types[i] = jsonNullInt64Type
			case nullFloat64Type:
				types[i] = jsonNullFloat64Type
			default:
				types[i] = st
			}
		}
	}

	values := make([]interface{}, len(tt))

	for i := range values {
		values[i] = reflect.New(types[i]).Interface()
	}

	return values, nil
}
