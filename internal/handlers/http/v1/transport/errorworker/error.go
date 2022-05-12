package errorworker

import (
	"encoding/json"
	"net/http"
)

type ErrorWorker struct {
}

func NewError() *ErrorWorker {
	return &ErrorWorker{}
}

type Errors struct {
	Error string `json:"error"`
}

func (e *ErrorWorker) ProcessingError(w http.ResponseWriter, err error) {
	errForResponse := Errors{
		Error: err.Error(),
	}

	errResp, _ := json.Marshal(errForResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(errResp)
}
