package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	// Подключаемся к NATS Streaming
	sc, err := stan.Connect("test-cluster", "publisher-123")
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Читаем JSON файл
	data, err := os.ReadFile("order.json")
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Отправляем данные в канал order_updates
	err = sc.Publish("order_updates", data)
	if err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	fmt.Println("Order published successfully!")
}
