package main

import (
	"WB-TECH-level-0/internal/database"
	"WB-TECH-level-0/internal/server"
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	database.InitDB()
	// Инициализация подключения к NATS Streaming
	database.InitNATS()
	// Загрузка данных из базы данных в кэш
	database.LoadCacheFromDB()

	// Запуск HTTP-сервера
	log.Println("Serving static files from ./static on http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Обработка запросов на поиск заказа по ID
	http.HandleFunc("/order/", server.OrderHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
