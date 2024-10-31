package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func MySQLStorage(conf mysql.Config) ( *sql.DB, error) {
	db, err := sql.Open("mysql", conf.FormatDSN())

	if err != nil{
		log.Fatal(err)
	}

	return db, err
}