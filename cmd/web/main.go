package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"my.com/one/pkg/models/mysql"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	// root:root@tcp(172.17.0.2:3306)/test-db
	//dsn := flag.String("dsn", "web:pass@/tcp(localhost:3306)/snippetbox?parseTime=true", "MySQL data source name")
	// docker run --name=my-mysql -e MYSQL_ROOT_PASSWORD=root -d -p 33060:3306 mysql/mysql-server:latest
	// ALTER USER 'web'@'172.17.0.1' IDENTIFIED BY 'pass';
	// docker inspect my-mysql | grep '172.17'
	dsn := flag.String("dsn", "web:pass@tcp(localhost:33060)/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	// db.SetMaxOpenConns(50)
	// db.SetMaxIdleConns(5)

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	// app := &config.Application{
	// 	ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// }

	// mux := http.NewServeMux()
	// mux.HandleFunc("/", app.home)
	// // mux.Handle("/", handlers.Home(app))
	// mux.HandleFunc("/snippet", app.showSnippet)
	// mux.HandleFunc("/snippet/create", app.createSnippet)

	// fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		// Handler:  mux,
		Handler: app.routes(), // Call the new app.routes() method
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	// err := http.ListenAndServe(cfg.Addr, mux)
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
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
