package Dialect

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Sqlite3 struct {}


/*
CREATE TABLE student(
   ID INT PRIMARY KEY     ,
   NAME           TEXT    NOT NULL,
   AGE            INT     NOT NULL,
   ADDRESS        CHAR(50),
   FEES         REAL
);
*/
func (s *Sqlite3) CreateTableSql(table *Schema) (string, []interface{}) {
	var columns []string
	for _,item := range table.Fields{
		columns=append(columns,fmt.Sprintf("%s %s %s",item.Name,item.SqlType,item.Tag))
	}
	desc:=strings.Join(columns,",")
	sql:=fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)
	return sql,nil
}

func (s *Sqlite3) DropTable(name string) (string, []interface{}) {
	sql:=fmt.Sprintf("DROP TABLE IF EXISTS %s", name)
	return sql,nil
}

func (s *Sqlite3) DataTypeOf(typ reflect.Value) string {
	//kind是底层类型
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid Sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (s *Sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT Name FROM sqlite_master WHERE type='table' and Name = ?", args
}

