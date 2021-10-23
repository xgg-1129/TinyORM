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