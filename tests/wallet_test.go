// tests/wallet_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"wallet-service/handlers"
	"wallet-service/repository"
	"wallet-service/service"
)

// Выполняется перед запуском всех тестов и инициализирует подключение к БД.
func TestMain(m *testing.M) {

	err := godotenv.Load("../config.env")
	if err != nil {
		println("Warning: no config.env file found, using system environment variables")
	}

	os.Setenv("DB_HOST", "localhost")

	err = repository.InitDB()
	if err != nil {
		panic("failed to initialize DB: " + err.Error())
	}

	time.Sleep(2 * time.Second)

	code := m.Run()

	repository.DB.Close()

	os.Exit(code)
}

// Проверяет валидацию суммы в бизнес-логике.
func TestProcessWalletOperationInvalidAmount(t *testing.T) {
	err := service.ProcessWalletOperation("test-wallet", "DEPOSIT", -10)
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

// Проверяет, что неверный тип операции приводит к ошибке.
func TestProcessWalletOperationInvalidOperation(t *testing.T) {
	err := service.ProcessWalletOperation("test-wallet", "INVALID", 100)
	if err == nil {
		t.Error("expected error for invalid operation type, got nil")
	}
}

// Проверяет обработчик POST /api/v1/wallet.
func TestWalletOperationHandler(t *testing.T) {
	// Формируем тело запроса в формате JSON
	reqBody := map[string]interface{}{
		"walletId":      "test-wallet",
		"operationType": "DEPOSIT",
		"amount":        50,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal json: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v1/wallet", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("failed to create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.WalletOperationHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}
}

// Проверяет обработчик GET /api/v1/wallets/{walletId}.
func TestGetWalletHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/wallets/test-wallet", nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/wallets/{walletId}", handlers.GetWalletHandler)
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response: %v", err)
	}
	if response["walletId"] != "test-wallet" {
		t.Errorf("expected walletId 'test-wallet', got %v", response["walletId"])
	}
	if _, ok := response["balance"].(float64); !ok {
		t.Errorf("expected balance to be a float64, got %T", response["balance"])
	}
}
