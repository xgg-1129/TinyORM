package Dialect

import (
	"fmt"
	"strings"
)

//https://www.runoob.com/sqlite/sqlite-update.html

type generator func(values ...interface{}) (string, []interface{})


var generators map[Type]generator
type  Type int
const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

func _insert(values ...interface{})(string, []interface{}) {
	// INSERT INTO $tableName ($fields),第一个参数是表名，第二个是[]fields
	tableName:=values[0].(string)
	str := strings.Join(values[1].([]string), ",")
	sql:=fmt.Sprintf("INSERT INTO %s (%s)", tableName, str)
	return sql,[]interface{}{}
}
//获取num个占位符 2  -- ?,?
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{})(string, []interface{}) {

	var sql strings.Builder
	var vars []interface{}

	sql.WriteString("VALUES")
	v:=values[0].([]interface{})
	sql.WriteString(fmt.Sprintf("(%v)",genBindVars(len(v))))
	vars=append(vars,v...)
	return sql.String(),vars
}
func _select(values ...interface{}) (string, []interface{}) {
	// SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}
func _limit(values ...interface{}) (string, []interface{}) {
	// LIMIT $num
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}


type Clause struct {
	sqls map[Type]string
	Vars map[Type][]interface{}
}
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sqls == nil {
		c.sqls = make(map[Type]string)
		c.Vars = make(map[Type][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sqls[name] = sql
	c.Vars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sqls[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.Vars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}
func _delete(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	return _select(values[0], []string{"count(*)"})
}
