package Dialect

import (

	"testing"
)

// schema_test.go
type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}



func TestParse(t *testing.T) {
	var TestDial  =  &Sqlite3{}
	s := ParseObect(&User{}, TestDial)
	if s.Name != "User" || len(s.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if s.FieldsMap["Name"].Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
}