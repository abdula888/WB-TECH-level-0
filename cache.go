package main

import "sync"

var cache = make(map[string]Order)
var mu sync.RWMutex

func setCache(order Order) {
	mu.Lock()
	defer mu.Unlock()
	cache[order.OrderUID] = order
}

func getCache(orderUID string) (Order, bool) {
	mu.RLock()
	defer mu.RUnlock()
	order, found := cache[orderUID]
	return order, found
}
