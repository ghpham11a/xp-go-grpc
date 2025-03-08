package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	pb "xp-go-grpc-client/proto"

	datepb "google.golang.org/genproto/googleapis/type/date"
)

func (app *App) CreateAccount(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse JSON into a struct
	var account Account
	if err := json.Unmarshal(reqBody, &account); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the struct to a protobuf message
	conertedId, conertedIdErr := strconv.Atoi(account.ID)
	if conertedIdErr != nil {
		log.Fatalf("Failed to convert string to int32: %v", err)
	}

	// Convert the date of birth to a protobuf date
	dob := &datepb.Date{
		Year:  int32(account.DateOfBirth.Year()),
		Month: int32(account.DateOfBirth.Month()),
		Day:   int32(account.DateOfBirth.Day()),
	}

	pbAccount := &pb.Account{
		Id:            int32(conertedId),
		Email:         account.Email,
		DateOfBirth:   dob,
		AccountNumber: account.AccountNumber,
		Balance:       fmt.Sprintf("%f", account.Balance),
		CreatedAt:     nil,
	}

	acct, acctErr := app.AccountsClient.CreateAccountRPC(r.Context(), pbAccount)
	if acctErr != nil {
		log.Fatalf("could not greet: %v", acctErr)
	}
	fmt.Printf("Account Number: %s\n", acct.GetAccountNumber())

}

func (app *App) ReadAccounts(w http.ResponseWriter, r *http.Request) {

}

func (app *App) ReadAccount(w http.ResponseWriter, r *http.Request) {

}

func (app *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {

}

func (app *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {

}
