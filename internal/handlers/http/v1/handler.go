package v1

import (
	"net/http"

	"github.com/gorilla/mux"

	"finstat/internal/service"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2

type Handlers struct {
	UserHandler    *UserHandler
	BalanceHandler *BalanceHandler
}

func NewHandler(service *service.Service) *Handlers {
	return &Handlers{
		UserHandler:    NewUserHandler(service),
		BalanceHandler: NewBalanceHandler(service),
	}
}

func (h *Handlers) Init() *mux.Router {
	route := mux.NewRouter()
	r := route.PathPrefix("/v1").Subrouter()
	h.init(r)

	return r
}

func (h *Handlers) init(r *mux.Router) *mux.Router {
	h.initUserRoutes(r)
	h.initBalanceRoutes(r)

	return r
}

func (h *Handlers) initUserRoutes(r *mux.Router) {
	r.HandleFunc("/user/", h.UserHandler.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/user/{id}", h.UserHandler.GetUserFromID).Methods(http.MethodGet)
	r.HandleFunc("/user/transactions/{id}", h.UserHandler.GetUserTransactions).Methods(http.MethodGet)
}

func (h *Handlers) initBalanceRoutes(r *mux.Router) {
	r.HandleFunc("/balance/add", h.BalanceHandler.AddMoneyToUser).Methods(http.MethodPut)
	r.HandleFunc("/balance/send", h.BalanceHandler.SendMoneyFromUser).Methods(http.MethodPut)
}
