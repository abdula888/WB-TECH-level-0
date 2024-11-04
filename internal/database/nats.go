package database

import (
	"encoding/json"
	"log"

	"WB-TECH-level-0/internal/cache"
	"WB-TECH-level-0/internal/models"

	"github.com/nats-io/stan.go"
)

var sc stan.Conn

func InitNATS() {
	var err error
	sc, err = stan.Connect("test-cluster", "client-123")
	if err != nil {
		log.Fatal(err)
	}
	_, err = sc.Subscribe("order_updates", handleOrderMessage)
	if err != nil {
		log.Fatal(err)
	}
}

func handleOrderMessage(msg *stan.Msg) {
	var order models.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return
	}
	// Запись в базу данных
	query := `
    INSERT INTO orders (
        order_uid, track_number, entry, locale, internal_signature, customer_id, 
        delivery_service, shardkey, sm_id, date_created, oof_shard
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`

	_, err = db.Exec(query,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.ShardKey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		log.Printf("Error inserting data: %v", err)
	}
	// Кэширование данных
	cache.SetCache(order)
}
