package service

import (
	"context"

	"finstat/internal/entity"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u entity.User) (userUid string, err error)
	GetUserFromId(UserUid string) (*entity.User, error)
	GetUsersTransactions(UserUid string) ([]entity.Transaction, error)
}

type BalanceRepo interface {
	AddMoneyToUser(userUid string, amount int64) error
	SendMoneyFromUser(fromUserUid, toUserUid string, amount int64) error
}

type Service struct {
	UserRepo    UserRepo
	BalanceRepo BalanceRepo
}

func NewService(userRepo UserRepo, balanceRepo BalanceRepo) *Service {
	return &Service{
		UserRepo:    userRepo,
		BalanceRepo: balanceRepo,
	}
}
