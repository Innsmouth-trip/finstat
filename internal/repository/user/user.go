package user

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"finstat/internal/entity"
)

type Repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo *Repo) CreateUser(ctx context.Context, u entity.User) (userUid string, err error) {
	tx, err := repo.db.Beginx()
	if err != nil {
		return
	}

	userUid, err = repo.createUser(ctx, tx, u)
	if err != nil {
		_ = tx.Rollback()
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	return
}

func (repo *Repo) createUser(ctx context.Context, q sqlx.QueryerContext, u entity.User) (userUid string, err error) {
	const userQuery = `
			INSERT INTO
					t_user
			VALUES
					($1, $2, $3, $4)
			RETURNING user_uid`

	err = q.QueryRowxContext(ctx, userQuery, u.UserUid, u.Firstname, u.Lastname, u.Username).Scan(&userUid)
	if err != nil {
		return userUid, err
	}

	return
}

func (repo *Repo) GetUserFromId(UserUid string) (*entity.User, error) {
	const userQuery = `
				SELECT * FROM t_user
				WHERE user_uid = $1
	`
	var user entity.User
	if err := repo.db.Get(&user, userQuery, UserUid); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &user, nil
}

func (repo *Repo) GetUsersTransactions(UserUid string) ([]entity.Transaction, error) {
	const transactionQuery = `
				SELECT * FROM t_users_transaction
				WHERE user_uid = $1
	`
	var trs []entity.Transaction
	if err := repo.db.Select(&trs, transactionQuery, UserUid); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return trs, nil
}
