package user

import (
	"encoding/json"
	"net/http"

	"finstat/internal/entity"
)

type Transport struct {
}

func NewTransport() *Transport {
	return &Transport{}
}

func (t *Transport) CreateUserDecode(r *http.Request) (user entity.User, err error) {

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return user, err
	}
	defer r.Body.Close()

	if !user.IsNice() {
		return user, entity.ErrInvalidDataLength
	}

	user.GenerateUid()

	return
}

func (t *Transport) CreateUserEncode(w http.ResponseWriter, user entity.User) error {
	userdata, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(userdata)
	return nil
}

func (t *Transport) UserEncode(w http.ResponseWriter, user *entity.User) error {
	userdata, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(userdata)
	return nil
}

func (t *Transport) TransactionEncode(w http.ResponseWriter, transactions []entity.Transaction) error {
	tr, _ := json.Marshal(transactions)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(tr)
	return nil
}
