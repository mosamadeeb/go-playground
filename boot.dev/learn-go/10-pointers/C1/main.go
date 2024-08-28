package main

import (
	"errors"
)

type customer struct {
	id      int
	balance float64
}

type transactionType string

const (
	transactionDeposit    transactionType = "deposit"
	transactionWithdrawal transactionType = "withdrawal"
)

type transaction struct {
	customerID      int
	amount          float64
	transactionType transactionType
}

// Don't touch above this line

func updateBalance(cust *customer, transact transaction) error {
	// It was not requested to check the customer ID in the transaction

	switch transact.transactionType {
	case transactionDeposit:
		cust.balance += transact.amount
	case transactionWithdrawal:
		if cust.balance < transact.amount {
			return errors.New("insufficient funds")
		}

		cust.balance -= transact.amount
	default:
		return errors.New("unknown transaction type")
	}

	return nil
}
