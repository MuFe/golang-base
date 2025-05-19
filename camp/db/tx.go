package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mufe/golang-base/camp/xlog"
)

type Tx struct {
	tx    *sql.Tx
	Print bool
}

func (s *Tx) Query(sql string, args ...interface{}) (rows *sql.Rows, err error) {
	rows, err = s.tx.Query(sql, args...)
	if err != nil {
		xlog.ErrorP(err)
	}
	return rows, err
}

func (s *Tx) QueryRow(sql string, args ...interface{}) (result *sql.Row) {
	result = s.tx.QueryRow(sql, args...)
	return result
}

func (s *Tx) Exec(sql string, args ...interface{}) (result sql.Result, err error) {
	result, err = s.tx.Exec(sql, args...)
	return result, err
}

func (s *Tx) GetTx() *sql.Tx {
	return s.tx
}
