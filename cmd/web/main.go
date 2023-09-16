package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
}

func main() {
	addr := flag.String("addr", ":666", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	app := &application {
		errorLog: errorLog,
		infoLog: infoLog,
	}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: app.errorLog,
		Handler: app.routes(),
	}

	app.infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
