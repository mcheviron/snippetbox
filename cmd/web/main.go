package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"snippetbox/internal/models"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

// * To counter directory listing attacks
type neuteredFS struct {
	fs http.FileSystem
}

func (nfs neuteredFS) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		//* If index.html isn't found, it'll return, otherwise the index.html will be served
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}

func main() {
	addr := flag.String(
		"addr",
		":8000",
		"HTTP network address to bind to",
	)
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

	app := &application{
		errorLogger:   errorLogger,
		infoLogger:    infoLogger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
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
