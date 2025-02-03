package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"wallet-service/service"

	"github.com/gorilla/mux"
)

// Структура POST-запроса
type WalletRequest struct {
	WalletID      string  `json:"walletId"`
	OperationType string  `json:"operationType"` // "DEPOSIT" или "WITHDRAW"
	Amount        float64 `json:"amount"`
}

// Обработка POST /api/v1/wallet
func WalletOperationHandler(w http.ResponseWriter, r *http.Request) {

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "cannot read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Разбор JSON-запроса
	var req WalletRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	//Обработка операции
	if err := service.ProcessWalletOperation(req.WalletID, req.OperationType, req.Amount); err != nil {
		log.Printf("Error processing operation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Operation successful"))
}

// Обработка GET /api/v1/wallets/{walletId}
func GetWalletHandler(w http.ResponseWriter, r *http.Request) {

	// Извлечение walletId из URL.
	vars := mux.Vars(r)
	walletId := vars["walletId"]

	// Получение баланса через бизнес-логику.
	balance, err := service.GetWalletBalance(walletId)
	if err != nil {
		log.Printf("Error getting wallet balance: %v", err)
		http.Error(w, "Error getting balance", http.StatusInternalServerError)
		return
	}

	// Формирование и отправка JSON-ответа.
	resp := map[string]interface{}{
		"walletId": walletId,
		"balance":  balance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
