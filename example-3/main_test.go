package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

//pass test case
func TestCancelOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error in creating new sql mock: %v", err.Error())
	}
	defer db.Close()

	o := Order{Id: 1, Value: 24.34, ReservedFee: 1.5, Status: ORDER_PENDING}
	u := User{Id: 1, Username: "Akram", Balance: 12.5}

	returningColumns := []string{"o_id", "o_reserved_fee", "o_status", "o_value", "u_balance", "u_id", " u_username"}
	expectedRows := sqlmock.NewRows(returningColumns).
		AddRow(o.Id, o.ReservedFee, o.Status, o.Value, u.Balance, u.Id, u.Username)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM orders AS o	INNER JOIN users AS u (.+)	FOR UPDATE").WithArgs(o.Id).WillReturnRows(expectedRows)
	mock.ExpectPrepare("UPDATE users SET ").ExpectExec().WithArgs(o.Value+o.ReservedFee, u.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	o.Status = ORDER_CANCELLED
	mock.ExpectPrepare("UPDATE orders SET ").ExpectExec().WithArgs(o.Status, o.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	cancelOrder(o.Id, db)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("failed to verify expectations of mock %s", err.Error())
	}
}
