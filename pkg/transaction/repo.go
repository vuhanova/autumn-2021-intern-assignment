package transaction

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"
)

type RepositoryItem struct {
	DB *sql.DB
}

type TransactionInterface interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewRepository(db *sql.DB) *RepositoryItem {
	return &RepositoryItem{
		DB: db,
	}
}

func (r *RepositoryItem) GetUsersBalance(userID int, currency string) (*User, error) {
	tr := &User{
		UserID:  userID,
		Balance: 0,
	}

	err := r.DB.QueryRow(`SELECT balance FROM users WHERE id = $1`, userID).Scan(&tr.Balance)
	if err != nil {
		return nil, err
	}

	if currency == "" || currency == "RUB" {
		return tr, nil
	}

	value, err := getCurrencyFromRub(currency)
	if err != nil {
		return nil, fmt.Errorf("didn`t convert currency: %s", currency)
	}

	tr.Balance = tr.Balance / value

	return tr, nil
}

func (r *RepositoryItem) CreateUsers(userID int) error {
	var id int
	balance := 0
	err := r.DB.QueryRow("INSERT INTO users (id, balance) VALUES ($1, $2) returning id", userID, balance).Scan(&id)
	if err != nil {
		return fmt.Errorf("dont create such user")
	}

	return nil
}

func (r *RepositoryItem) AddMoney(userID int, money float64) error {
	if money < 0 {
		return fmt.Errorf("negative amount")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	err = r.appendMoneyToUser(userID, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = writeTransaction(&userID, nil, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *RepositoryItem) appendMoneyToUser(userID int, money float64, db TransactionInterface) error {
	if money < 0 {
		return fmt.Errorf("negative amount")
	}
	_, err := r.GetUsersBalance(userID, "")
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		err = r.CreateUsers(userID)
		if err != nil {
			return err
		}
	}

	_, err = db.Exec("UPDATE users SET balance = balance + $1 WHERE id = $2", money, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryItem) getMoneyFromDB(userID int, money float64, db TransactionInterface) error {
	if money < 0 {
		return fmt.Errorf("negative amount")
	}

	tr, err := r.GetUsersBalance(userID, "")
	if err != nil {
		return err
	}

	if tr.Balance < money {
		return fmt.Errorf("not enough money")
	}

	_, err = db.Exec("UPDATE users SET balance = balance - $1 WHERE id = $2", money, userID)
	if err != nil {
		return err //failed to withdraw money
	}

	return nil
}

func (r *RepositoryItem) WithdrawMoney(userID int, money float64) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	err = r.getMoneyFromDB(userID, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = writeTransaction(nil, &userID, -money, r.DB)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (r *RepositoryItem) TransferMoney(fromUserID int, toUserID int, money float64) error {
	if money < 0 {
		return fmt.Errorf("negative amount")
	}
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	err = r.getMoneyFromDB(fromUserID, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = r.appendMoneyToUser(toUserID, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = writeTransaction(&toUserID, &fromUserID, money, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func writeTransaction(toID, fromID *int, money float64, db TransactionInterface) error {
	var id int
	created := time.Now().Format("2006-01-02 15:01")

	err := db.QueryRow("INSERT INTO transaction (to_id, from_id, money, created) VALUES ($1, $2, $3, $4) returning id",
		toID, fromID, money, created).Scan(&id)
	if err != nil {
		return fmt.Errorf("dont create transaction: %v", err)
	}

	return nil
}

func (r *RepositoryItem) GetTransaction(userID int, orderBy string) ([]*Transaction, error) {
	rows, err := r.DB.Query("SELECT to_id, from_id, money, created FROM transaction where to_id = $1 or from_id = $1",
		userID)
	if err != nil {
		return nil, err
	}

	info := make([]*Transaction, 0, 10)
	for rows.Next() {
		curr := &Transaction{}
		err = rows.Scan(&curr.ToID, &curr.FromID, &curr.Money, &curr.Created)
		if err != nil {
			return nil, err
		}
		info = append(info, curr)
	}

	lowOrderBy := strings.ToLower(orderBy)
	if lowOrderBy == "date" {
		sort.Slice(info[:], func(i, j int) bool {
			return info[i].Created.Before(info[j].Created)
		})
	} else if lowOrderBy == "money" {
		sort.Slice(info[:], func(i, j int) bool {
			return info[i].Money < info[j].Money
		})
	}

	return info, nil
}
