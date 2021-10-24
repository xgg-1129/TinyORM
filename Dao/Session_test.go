package Dao

import (
	mylog "TinyORM/Log"
	"fmt"
	"testing"
)

func TestSelect(t *testing.T) {
	engine, _:= NewEngine("sqlite3","test.db")
	defer engine.Close()
	s:=engine.NewSession()
	s.SetSchema(&User{})
	var res []User
	if err := s.Select(&res);err!=nil{
		mylog.Error(err)
	}
	for _,item := range res{
		fmt.Printf("%+v\n",item)
	}
}
func TestUpdate(t *testing.T) {
	engine, _:= NewEngine("sqlite3","test.db")
	defer engine.Close()
	s:=engine.NewSession()
	s.SetSchema(&User{})
	if err := s.Where("Name='xgg'").Delete();err!=nil{
		mylog.Error(err)
		t.Fail()
	}

}

type Account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *Account) BeforeInsert(s *session) error {
	mylog.Info("before inert", account)
	account.ID += 1000
	return nil
}

func (account *Account) AfterQuery(s *session) error {
	mylog.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	engine, _:= NewEngine("sqlite3","test.db")
	defer engine.Close()
	s:=engine.NewSession()
	s.SetSchema(&Account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_ = s.insert(&Account{1, "123456"})
	_ = s.insert(&Account{2, "qwerty"})


	u := make([]Account,0)
	err := s.Select(&u)
	if err != nil || u[0].ID != 1001 || u[0].Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
}