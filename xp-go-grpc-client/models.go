package main

import (
	"time"
)

type Account struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	DateOfBirth   time.Time `json:"date_of_birth"`
	AccountNumber string    `json:"account_number"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
}
