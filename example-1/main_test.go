package main

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub connection", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = recordStats(db, 3, 2)
	if err != nil {
		t.Errorf("error in updating record stats: %s", err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("error in meeting the expectations: %s", err.Error())
	}

}

//fail case
func TestShouldRollBackUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening stub connection", err.Error())
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).
		WillReturnError(fmt.Errorf("unable to update product_viewers"))
	mock.ExpectRollback()

	err = recordStats(db, 3, 2)
	if err == nil {
		t.Errorf("expeting an error in updating records stats: but there was none")
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("error in meeting expectations: %s", err.Error())
	}
}
