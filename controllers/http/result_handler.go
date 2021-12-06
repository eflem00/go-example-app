package http

import (
	"fmt"
	"net/http"

	"github.com/eflem00/go-example-app/usecases"
	"github.com/go-chi/chi/v5"
)

type ResultHandler struct {
	resultUsecase *usecases.ResultUsecase
}

func NewResultHandler(resultUsecase *usecases.ResultUsecase) *ResultHandler {
	return &ResultHandler{
		resultUsecase,
	}
}

func (handler *ResultHandler) GetResultById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	val, err := handler.resultUsecase.GetResultById(r.Context(), id)

	if err != nil {
		http.Error(w, "Error reading key", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, val)
}
