package balance

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"finstat/internal/entity"
	"finstat/internal/errs"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo *Repo) AddMoneyToUser(userUid string, amount int64) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	err = repo.addmoney(tx, userUid, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) addmoney(tx *sqlx.Tx, userUid string, amount int64) error {
	const addmoneyQuerty = `
		UPDATE t_user
		SET balance = balance + $2
		WHERE user_uid = $1
		`

	_, err := tx.Exec(addmoneyQuerty, userUid, amount)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repo) SendMoneyFromUser(fromUserUid, toUserUid string, amount int64) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	balance, err := repo.getUserFromId(fromUserUid)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if balance-amount < 0 {
		return errs.ErrUserIsUnableToPay
	}

	_, err = repo.SubmoneyFromUser(tx, fromUserUid, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = repo.addMoneyToUser(tx, toUserUid, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	_, err = repo.CreateTransaction(tx, fromUserUid, toUserUid, amount)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repo) SubmoneyFromUser(tx *sqlx.Tx, fromUserUid string, amount int64) (balance int64, err error) {
	const query = `
			UPDATE t_user
			SET balance = balance - $1
			WHERE user_uid = $2
			RETURNING balance
	`
	err = tx.QueryRowx(query, amount, fromUserUid).Scan(&balance)
	return
}

func (repo *Repo) addMoneyToUser(tx *sqlx.Tx, fromUserUid string, amount int64) (balance int64, err error) {
	const query = `
			UPDATE t_user
			SET balance = balance + $1
			WHERE user_uid = $2
			RETURNING balance
	`
	err = tx.QueryRowx(query, amount, fromUserUid).Scan(&balance)
	return
}

func (repo *Repo) CreateTransaction(tx *sqlx.Tx, fromUserUid string, toUserUid string, amount int64) (transactionID string,
	err error,
) {
	var t entity.Transaction

	t.TransactionID = uuid.NewString()
	t.UserUid = fromUserUid
	t.Amount = amount
	t.ToUserUid = toUserUid
	t.FromUserUid = fromUserUid

	loc, _ := time.LoadLocation("Europe/Moscow")
	t.TransactionDate = time.Now().In(loc)

	transactionID, err = repo.createTransaction(tx, t)

	return
}

func (repo *Repo) createTransaction(tx *sqlx.Tx, t entity.Transaction) (transactionID string, err error) {
	const query = `
			INSERT INTO
					t_users_transaction
			VALUES
					($1, $2, $3, $4, $5, $6)
			RETURNING transaction_id`

	err = tx.QueryRowx(query,
		t.UserUid,
		t.TransactionID,
		t.Amount,
		t.FromUserUid,
		t.ToUserUid,
		t.TransactionDate,
	).Scan(&transactionID)
	if err != nil {
		return "", err
	}

	return
}

func (repo *Repo) getUserFromId(UserUid string) (userUid int64, err error) {
	const userQuery = `
				SELECT balance FROM t_user
				WHERE user_uid = $1
	`
	var user entity.User
	if err := repo.db.Get(&user, userQuery, UserUid); err != nil {
		if err == sql.ErrNoRows {
			return 0, err
		}
		return 0, err
	}

	return user.Balance, nil
}
