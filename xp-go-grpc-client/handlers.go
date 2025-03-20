package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "xp-go-grpc-client/proto"

	"github.com/gorilla/mux"
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

	// Parse dateOfBirth from "YYYY-MM-DD" string
	var dob *datepb.Date
	if account.DateOfBirth != "" {
		parsedDOB, err := time.Parse("2006-01-02", account.DateOfBirth)
		if err != nil {
			http.Error(w, "Invalid DateOfBirth format (expected YYYY-MM-DD): "+err.Error(), http.StatusBadRequest)
			return
		}
		dob = &datepb.Date{
			Year:  int32(parsedDOB.Year()),
			Month: int32(parsedDOB.Month()),
			Day:   int32(parsedDOB.Day()),
		}
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

// ReadAccount fetches an account from the gRPC server
func (app *App) ReadAccount(w http.ResponseWriter, r *http.Request) {

	fmt.Println("ReadAccount called")

	vars := mux.Vars(r)
	id := vars["accountId"]

	if id == "" {
		http.Error(w, "Missing account number", http.StatusBadRequest)
		return
	}
	// Call gRPC GetAccountRPC
	convertedID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID format: "+err.Error(), http.StatusBadRequest)
		return
	}
	account, err := app.AccountsClient.ReadAccountRPC(r.Context(), &pb.ReadAccountRequest{Id: int32(convertedID)})
	if err != nil {
		http.Error(w, "Failed to fetch account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convert response to JSON
	response, _ := json.Marshal(account)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// UpdateAccount updates an account using the gRPC server
func (app *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	var account Account
	if err := json.Unmarshal(reqBody, &account); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Convert ID to int32
	convertedID, err := strconv.Atoi(account.ID)
	if err != nil {
		http.Error(w, "Invalid ID format: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse dateOfBirth from "YYYY-MM-DD" string
	var dob *datepb.Date
	if account.DateOfBirth != "" {
		parsedDOB, err := time.Parse("2006-01-02", account.DateOfBirth)
		if err != nil {
			http.Error(w, "Invalid DateOfBirth format (expected YYYY-MM-DD): "+err.Error(), http.StatusBadRequest)
			return
		}
		dob = &datepb.Date{
			Year:  int32(parsedDOB.Year()),
			Month: int32(parsedDOB.Month()),
			Day:   int32(parsedDOB.Day()),
		}
	}

	pbAccount := &pb.Account{
		Id:            int32(convertedID),
		Email:         account.Email,
		DateOfBirth:   dob,
		AccountNumber: account.AccountNumber,
		Balance:       fmt.Sprintf("%f", account.Balance),
		CreatedAt:     nil,
	}

	// Call gRPC server
	updatedAccount, err := app.AccountsClient.UpdateAccountRPC(r.Context(), pbAccount)
	if err != nil {
		http.Error(w, "Failed to update account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	response, _ := json.Marshal(updatedAccount)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// DeleteAccount deletes an account using the gRPC server
func (app *App) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	accountNumber := r.URL.Query().Get("account_number")
	if accountNumber == "" {
		http.Error(w, "Missing account number", http.StatusBadRequest)
		return
	}

	// Call gRPC DeleteAccountRPC
	convertedID, err := strconv.Atoi(accountNumber)
	if err != nil {
		http.Error(w, "Invalid account number format: "+err.Error(), http.StatusBadRequest)
		return
	}
	_, err = app.AccountsClient.DeleteAccountRPC(r.Context(), &pb.DeleteAccountRequest{Id: int32(convertedID)})
	if err != nil {
		http.Error(w, "Failed to delete account: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Account deleted successfully"}`))
}
