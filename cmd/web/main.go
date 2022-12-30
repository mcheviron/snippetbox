package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"snippetbox/internal/models"
	"strings"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	snippets       *models.Snippets
	users          *models.Users
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Ensures that cookies are sent over HTTPS only
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLogger:    errorLogger,
		infoLogger:     infoLogger,
		snippets:       &models.Snippets{DB: db},
		users:          &models.Users{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	server := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLogger,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLogger.Printf("Starting server on %s", *addr)
	errorLogger.Fatal(server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
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
