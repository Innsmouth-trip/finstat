package balance

import (
	"encoding/json"
	"io"
	"net/http"

	"finstat/internal/entity"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (t *Transport) CheckAddMoney(r *http.Request, c entity.Currency) (balance entity.AddMoneyToUserOpts, err error) {
	err = json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		return balance, err
	}
	defer r.Body.Close()

	balance.AmountInt, err = c.Set(balance.Amount)
	if err != nil {
		return balance, err
	}

	if balance.AddmoneyOptsCheck() {
		return balance, entity.ErrWrongAmount
	}

	return
}

func (t *Transport) CheckSendMoney(r *http.Request, c entity.Currency) (balance entity.SendMoneyFromUserOpts, err error) {
	err = json.NewDecoder(r.Body).Decode(&balance)
	if err != nil {
		return balance, err
	}

	balance.AmountInt, err = c.Set(balance.Amount)
	if err != nil {
		return balance, err
	}

	if balance.SendMonetOptsCheck() {
		return balance, entity.ErrWrongAmount
	}

	defer r.Body.Close()

	return
}

func (t *Transport) BalanceEncode(w http.ResponseWriter, str string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, str)
	return nil
}
