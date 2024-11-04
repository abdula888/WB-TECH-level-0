package cache

import (
	"WB-TECH-level-0/internal/models"
	"sync"
)

var cache = make(map[string]models.Order)
var mu sync.RWMutex

func SetCache(order models.Order) {
	mu.Lock()
	defer mu.Unlock()
	cache[order.OrderUID] = order
}

func GetCache(orderUID string) (models.Order, bool) {
	mu.RLock()
	defer mu.RUnlock()
	order, found := cache[orderUID]
	return order, found
}
