package service

import (
	"errors"
	"wallet-service/repository"
)

// Проверка входных данных
func ProcessWalletOperation(walletID, operationType string, amount float64) error {
	if walletID == "" {
		return errors.New("walletID is required")
	}
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	return repository.UpdateWalletBalance(walletID, amount, operationType)
}

// Возврат баланс
func GetWalletBalance(walletID string) (float64, error) {
	if walletID == "" {
		return 0, errors.New("walletID is required")
	}
	return repository.GetWalletBalance(walletID)
}
