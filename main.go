package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(localhost:3306)/benchmark_sql")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(2 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)

	go sqlStat(db)

	http.HandleFunc("/static", handlerStatic)
	http.HandleFunc("/read", handlerRead(db))
	http.HandleFunc("/write", handlerWrite(db))
	log.Println("Listen 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func sqlStat(db *sql.DB) {
	for {
		<-time.After(1 * time.Second)
		printDbStat(db.Stats())
	}
}

func printDbStat(stat sql.DBStats) {
	log.Println("db stat")
	fmt.Println("MaxOpenConnections", stat.MaxOpenConnections)
	fmt.Println("OpenConnections", stat.OpenConnections)
	fmt.Println("InUse", stat.InUse)
	fmt.Println("Idle", stat.Idle)
	fmt.Println("WaitCount", stat.WaitCount)
	fmt.Println("WaitDuration", stat.WaitDuration)
	fmt.Println("MaxIdleClosed", stat.MaxIdleClosed)
	fmt.Println("MaxIdleTimeClosed", stat.MaxIdleTimeClosed)
	fmt.Println("MaxLifetimeClosed", stat.MaxLifetimeClosed)
	fmt.Println("==================================================")
}

func handlerStatic(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "ok")
	if err != nil {
		log.Println(err)
	}
}

func handlerWrite(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		stmt, prepCancel, err := DbPrep(db, "INSERT INTO benchmark_write(ip, created_at) VALUES(?, ?)")
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer prepCancel()
		defer stmt.Close()

		cancel, err := DbExecStmt(stmt, ip, time.Now())
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer cancel()

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "ok")
	}
}

func handlerRead(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, cancelRow, err := DbQuery(db, "SELECT id, ip, url, redirect_to, created_at FROM benchmark_read")

		defer cancelRow()
		defer rows.Close()

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var ar []Read

		for rows.Next() {
			var read Read
			var t string
			err := rows.Scan(&read.Id, &read.Ip, &read.Url, &read.RedirectTo, &t)
			if err != nil {
				log.Println(err)
				continue
			}

			read.CreateAt, err = time.Parse("2006-01-02 15:04:05", t)
			if err != nil {
				log.Println(err)
				continue
			}

			ar = append(ar, read)
		}

		res, err := json.Marshal(ar)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = fmt.Fprintf(w, "%s", res)
		if err != nil {
			log.Println(err)
		}
	}
}
