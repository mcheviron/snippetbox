package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"strings"

	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address to bind to")
	dsn := flag.String(
		"dsn",
		"web:Slavendral_1996@/snippetbox?parseTime=true",
		"MySQL data source name",
	)
	flag.Parse()

	if !strings.Contains(*addr, ":") {
		*addr = ":" + *addr
	}
	infoLogger := log.New(os.Stdout, "INFO\t",
		log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR\t",
		log.Ldate|log.Ltime|log.Lshortfile)
	db, err := openDB(*dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLogger:   errorLogger,
		infoLogger:    infoLogger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}
	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLogger,
		Handler:  app.routes(),
	}

	infoLogger.Printf("Starting server on %s", *addr)
	errorLogger.Fatal(server.ListenAndServe())
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		dsn,
	)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
