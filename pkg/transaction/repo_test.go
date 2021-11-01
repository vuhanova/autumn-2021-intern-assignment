package transaction

import (
	"database/sql"
	"fmt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
	"time"
)

func TestGetUsersBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	//good query
	elemID := 1
	rows := sqlmock.NewRows([]string{"balance"})
	expect := []*User{&User{
		UserID:  elemID,
		Balance: 67.4,
	}}

	for _, item := range expect {
		rows = rows.AddRow(item.Balance)
	}

	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewRepository(db)

	tr, err := repo.GetUsersBalance(elemID, "")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(tr, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], tr)
		return
	}

}

func TestGetUsersBalanceErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	elemID := 1

	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(elemID).
		WillReturnError(fmt.Errorf("db_error"))

	repo := NewRepository(db)
	_, err = repo.GetUsersBalance(elemID, "")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestCreateUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	//good query
	elemID := 1
	rows := sqlmock.NewRows([]string{"id"})
	expect := []int{elemID}

	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	mock.
		ExpectQuery("INSERT INTO users").
		WithArgs(elemID, 0).
		WillReturnRows(rows)

	repo := NewRepository(db)

	err = repo.CreateUsers(elemID)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	// error

	mock.
		ExpectQuery("INSERT INTO users").
		WithArgs(elemID, 0).
		WillReturnError(fmt.Errorf("dont create such user"))

	err = repo.CreateUsers(elemID)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}

func TestAddMoney(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id"})
	elemID := 1
	expect := []int{elemID}
	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(55.3, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rows.AddRow(1)
	mock.
		ExpectQuery("INSERT INTO transaction").
		WithArgs(&elemID, nil, 55.3, time.Now().Format("2006-01-02 15:01")).
		WillReturnRows(rows)

	mock.ExpectCommit()
	//ok query
	err = repo.AddMoney(1, 55.3)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestAddMoneyError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewRepository(db)

	err = repo.AddMoney(1, -23.12)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	// expect begin
	mock.ExpectBegin().WillReturnError(fmt.Errorf("shahajskd"))

	err = repo.AddMoney(1, 23.12)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).WillReturnError(fmt.Errorf("no rows"))
	err = repo.AddMoney(1, 23.12)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	//sql.NoRows

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnError(sql.ErrNoRows)
	mock.
		ExpectQuery("INSERT INTO users").
		WithArgs(1, 0).
		WillReturnError(fmt.Errorf("dont create such user"))
	err = repo.AddMoney(1, 23.12)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//error update
	rows := sqlmock.NewRows([]string{"id"})
	elemID := 1
	expect := []int{elemID}
	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(55.3, 1).
		WillReturnError(fmt.Errorf("error"))

	err = repo.AddMoney(1, 55.3)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//error write transaction
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(55.3, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rows.AddRow(1)
	mock.
		ExpectQuery("INSERT INTO transaction").
		WithArgs(&elemID, nil, 55.3, time.Now().Format("2006-01-02 15:01")).
		WillReturnError(fmt.Errorf("don`t write transaction"))

	err = repo.AddMoney(1, 55.3)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithdrawMoney(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewRepository(db)

	rows := sqlmock.NewRows([]string{"id"})
	elemID := 1
	expect := []int{elemID}
	for _, item := range expect {
		rows = rows.AddRow(item)
	}
	//err = repo.AddMoney(1, 1000)
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(0.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rows.AddRow(1)
	mock.
		ExpectQuery("INSERT INTO transaction").
		WithArgs(nil, &elemID, 0.0, time.Now().Format("2006-01-02 15:01")).
		WillReturnRows(rows)

	mock.ExpectCommit()
	//ok query
	err = repo.WithdrawMoney(1, 0.0)

	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWithdrawMoneyError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()
	repo := NewRepository(db)

	mock.ExpectBegin()
	err = repo.WithdrawMoney(1, -223)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	rows := sqlmock.NewRows([]string{"id"})
	elemID := 1
	expect := []int{elemID}
	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	//not enough money
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	err = repo.WithdrawMoney(1, 500)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// expect begin
	mock.ExpectBegin().WillReturnError(fmt.Errorf("error"))

	err = repo.WithdrawMoney(1, 0.0)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// selec
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).WillReturnError(fmt.Errorf("no rows"))
	err = repo.WithdrawMoney(1, 0.0)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//error update

	for _, item := range expect {
		rows = rows.AddRow(item)
	}

	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(0.0, 1).
		WillReturnError(fmt.Errorf("error"))

	err = repo.WithdrawMoney(1, 0.0)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//error write transaction
	mock.ExpectBegin()
	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(1).
		WillReturnRows(rows)

	mock.
		ExpectExec("UPDATE users SET").
		WithArgs(0.0, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	rows.AddRow(1)
	mock.
		ExpectQuery("INSERT INTO transaction").
		WithArgs(nil, &elemID, 0.0, time.Now().Format("2006-01-02 15:01")).
		WillReturnError(fmt.Errorf("don`t write transaction"))

	err = repo.WithdrawMoney(1, 0.0)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
