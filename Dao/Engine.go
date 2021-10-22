package Dao

import (
	mylog "TinyORM/Log"
	"database/sql"
	"TinyORM/Dialect"
	"errors"
)

//Engine负责数据库的连接，关闭
type Engine struct {
	db *sql.DB
	dialect Dialect.Dialect
}
func NewEngine(driver, source string) (*Engine, error) {
	db, err := sql.Open(driver, source)
	if err!=nil{
		mylog.Error(err)
		return nil,err
	}
	//open未必会创建连接，通过ping能够确保连接生成
	if err = db.Ping();err!=nil{
		mylog.Error(err)
		return nil,err
	}

	dial:=Dialect.GetDialect(driver)
	if dial==nil {
		mylog.Error("the dialect dose not in DialectMap")
		return nil,errors.New("the dialect dose not in DialectMap")
	}
	e:=new(Engine)
	e.db=db
	e.dialect=dial
	mylog.Infof("Connect dataBase success: [source: %s ]",source)
	return e,nil
}
func (e *Engine) Close() error {
	err:=e.db.Close()
	if err!=nil{
		mylog.Error(err)
	}
	return err
}
func (e *Engine) NewSession()*session{
	s:=NewSession(e.db)
	s.SqlGenerator=e.dialect
	mylog.Info("Create New Session")
	return s
}