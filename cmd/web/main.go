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
	snippets       models.SnippetModelInterface
	users          models.UserModelInterface
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
	hosts := strings.Split(*addr, ":")
	for i, v := range hosts {
		if v == "" {
			hosts = append(hosts[:i], hosts[i+1:]...)
		}
	}
	var host string
	if len(hosts) < 2 {
		host = "localhost"
	} else {
		host = strings.Split(*addr, ":")[0]
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
	// NOTE: This would pretty much eleminate CSRF attacks because
	// only cookies comming from same site will be used but this will mean
	// that links pointing to your site on other sites will not load the cookies
	// and the used won't be logged in even if he did already login, which can be bad UX.
	// sessionManager.Cookie.SameSite = http.SameSiteStrictMode

	// Ensures that cookies are sent over HTTPS only
	sessionManager.Cookie.Secure = true

	app := &application{
		errorLogger:    errorLogger,
		infoLogger:     infoLogger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
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
	infoLogger.Printf("Starting server on https://%s%s", host, *addr)
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

// TODO: Implement an HTTP server that redirects to the HTTPS one if you have time
// Generally this will be offloaded to an external party if you're using ACME
// certs

// func redirectToHTTPS(tlsPort string) {
// 	httpSrv := http.Server{
// 		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			host, _, err := net.SplitHostPort(r.Host)
// 			if err != nil {
// 				return
// 			}
// 			u := r.URL
// 			u.Host = net.JoinHostPort(host, tlsPort)
// 			u.Scheme = "https"
// 			http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
// 		}),
// 	}
// 	log.Fatal(httpSrv.ListenAndServe())
// }
