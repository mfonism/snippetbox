package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":666", "HTTP network address")

	dbPass := os.Getenv("dbPass")
	defaultDSN := fmt.Sprintf("web:%s@/snippetbox?parseTime=true", dbPass)
	dsn := flag.String("dsn", defaultDSN, "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
	}

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	defer db.Close()

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: app.errorLog,
		Handler: app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
