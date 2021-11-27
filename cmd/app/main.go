package main

import (
	"autumn-2021-intern-assignment/pkg/handlers"
	"autumn-2021-intern-assignment/pkg/transaction"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
)

func main() {
	/*
		- DB_PASSWORD=qwerty123
		      - PG_USER=postgres
		      - PG_DB=postgres
	*/
	fmt.Println("here2")

	name := os.Getenv("PG_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("PG_DB")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		name, password, dbname, "disable", "db", "5432")
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("no open bd:", err)
		return
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync() // flushes buffer, if any
	logger := zapLogger.Sugar()

	repo := transaction.NewRepository(db)
	handler := handlers.ItemsHandler{ItemRepo: repo, Logger: logger}
	r := mux.NewRouter()
	r.HandleFunc("/user", handler.GetBalanceFromUser)
	r.HandleFunc("/balance/add", handler.IncreaseBalance)
	r.HandleFunc("/balance/reduce", handler.DecreaseBalance)
	r.HandleFunc("/balance/transfer", handler.TransferBalance)
	r.HandleFunc("/info", handler.ListTransaction)

	addr := ":8000"

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
