package Dao

import (
	"database/sql"
	"strings"
	mylog "TinyORM/Log"
)

type session struct {
	db *sql.DB
	sql strings.Builder

	//用来替换sql语句中的占位符
	placeHolder []interface{}
}

func NewSession(db *sql.DB)*session{
	return &session{
		db:db,
	}
}
func (s *session) Exec()(res sql.Result,err error){
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	if res , err = s.db.Exec(s.sql.String(),s.placeHolder);err!=nil{
		mylog.Error(err)
	}
	return
}
func (s *session) QueryRow() (row *sql.Row) {
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	row = s.db.QueryRow(s.sql.String(), s.placeHolder)
	return
}
func (s *session) QueryRows() (rows *sql.Rows, err error) {
	defer s.sql.Reset()
	mylog.Info(s.sql.String(), s.placeHolder)
	if rows,err = s.db.Query(s.sql.String(),s.placeHolder);err!=nil{
		mylog.Error(err)
	}
	return
}

