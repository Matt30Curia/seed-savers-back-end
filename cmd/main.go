package main

import (
	"backend/seed-savers/cmd/api"
	"backend/seed-savers/config"
	"backend/seed-savers/db"
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {

	db, err := db.MySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil{
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewServer(":3000", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB){
	err := db.Ping()
	if err != nil{
		log.Fatal(err)
	}
	log.Println("succes conneted to db")
}