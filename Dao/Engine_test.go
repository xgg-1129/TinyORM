package Dao

import (
	mylog "TinyORM/Log"
	_ "TinyORM/go-sqlite3-master/go-sqlite3-master"
	"errors"
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
func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "test.db")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}
func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})

}
func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.SetSchema(&User{})
	s.DropTable()
	_, err := engine.Transaction(func(s *session) (result interface{}, err error) {
		s.SetSchema(&User{})
		s.CreateTable()
		err = s.Insert(&User{"Tom", 18})
		return nil, errors.New("Error")
	})
	if err == nil  {
		t.Fatal("failed to rollback")
	}
}
func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	s.SetSchema(&User{})
	s.DropTable()
	_, err := engine.Transaction(func(s *session) (result interface{}, err error) {
		s.SetSchema(&User{})
		s.CreateTable()
		err = s.Insert(&User{"Tom", 18})
		return
	})
	u := make([]User,0)
	s.Select(&u)
	if err != nil || u[0].Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
func TestSession(t *testing.T) {
	engine, _:= NewEngine("sqlite3","test.DB")
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