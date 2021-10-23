package Dao

import (
	mylog "TinyORM/Log"
	_ "TinyORM/go-sqlite3-master/go-sqlite3-master"
	"fmt"
	"testing"
)

func TestEngine(t *testing.T) {
	engine, _:= NewEngine("sqlite3","xggdb")
	defer engine.Close()
	s:=engine.NewSession()
	s.SetSql("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.SetSql("CREATE TABLE User(Name text);").Exec()
	_, _ = s.SetSql("CREATE TABLE User(Name text);").Exec()
	result, _ := s.SetSql("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
type User struct {
	Name string `orm:"PRIMARY KEY"`
	Age  int
}
func TestSession(t *testing.T) {
	engine, _:= NewEngine("sqlite3","test.db")
	defer engine.Close()
	s:=engine.NewSession()
	s.SetSchema(&User{})
	s.CreateTable()
	if _,err := s.HasTable();err!=nil{
		mylog.Error(err)
	}
	u1:=User{
		Name: "zzb",
		Age:  21,
	}
	u2:=User{
		Name: "zx",
		Age:  11,
	}
	if err := s.Insert(u1,&u2);err!=nil{
		t.Fatal(err)
	}
}