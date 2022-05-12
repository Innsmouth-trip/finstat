package service

func (svc *Service) AddMoneyToUser(userUid string, amount int64) error {
	err := svc.BalanceRepo.AddMoneyToUser(userUid, amount)
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) SendMoneyFromUser(fromUserUid, toUserUid string, amount int64) error {
	err := svc.BalanceRepo.SendMoneyFromUser(fromUserUid, toUserUid, amount)
	if err != nil {
		return err
	}
	return nil
}
