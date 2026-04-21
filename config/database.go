package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() *sql.DB {
	dbUser := "root"
	dbPass := "root"
	dbName := "financial_record"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbDriver := "mysql"

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true&loc=Local"
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		log.Fatalln("Gagal koneksi ke database", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln("Gagal ping ke database", err)
	}
	//errPing:= db.Ping()
	//if errPing!= nil{}
	log.Println("Berhasil konek ke database")
	return db
}

//parseTime=true, untuk insert data waktu ke database
//loc=local, untuk mengatur zona waktu sesuai sistem
