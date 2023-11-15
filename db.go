package main

import (
	"context"
	"database/sql"
	"time"
)
import _ "github.com/go-sql-driver/mysql"

// var RoundRobinDb *RoundRobin

const dbTimeout = 100 * time.Millisecond

func init() {
	var err error
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:3306)/benchmark_sql")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1 * time.Second)
	db.SetConnMaxLifetime(60 * time.Second)

	cancel1, err := DbExec(db, CreateReadTable)
	if err != nil {
		panic(err)
	}
	defer cancel1()

	cancel2, err := DbExec(db, CreateWriteTable)
	if err != nil {
		panic(err)
	}
	defer cancel2()

	cancel3, err := DbExec(db, CreateSimpleData, "127.0.0.1", "/test-read", "www.google.com", time.Now())
	if err != nil {
		panic(err)
	}
	defer cancel3()

	//RoundRobinDb, err = New(db, 10)
	//if err != nil {
	//	panic(err)
	//}
}

func DbPrep(db *sql.DB, query string) (*sql.Stmt, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, cancel, err
	}

	return stmt, cancel, err
}

func DbExecStmt(stmt *sql.Stmt, args ...any) (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	_, err := stmt.ExecContext(ctx, args...)
	return cancel, err
}

func DbExec(db *sql.DB, query string, args ...any) (context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	_, err := db.ExecContext(ctx, query, args...)
	return cancel, err
}

func DbQuery(db *sql.DB, query string, args ...any) (*sql.Rows, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	rows, err := db.QueryContext(ctx, query, args...)

	return rows, cancel, err
}
