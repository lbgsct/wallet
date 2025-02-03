package repository

import (
	"database/sql"
	"errors"
)

// Транзакция для изменения баланса кошелька
func UpdateWalletBalance(walletID string, amount float64, operationType string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Чтение текущего баланса (с блокировкой на запись)
	var currentBalance float64
	err = tx.QueryRow(`SELECT balance FROM wallets WHERE wallet_id = $1 FOR UPDATE`, walletID).Scan(&currentBalance)
	if err == sql.ErrNoRows {
		// Если кошелёк не найден, создаём с нулевым балансом
		_, err = tx.Exec(`INSERT INTO wallets (wallet_id, balance) VALUES ($1, 0)`, walletID)
		if err != nil {
			return err
		}
		currentBalance = 0
	} else if err != nil {
		return err
	}

	// Новый баланс
	var newBalance float64
	switch operationType {
	case "DEPOSIT":
		newBalance = currentBalance + amount
	case "WITHDRAW":
		newBalance = currentBalance - amount
		if newBalance < 0 {
			return errors.New("not enough funds")
		}
	default:
		return errors.New("unknown operation type")
	}

	// Обновление баланса
	_, err = tx.Exec(`UPDATE wallets SET balance = $1 WHERE wallet_id = $2`, newBalance, walletID)
	if err != nil {
		return err
	}

	// Подтверждение транзакции
	return tx.Commit()
}

// Текущий баланс кошелька
func GetWalletBalance(walletID string) (float64, error) {
	var balance float64
	err := DB.QueryRow(`SELECT balance FROM wallets WHERE wallet_id = $1`, walletID).Scan(&balance)
	if err == sql.ErrNoRows {
		// Если кошелёк не найден, пусть вернётся 0
		return 0, nil
	}
	return balance, err
}
