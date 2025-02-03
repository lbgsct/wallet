package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"wallet-service/handlers"
	"wallet-service/repository"
)

func main() {

	// Загрузка переменных окружения из config.env
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("Warning: No config.env file found, using system env")
	}

	// Подключение к базе данных
	if err := repository.InitDB(); err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}
	defer repository.DB.Close()

	// Роутер
	r := mux.NewRouter()

	// Настройка маршрутов
	r.HandleFunc("/api/v1/wallet", handlers.WalletOperationHandler).Methods("POST")
	r.HandleFunc("/api/v1/wallets/{walletId}", handlers.GetWalletHandler).Methods("GET")

	// Порт из переменных окружения
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
