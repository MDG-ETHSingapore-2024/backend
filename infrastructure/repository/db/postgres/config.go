package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ResultStatus int 
const (
	Success ResultStatus = iota
	Failure
)

type Repository struct {
	db *sql.DB
}

type Result struct {
	status ResultStatus
	rows   *sql.Rows
}

func (r *Result) Rows() *sql.Rows {
	return r.rows
}

func OpenDatabase(dbName string, connstr string) *Repository {
	if dbName == "" {
		dbName = "postgres"
	}

	db, err := sql.Open(dbName, connstr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(err.Error())
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Repository {
		db: db,
	}
}

func (r *Repository) CloseDatabase() {
	err := r.db.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (r *Repository) ExecuteQuery(stmt string, args ...any) *Result {
	rows, err := r.db.Query(stmt, args)
	if err != nil {
		log.Fatalln(err.Error())
		return &Result {
			status: Failure,
		}
	}

	return &Result {
		status: Success,
		rows: rows,
	}
}

func (r *Result) Status() ResultStatus {
	return r.status
}

func (r *Result) Rows() *sql.Rows {
	return r.rows
}