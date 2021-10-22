package Dao

import (
	"TinyORM/Dialect"
	mylog "TinyORM/Log"
	"database/sql"
	"fmt"
	"strings"
)
//和数据库进行交互
type session struct {
	db *sql.DB
	sql strings.Builder

	//用来替换sql语句中的占位符
	placeHolder []interface{}

	SqlGenerator Dialect.Dialect
	Table *Dialect.Schema

}

func NewSession(db *sql.DB)*session{
	return &session{
		db:db,
	}
}
func (s *session) Exec()(res sql.Result,err error){
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	if res , err = s.db.Exec(s.sql.String(),s.placeHolder...);err!=nil{
		mylog.Error(err)
	}
	return
}
func (s *session) QueryRow() (row *sql.Row) {
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	row = s.db.QueryRow(s.sql.String(), s.placeHolder...)
	return
}
func (s *session) QueryRows() (rows *sql.Rows, err error) {
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	if rows,err = s.db.Query(s.sql.String(),s.placeHolder...);err!=nil{
		mylog.Error(err)
	}
	return
}
func (s *session) SetSql(sql string,values ...interface{})*session {
	s.sql.WriteString(sql)
	s.placeHolder=nil
	s.placeHolder=append(s.placeHolder,values...)
	return s
}
func (s *session) SetSchema(object interface{}){
	s.Table=Dialect.ParseObect(object,s.SqlGenerator)
}
func (s *session) CreateTable() error {
	Sql, i := s.SqlGenerator.CreateTableSql(s.Table)
	s.SetSql(Sql,i...)
	_, err := s.Exec()
	return err
}
func (s *session) DropTable()error{
	Sql, i := s.SqlGenerator.DropTable(s.Table.Name)
	s.SetSql(Sql,i...)
	_, err := s.Exec()
	return err
}
func (s *session) HasTable() (bool,error) {
	Sql,i := s.SqlGenerator.TableExistSQL(s.Table.Name)
	s.SetSql(Sql,i...)
	row := s.QueryRow()
	var name string
	err := row.Scan(&name)
	fmt.Println(name)
	return name==s.Table.Name,err
}
func (s *session) RefTable()*Dialect.Schema {
	return s.Table
}
func (s *session) Insert(objects... interface{}){

}
func (s *session) insert(object interface{}) {

}

