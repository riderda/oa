package dbs

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	USERNAME = "root"
	PASSWORD = ""
	HOST = "localhost"
	PORT = "3306"
	DATABASE = "demo"
	CHARSET = "utf8"

)

func GetConnect()(*sqlx.DB,error){
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",USERNAME,PASSWORD,HOST,PORT,DATABASE,CHARSET)
	Db ,err := sqlx.Open("mysql",dbDSN)
	if err != nil {
		return Db, err
	}
	Db.SetConnMaxLifetime(500)
	Db.SetMaxIdleConns(20)
	Db.SetMaxOpenConns(100)
	return Db, nil
}
