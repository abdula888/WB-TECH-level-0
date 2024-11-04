package database

import (
	"database/sql"
	"log"

	"WB-TECH-level-0/internal/cache"
	"WB-TECH-level-0/internal/models"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	connStr := "user=order_user password=123 dbname=order_service sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

func LoadCacheFromDB() {
	rows, err := db.Query(`SELECT 
        o.order_uid,
        o.track_number,
        o.entry,
        o.locale,
        o.internal_signature,
        o.customer_id,
        o.delivery_service,
        o.shardkey,
        o.sm_id,
        o.date_created,
        o.oof_shard,
        d.name AS delivery_name,
        d.phone AS delivery_phone,
        d.zip AS delivery_zip,
        d.city AS delivery_city,
        d.address AS delivery_address,
        d.region AS delivery_region,
        d.email AS delivery_email,
        p.transaction AS payment_transaction,
        p.currency AS payment_currency,
        p.provider AS payment_provider,
        p.amount AS payment_amount,
        p.payment_dt AS payment_date,
        p.bank AS payment_bank,
        p.delivery_cost AS payment_delivery_cost,
        p.goods_total AS payment_goods_total,
        p.custom_fee AS payment_custom_fee,
        i.chrt_id AS item_chrt_id,
        i.track_number AS item_track_number,
        i.price AS item_price,
        i.rid AS item_rid,
        i.name AS item_name,
        i.sale AS item_sale,
        i.size AS item_size,
        i.total_price AS item_total_price,
        i.nm_id AS item_nm_id,
        i.brand AS item_brand,
        i.status AS item_status
    FROM orders o
    LEFT JOIN delivery d ON o.order_uid = d.order_uid
    LEFT JOIN payment p ON o.order_uid = p.order_uid
    LEFT JOIN items i ON o.order_uid = i.order_uid`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Создаём карту для хранения заказов по order_uid
	orderMap := make(map[string]*models.Order)

	for rows.Next() {
		var order models.Order
		var delivery models.Delivery
		var payment models.Payment
		var item models.Item

		err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.ShardKey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&delivery.Name,
			&delivery.Phone,
			&delivery.Zip,
			&delivery.City,
			&delivery.Address,
			&delivery.Region,
			&delivery.Email,
			&payment.Transaction,
			&payment.Currency,
			&payment.Provider,
			&payment.Amount,
			&payment.PaymentDT,
			&payment.Bank,
			&payment.DeliveryCost,
			&payment.GoodsTotal,
			&payment.CustomFee,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			log.Fatal(err)
		}

		// Если заказ уже существует в карте, добавляем товар
		if existingOrder, exists := orderMap[order.OrderUID]; exists {
			// Добавление товара в существующий заказ
			if item.ChrtID != 0 {
				existingOrder.Items = append(existingOrder.Items, item)
			}
		} else {
			// Создание нового заказа с пустым срезом товаров
			newOrder := models.Order{
				OrderUID:          order.OrderUID,
				TrackNumber:       order.TrackNumber,
				Entry:             order.Entry,
				Locale:            order.Locale,
				InternalSignature: order.InternalSignature,
				CustomerID:        order.CustomerID,
				DeliveryService:   order.DeliveryService,
				ShardKey:          order.ShardKey,
				SmID:              order.SmID,
				DateCreated:       order.DateCreated,
				OofShard:          order.OofShard,
				Delivery:          delivery,
				Payment:           payment,
				Items:             []models.Item{}, // Инициализация пустого среза
			}
			// Добавление нового заказа в карту
			orderMap[order.OrderUID] = &newOrder

			// Добавление товара, если он существует
			if item.ChrtID != 0 {
				orderMap[order.OrderUID].Items = append(orderMap[order.OrderUID].Items, item)
			}
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Преобразование карты в срез и запись в кэш
	for _, order := range orderMap {
		cache.SetCache(*order)
	}
}
