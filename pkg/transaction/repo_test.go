package transaction

import (
	"fmt"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"reflect"
	"testing"
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
