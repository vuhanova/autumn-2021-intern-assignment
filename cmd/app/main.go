package main

import (
	"autumn-2021-intern-assignment/pkg/handlers"
	"autumn-2021-intern-assignment/pkg/transaction"
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	var name, password, dbname string
	flag.StringVar(&name, "user", "", "The name of user")
	flag.StringVar(&password, "password", "", "password")
	flag.StringVar(&dbname, "db", "", "The name of database")
	flag.Parse()
	connStr := "user=" + name + " password=" + password + " dbname=" + dbname +
		" sslmode=disable host=localhost port=5432"
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

	addr := ":9000"

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}
