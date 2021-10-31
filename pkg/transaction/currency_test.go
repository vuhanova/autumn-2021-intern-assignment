package transaction

import (
	"fmt"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"
)

func TestCurrency(t *testing.T) {
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
		Balance: 60,
	}}

	for _, item := range expect {
		rows = rows.AddRow(item.Balance)
	}

	mock.
		ExpectQuery("SELECT balance FROM users WHERE").
		WithArgs(elemID).
		WillReturnRows(rows)

	repo := NewRepository(db)

	tr, err := repo.GetUsersBalance(elemID, "GBP")
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}

	fmt.Println(tr.Balance)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	/*	if !reflect.DeepEqual(tr, expect[0]) {
			t.Errorf("results not match, want %v, have %v", expect[0], tr)
			return
		}
	*/

}

func TestCurrencyErrors(t *testing.T) {
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
	s, err := repo.GetUsersBalance(elemID, "mock11111111")
	fmt.Println(s)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

}
