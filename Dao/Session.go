package Dao

import (
	"TinyORM/Dialect"
	mylog "TinyORM/Log"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)
//和数据库进行交互
type session struct {
	db *sql.DB
	sql strings.Builder

	//执行事务引擎
	tx *sql.Tx

	//用来替换sql语句中的占位符
	placeHolder []interface{}

	SqlGenerator Dialect.Dialect
	Table *Dialect.Schema
	clause   Dialect.Clause
}

func NewSession(db *sql.DB)*session{
	return &session{
		db:db,
	}
}
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}
func (s *session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}
func (s *session) Exec()(res sql.Result,err error){
	defer s.Clear()
	mylog.Info(s.sql.String(), s.placeHolder)
	if res , err = s.DB().Exec(s.sql.String(),s.placeHolder...);err!=nil{
		mylog.Error(err)
	}
	return
}
func (s *session) QueryRow() (row *sql.Row) {
	defer s.Clear()
	mylog.Info(s.sql.String(), s.placeHolder)
	row = s.DB().QueryRow(s.sql.String(), s.placeHolder...)
	return
}
func (s *session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	mylog.Info(s.sql.String(), s.placeHolder)
	if rows,err = s.DB().Query(s.sql.String(),s.placeHolder...);err!=nil{
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
func (s *session) Clear() {
	s.sql.Reset()
	s.placeHolder = nil
	s.clause = Dialect.Clause{}
}
func (s *session) Insert(objects ...interface{}) error{
	for _,object :=range objects{
		if err := s.insert(object);err!=nil{
			return err
		}
	}
	return nil
}
func (s *session) insert(object interface{})error{
	s.CallMethod(BeforeInsert,object)
	recordValues:=s.Table.RecordValues(object)
	s.clause.Set(Dialect.INSERT, s.Table.Name, s.Table.FieldsName)
	s.clause.Set(Dialect.VALUES, recordValues)
	sql, vars := s.clause.Build(Dialect.INSERT, Dialect.VALUES)
	s.SetSql(sql,vars...)
	_, err := s.Exec()
	return err
}
func (s *session) Select(slice interface{})error{
	//
	s.CallMethod(BeforeQuery,nil)
	destSlice:=reflect.Indirect(reflect.ValueOf(slice))
	ElemType:=destSlice.Type().Elem()

	s.clause.Set(Dialect.SELECT, s.Table.Name, s.Table.FieldsName)
	sql, vars := s.clause.Build(Dialect.SELECT, Dialect.WHERE, Dialect.ORDERBY, Dialect.LIMIT)
	s.SetSql(sql,vars...)
	rows, err := s.QueryRows()
	if err!=nil{
		return err
	}
	for rows.Next(){
		dest:=reflect.New(ElemType).Elem()
		//how to 构造dest。。。。
		var values []interface{}
		for _, name := range s.Table.FieldsName {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err = rows.Scan(values...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery,dest.Addr().Interface())
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
func (s *session) Update( m map[string]interface{})error{
	s.clause.Set(Dialect.UPDATE,s.Table.Name,m)
	sql, vars := s.clause.Build(Dialect.UPDATE, Dialect.WHERE)
	s.SetSql(sql,vars...)
	_, err := s.Exec()
	return err
}
func (s *session) Delete() error {
	s.clause.Set(Dialect.DELETE,s.Table.Name)
	sql, vars := s.clause.Build(Dialect.DELETE, Dialect.WHERE)
	s.SetSql(sql,vars...)
	_, err := s.Exec()
	return err
}
func (s *session) Where(str string)*session{
	s.clause.Set(Dialect.WHERE,str)
	return s
}

func (s *session) Begin() (err error) {
	mylog.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		mylog.Error(err)
		return
	}
	return
}

func (s *session) Commit() (err error) {
	mylog.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		mylog.Error(err)
	}
	return
}

func (s *session) Rollback() (err error) {
	mylog.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		mylog.Error(err)
	}
	return
}

