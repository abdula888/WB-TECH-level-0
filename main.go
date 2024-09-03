package main

import (
	"log"
	"net/http"
)

func main() {
	// Инициализация базы данных
	initDB()
	// Инициализация подключения к NATS Streaming
	initNATS()
	// Загрузка данных из базы данных в кэш
	loadCacheFromDB()

	// Запуск HTTP-сервера
	log.Println("Serving static files from ./static on http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Обработка запросов на поиск заказа по ID
	http.HandleFunc("/order/", orderHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
