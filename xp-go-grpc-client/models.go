package main

import (
	"time"
)

type Account struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	DateOfBirth   string    `json:"dateOfBirth"`
	AccountNumber string    `json:"accountNumber"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"createdAt"`
}
