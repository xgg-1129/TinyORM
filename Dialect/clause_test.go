package Dialect

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_genBindVars(t *testing.T) {
	type args struct {
		num int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{name: "1", args: struct{ num int }{num:2 }  },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := genBindVars(tt.args.num); got != tt.want {
				fmt.Println(got)
			}
		})
	}
}

func TestSelect(t *testing.T) {
	clause:=&Clause{
		sqls: map[Type]string{},
		Vars: map[Type][]interface{}{},
	}
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql)
	t.Log(vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}
