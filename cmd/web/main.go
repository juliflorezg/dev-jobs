package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/juliflorezg/dev-jobs/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	jobPosts       models.JobPostModelInterface
	users          models.UserModelInterface
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		// content, err := ioutil.ReadFile(os.Getenv("DATABASE_URL_FILE"))
		// content, err := os.ReadFile(os.Getenv("DATABASE_URL_FILE"))
		content, err := os.ReadFile(os.Getenv("SECRET_URL_PATH"))
		if err != nil {
			log.Fatal(err)
		}
		databaseUrl = string(content)
	}

	fmt.Println("database URL:::", databaseUrl)

	addr := flag.String("port", ":8080", "HTTP Network Address")

	dsn := flag.String("dsn", databaseUrl, "MySQL data source name")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		templateCache:  templateCache,
		jobPosts:       &models.JobPostModel{DB: db},
		users:          &models.UserModel{DB: db},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * (time.Second),
		WriteTimeout: 5 * (time.Second),
	}

	logger.Info("Starting server, listening in port", "port", *addr)

	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(dsn string) (*sql.DB, error) {

	var db *sql.DB
	var err error

	for i := 0; i < 10; i++ { // Try up to 10 times
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return db, nil
			}
		}
		log.Println("Failed to connect to database. Retrying...", err)
		time.Sleep(2 * time.Second)
	}

	if db != nil {
		db.Close()
	}
	return nil, err

}
