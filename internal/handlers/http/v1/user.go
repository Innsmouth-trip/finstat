package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"finstat/internal/entity"
	"finstat/internal/handlers/http/v1/transport/errorworker"
	"finstat/internal/handlers/http/v1/transport/user"
)

type UserSvc interface {
	CreateUser(ctx context.Context, u entity.User) (UserUid string, err error)
	GetUser(userUid string) (user *entity.User, err error)
	GetUserTransaction(userUid string) ([]entity.Transaction, error)
}

type UserHandler struct {
	userSvc       UserSvc
	userTransport UserTransport
	errorWorker   ErrorWorker
}

type UserTransport interface {
	CreateUserDecode(r *http.Request) (entity.User, error)
	CreateUserEncode(w http.ResponseWriter, user entity.User) error
	UserEncode(w http.ResponseWriter, user *entity.User) error
	TransactionEncode(w http.ResponseWriter, transactions []entity.Transaction) error
}

type ErrorWorker interface {
	ProcessingError(w http.ResponseWriter, err error)
}

func NewUserHandler(svc UserSvc) *UserHandler {
	return &UserHandler{
		userSvc:       svc,
		userTransport: user.NewTransport(),
		errorWorker:   errorworker.NewError(),
	}
}

/**
 * @api {post} /user/ CreateUser
 * @apiName CreateUser
 * @apiGroup User
 *
 * @apiBody {String} firstname firstname.
 * @apiBody {String} lastname  lastname.
 * @apiBody {String} username  username.
 *
 * @apiSuccess (200) {Object} entity.Transaction full transaction object.
 * @apiError  (400) {String} ErrInvalidField invalid field length.
 */
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.userTransport.CreateUserDecode(r)
	if err != nil {
		h.errorWorker.ProcessingError(w, err)
		return
	}

	_, err = h.userSvc.CreateUser(r.Context(), user)
	if err != nil {
		h.errorWorker.ProcessingError(w, err)
		return
	}

	err = h.userTransport.CreateUserEncode(w, user)
	if err != nil {
		h.errorWorker.ProcessingError(w, err)
		return
	}
	return
}

/**
 * @api {get} /user/{id} GetUserFromUid
 * @apiName GetUserFromUid
 * @apiGroup User
 *
 * @apiParam {String} user_uid unique user uid.
 *
 * @apiSuccess (200) {Object} entity.User full user object.
 * @apiError (404) {String} ErrUserNotFound user not found.
 */
func (h *UserHandler) GetUserFromID(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["id"]
	usr, err := h.userSvc.GetUser(uid)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "user not found")
		return
	}

	err = h.userTransport.UserEncode(w, usr)
	if err != nil {
		h.errorWorker.ProcessingError(w, err)
		return
	}
	return
}

/**
 * @api {get} /user/transactions/{user_uid} GetUserTransactions
 * @apiName GetUserTransactions
 * @apiGroup User
 *
 * @apiParam {String} user_uid user uid.
 *
 * @apiSuccess (200) {Object} entity.User full user object.
 * @apiError  (404) {String}  ErrTransactionNotFound transaction not found.
 */
func (h *UserHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	uid := mux.Vars(r)["id"]
	transactions, err := h.userSvc.GetUserTransaction(uid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.userTransport.TransactionEncode(w, transactions)
	if err != nil {
		h.errorWorker.ProcessingError(w, err)
		return
	}
	return
}
