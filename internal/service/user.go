package service

import (
	"context"

	"finstat/internal/entity"
)

func (svc *Service) CreateUser(ctx context.Context, u entity.User) (UserUid string, err error) {
	userUid, err := svc.UserRepo.CreateUser(ctx, u)
	if err != nil {
		return "create user error", err
	}

	return userUid, nil
}

func (svc *Service) GetUser(userUid string) (user *entity.User, err error) {
	u, err := svc.UserRepo.GetUserFromId(userUid)
	if err != nil {
		return nil, err
	}

	return u, nil

}

func (svc *Service) GetUserTransaction(userUid string) ([]entity.Transaction, error) {
	tr, err := svc.UserRepo.GetUsersTransactions(userUid)
	if err != nil {
		return nil, err
	}

	return tr, nil
}
