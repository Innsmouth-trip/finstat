package entity

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	minLength = 2
	maxLength = 32
)

var (
	ErrInvalidDataLength    = errors.New("invalid username length")
	ErrWrongAmount          = errors.New("wrong transfer amount")
	ErrInvalidUserUidLength = errors.New("invalid userUid length")
	ErrInvalidSet           = errors.New("invalid amount set")
)

type User struct {
	UserUid   string `json:"user_uid" db:"user_uid"`
	Firstname string `json:"firstname" db:"firsname"`
	Lastname  string `json:"lastname" db:"lastname"`
	Username  string `json:"username" db:"username"`
	Balance   int64  `json:"balance" db:"balance"`
}

type Transaction struct {
	UserUid         string    `json:"user_uid" db:"user_uid"`
	TransactionID   string    `json:"transaction_id" db:"transaction_id"`
	Amount          int64     `json:"amount" db:"amount"`
	FromUserUid     string    `json:"from_user_uid" db:"from_user"`
	ToUserUid       string    `json:"to_user" db:"to_user"`
	TransactionDate time.Time `json:"transaction_date" db:"transaction_date"`
}

type AddMoneyToUserOpts struct {
	UserUid string `json:"user_uid"`
	Amount  string `json:"amount"`

	AmountInt Currency
}

type SendMoneyFromUserOpts struct {
	FromUserUid string `json:"from_user_uid"`
	ToUserUid   string `json:"to_user_uid"`
	Amount      string `json:"amount"`

	AmountInt Currency
}

type Currency int64

func (u *User) IsNice() bool {
	if !niceName(u.Username) {
		return false
	}

	if !niceName(u.Lastname) {
		return false
	}

	if !niceName(u.Firstname) {
		return false
	}

	return true
}

func niceName(userdata string) bool {
	length := len(userdata)
	return minLength <= length && length <= maxLength
}

func (u *User) GenerateUid() {
	u.UserUid = uuid.NewString()
	return
}

func (a *AddMoneyToUserOpts) AddmoneyOptsCheck() bool {
	if len(a.UserUid) == 0 {
		return true
	}

	if a.AmountInt <= 0 {
		return true
	}

	return false
}

func (c *Currency) Set(s string) (Currency, error) {
	srubles, spenny, ok := strings.Cut(s, ".")
	if !ok {
		return 0, ErrInvalidSet
	}

	rubles, err := strconv.ParseInt(srubles, 10, 64)
	if err != nil {
		return 0, err
	}

	penny, err := strconv.ParseInt(spenny, 10, 64)
	if err != nil {
		return 0, err
	}

	p := Currency(rubles*100 + penny)

	return p, nil
}

func (s *SendMoneyFromUserOpts) SendMonetOptsCheck() bool {
	if len(s.FromUserUid) == 0 {
		return true
	}

	if len(s.ToUserUid) == 0 {
		return true
	}

	if s.AmountInt <= 0 {
		return true
	}

	return false
}
