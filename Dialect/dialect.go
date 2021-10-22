package Dialect

import (
	mylog "TinyORM/Log"
	"reflect"
)

//dialect适配器模型,主要的作用就是生成每个数据库标准的sql语句，然后赋予session去执行


var	dialectMap = map[string]Dialect{}

type Dialect interface {
	DataTypeOf(p reflect.Value)string
	TableExistSQL(name string)(string, []interface{})
	CreateTableSql(table *Schema)(string,[]interface{})
	DropTable(name string)(string,[]interface{})
}

func RegisterDialect(name string,d Dialect){
	dialectMap[name]=d
}
func GetDialect(name string) Dialect {
	if res, ok := dialectMap[name];ok{
		mylog.Info("the Name of Dialect does not register")
		return res
	}
	return nil
}
func init() {
	dialectMap["sqlite3"]=&Sqlite3{}
}


