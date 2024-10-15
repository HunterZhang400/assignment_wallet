package wallet

import (
	"assigement_wallet/basedata"
	"assigement_wallet/pkg/db_util"
	"assigement_wallet/pkg/redis_util"
	"errors"
	"fmt"
	"math"
	"time"
)

const MoneyTypeOfDebit = 1
const MoneyTypeOfCredit = 2

var lockHoldMaxTime = time.Duration(30) * time.Second
var maxWaitLockTime = time.Duration(5) * time.Second

func QueryBalance(userID string) (int64, error) {
	db := db_util.GetDB()
	balance := int64(0)
	err := db.Get(&balance, "select balance from wallet_balance where user_id = $1", userID)
	return balance, err
}

func CheckUserExist(userID string) (bool, error) {
	db := db_util.GetDB()
	count := int64(0)
	err := db.Get(&count, "select count(1) from wallet_balance where user_id = $1", userID)
	return count > 0, err
}

func Deposit(userID string, money int64) error {
	if money <= 0 {
		return errors.New(basedata.InvalidAmount)
	}
	locker, err := redis_util.GetLockWithTimeout(getUserLockKey(userID), lockHoldMaxTime, maxWaitLockTime)
	if err != nil {
		return errors.New(basedata.UserBusy)
	}
	defer locker.UnLock()
	balance, err := QueryBalance(userID)
	if err != nil {
		return err
	}
	//prevent overfill
	room := math.MaxInt64 - balance
	if money > room {
		return errors.New("your account overfill")
	}
	db := db_util.GetDB()
	tx := db.MustBegin()
	_, err = tx.Exec("INSERT INTO wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values "+
		"($1, $2, $3, $4, $5, $6)", userID, MoneyTypeOfDebit, money, balance+money, time.Now(), "deposit")
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	_, err = tx.Exec("update wallet_balance set balance = balance + $1 where user_id = $2", money, userID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	err = tx.Commit()
	return err
}

func Withdraw(userID string, money int64) error {
	if money <= 0 {
		return errors.New(basedata.InvalidAmount)
	}
	locker, err := redis_util.GetLockWithTimeout(getUserLockKey(userID), lockHoldMaxTime, maxWaitLockTime)
	if err != nil {
		return errors.New(basedata.UserBusy)
	}
	defer locker.UnLock()
	balance, err := QueryBalance(userID)
	if err != nil {
		return err
	}
	if money > balance {
		return errors.New("your account balance insufficient")
	}
	db := db_util.GetDB()
	tx := db.MustBegin()
	_, err = tx.Exec("INSERT INTO wallet_detail (user_id, flow_type, amount,balance, occur_time, summary) values "+
		"($1, $2, $3, $4, $5, $6)", userID, MoneyTypeOfCredit, money, balance-money, time.Now(), "withdraw")
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	_, err = tx.Exec("update wallet_balance set balance = balance - $1 where user_id = $2", money, userID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	err = tx.Commit()
	return err
}

type TransactionDetail struct {
	ID        int64     `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	FlowType  int16     `db:"flow_type" json:"flow_type"`
	Amount    int64     `db:"amount" json:"amount"`
	Balance   int64     `db:"balance" json:"balance"`
	OccurTime time.Time `db:"occur_time" json:"occur_time"`
	Summary   string    `db:"summary" json:"summary"`
}

func QueryHistory(userID string, size int64) ([]TransactionDetail, error) {
	if size <= 0 || size > 1000 {
		return nil, errors.New("invalid query size")
	}
	db := db_util.GetDB()
	details := make([]TransactionDetail, 0)
	err := db.Select(&details, "select id,user_id,flow_type,amount,balance,occur_time,summary from wallet_detail"+
		" where user_id = $1 order by id limit $2", userID, size)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func Transfer(fromUserID, toUserID string, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	fromLocker, err := redis_util.GetLockWithTimeout(getUserLockKey(fromUserID), lockHoldMaxTime, maxWaitLockTime)
	if err != nil {
		return errors.New(basedata.UserBusy)
	}
	defer fromLocker.UnLock()
	toLocker, err := redis_util.GetLockWithTimeout(getUserLockKey(toUserID), lockHoldMaxTime, maxWaitLockTime)
	if err != nil {
		return errors.New(basedata.UserBusy)
	}
	defer toLocker.UnLock()

	fromBalance, err := QueryBalance(fromUserID)
	if err != nil {
		return err
	}
	if amount > fromBalance {
		return errors.New("your account is insufficient")
	}
	toBalance, err := QueryBalance(toUserID)
	if err != nil {
		return err
	}
	room := math.MaxInt64 - toBalance
	if amount > room {
		return errors.New("your account overfill")
	}
	db := db_util.GetDB()
	tx := db.MustBegin()
	_, err = tx.Exec("INSERT INTO wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values "+
		"($1, $2, $3, $4, $5, $6)", fromUserID, MoneyTypeOfCredit, amount, fromBalance-amount, time.Now(), fmt.Sprintf("transfer to %s", toUserID))
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	_, err = tx.Exec("INSERT INTO wallet_detail (user_id, flow_type, amount, balance, occur_time, summary) values "+
		"($1, $2, $3, $4, $5, $6)", toUserID, MoneyTypeOfDebit, amount, toBalance+amount, time.Now(), fmt.Sprintf("transfer from %s", fromUserID))
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	_, err = tx.Exec("update wallet_balance set balance = balance - $1 where user_id = $2", amount, fromUserID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	_, err = tx.Exec("update wallet_balance set balance = balance + $1 where user_id = $2", amount, toUserID)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(rollbackErr, err)
		}
		return err
	}
	err = tx.Commit()
	return nil
}

func getUserLockKey(userID string) string {
	return "redis:distribute:locker:key:" + userID
}
