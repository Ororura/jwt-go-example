package handler

import (
	"encoding/json"
	"net/http"

	"jwt-go/internal/usecase"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(u *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.usecase.GetProducts()
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
