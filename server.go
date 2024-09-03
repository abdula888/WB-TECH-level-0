package main

import (
	"encoding/json"
	"net/http"
)

func orderHandler(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("id")
	order, found := getCache(orderUID)
	if !found {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
