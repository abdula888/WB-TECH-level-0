package server

import (
	"WB-TECH-level-0/internal/cache"
	"encoding/json"
	"net/http"
)

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("id")
	order, found := cache.GetCache(orderUID)
	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
