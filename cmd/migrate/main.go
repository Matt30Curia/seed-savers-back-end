package main

import (
	"backend/seed-savers/config"
	"backend/seed-savers/db"
	"os"
	"log"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	sql "github.com/golang-migrate/migrate/v4/database/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.MySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := sql.WithInstance(db, &sql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "do" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
