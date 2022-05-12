package v1

import (
	"net/http"

	"finstat/internal/entity"
	"finstat/internal/handlers/http/v1/transport/balance"
	"finstat/internal/handlers/http/v1/transport/errorworker"
)

type BalanceSvc interface {
	AddMoneyToUser(userUid string, amount int64) error
	SendMoneyFromUser(fromUserUid, toUserUid string, amount int64) error
}

type BalanceHandler struct {
	BalanceSvc       BalanceSvc
	BalanceTransport BalanceTransport
	ErrorWorker      ErrorWorker
}

type BalanceTransport interface {
	CheckAddMoney(r *http.Request, c entity.Currency) (balance entity.AddMoneyToUserOpts, err error)
	CheckSendMoney(r *http.Request, c entity.Currency) (balance entity.SendMoneyFromUserOpts, err error)
	BalanceEncode(w http.ResponseWriter, str string) error
}

func NewBalanceHandler(svc BalanceSvc) *BalanceHandler {
	return &BalanceHandler{
		BalanceSvc:       svc,
		BalanceTransport: balance.NewTransport(),
		ErrorWorker:      errorworker.NewError(),
	}
}

/**
 * @api {put} /balance/add Add money to user
 * @apiName AddMoneyToUser
 * @apiGroup Balance
 *
 * @apiBody {String} user_uid user uid.
 * @apiBody {String} amount amount (format = "50.00").
 *
 * @apiSuccess (200) {String} OK adding money to user is done correctly.
 * @apiError  (400) {String} ErrInvalidField invalid field.
 */
func (h *BalanceHandler) AddMoneyToUser(w http.ResponseWriter, r *http.Request) {
	var c entity.Currency

	balance, err := h.BalanceTransport.CheckAddMoney(r, c)
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}

	err = h.BalanceSvc.AddMoneyToUser(balance.UserUid, int64(balance.AmountInt))
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}

	err = h.BalanceTransport.BalanceEncode(w, "add money to user is done correctly")
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}
}

/**
 * @api {put} /balance/send Send money: user to user
 * @apiName SendMoneyFromUser
 * @apiGroup Balance
 *
 * @apiBody {String} from_user_uid from user uid.
 * @apiBody {String} to_user_uid to user uid.
 * @apiBody {String} amount amount (format = "50.00").
 *
 * @apiSuccess (200) {String} OK sending money to user is done correctly.
 * @apiError  (400) {String} ErrInvalidField invalid field length.
 */
func (h *BalanceHandler) SendMoneyFromUser(w http.ResponseWriter, r *http.Request) {
	var c entity.Currency
	balance, err := h.BalanceTransport.CheckSendMoney(r, c)
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}

	err = h.BalanceSvc.SendMoneyFromUser(balance.FromUserUid, balance.ToUserUid, int64(balance.AmountInt))
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}

	err = h.BalanceTransport.BalanceEncode(w, "send money to user is done correctly")
	if err != nil {
		h.ErrorWorker.ProcessingError(w, err)
		return
	}
}
