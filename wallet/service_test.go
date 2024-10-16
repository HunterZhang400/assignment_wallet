package wallet

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mySQLMatcher = sqlmock.QueryMatcherFunc(func(expectedSQL, actualSQL string) error { //nolint:gochecknoglobals,gochecknoglobals
		if expectedSQL != actualSQL {
			return fmt.Errorf("expectedSQL: %s, actualSQL: %s", expectedSQL, actualSQL) //nolint:err113
		}
		return nil
	})
)

func TestUserExist(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(mySQLMatcher))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("select count(1) as amount from wallet_balance where user_id = $1").WithArgs("user-abc").
		WillReturnRows(sqlmock.NewRows([]string{"amount"}).AddRow(1))
	exist, err := CheckUserExist(sqlxDB, "user-abc")
	if err != nil {
		t.Errorf("error was not expected while checking if user exists: %v", err)
	}
	assert.Equal(t, true, exist)
}

func TestQueryBalance(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(mySQLMatcher))
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mock.ExpectQuery("select balance from wallet_balance where user_id = $1").WithArgs("user-abc").
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(100))
	exist, err := QueryBalance(sqlxDB, "user-abc")
	if err != nil {
		t.Errorf("error was not expected while checking if user exists: %v", err)
	}
	assert.Equal(t, int64(100), exist)
}
