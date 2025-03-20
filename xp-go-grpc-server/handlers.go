package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	// Import the generated protobuf package
	pb "xp-go-grpc-server/proto"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/pgtype"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (app *App) CreateAccountRPC(ctx context.Context, in *pb.Account) (*pb.Account, error) {

	fmt.Println("CreateAccountRPC")

	// Convert date_of_birth (google.type.Date) to PostgreSQL DATE format
	dob := in.GetDateOfBirth()
	var dobStr *string
	if dob != nil {
		formattedDOB := fmt.Sprintf("%04d-%02d-%02d", dob.GetYear(), dob.GetMonth(), dob.GetDay())
		dobStr = &formattedDOB
	}

	// Convert created_at (google.protobuf.Timestamp) to time.Time
	createdAt := in.GetCreatedAt()
	var createdAtTime time.Time
	if createdAt != nil {
		createdAtTime = createdAt.AsTime()
	} else {
		createdAtTime = time.Now()
	}

	// Insert into the database
	query := `INSERT INTO accounts (email, date_of_birth, account_number, balance, created_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int32
	err := app.DB.QueryRow(ctx, query, in.GetEmail(), dobStr, in.GetAccountNumber(), in.GetBalance(), createdAtTime).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert account: %w", err)
	}

	// Return the created account with the new ID
	return &pb.Account{
		Id:            id,
		Email:         in.GetEmail(),
		DateOfBirth:   in.GetDateOfBirth(),
		AccountNumber: in.GetAccountNumber(),
		Balance:       in.GetBalance(),
		CreatedAt:     in.GetCreatedAt(),
	}, nil
}

func (app *App) ReadAccountRPC(ctx context.Context, in *pb.ReadAccountRequest) (*pb.Account, error) {
	// Query for the account
	query := `SELECT id, email, date_of_birth, account_number, balance, created_at FROM accounts WHERE id = $1`
	var id int32
	var email, accountNumber, balance string
	var dob pgtype.Date
	var createdAt time.Time

	// accountId := strconv.Itoa(int(in.GetId()))

	err := app.DB.QueryRow(ctx, query, in.GetId()).Scan(&id, &email, &dob, &accountNumber, &balance, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	// Convert PostgreSQL DATE to google.type.Date
	var pbDob *date.Date
	if dob.Status == pgtype.Present {
		pbDob = &date.Date{
			Year:  int32(dob.Time.Year()),
			Month: int32(dob.Time.Month()),
			Day:   int32(dob.Time.Day()),
		}
	}

	// Convert created_at timestamp to google.protobuf.Timestamp
	pbCreatedAt := timestamppb.New(createdAt)

	// Return the account
	return &pb.Account{
		Id:            id,
		Email:         email,
		DateOfBirth:   pbDob,
		AccountNumber: accountNumber,
		Balance:       balance,
		CreatedAt:     pbCreatedAt,
	}, nil
}

func (app *App) UpdateAccountRPC(ctx context.Context, in *pb.Account) (*pb.Account, error) {
	// Convert date_of_birth (google.type.Date) to PostgreSQL DATE format
	dob := in.GetDateOfBirth()
	var dobStr *string
	if dob != nil {
		formattedDOB := fmt.Sprintf("%04d-%02d-%02d", dob.GetYear(), dob.GetMonth(), dob.GetDay())
		dobStr = &formattedDOB
	}

	// Update the account in the database
	query := `UPDATE accounts SET email = $1, date_of_birth = $2, balance = $3 WHERE account_number = $4 RETURNING id`
	accountId := strconv.Itoa(int(in.GetId()))
	// var id int32
	err := app.DB.QueryRow(ctx, query, in.GetEmail(), dobStr, in.GetBalance(), in.GetAccountNumber()).Scan(accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	// Return the updated account
	return in, nil
}

func (app *App) DeleteAccountRPC(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.DeleteAccountResponse, error) {
	// Delete the account from the database
	query := `DELETE FROM accounts WHERE account_number = $1`
	accountId := strconv.Itoa(int(in.GetId()))
	res, err := app.DB.Exec(ctx, query, accountId)
	if err != nil {
		return nil, fmt.Errorf("failed to delete account: %w", err)
	}

	// Check if an account was deleted
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("account not found")
	}

	// Return success
	return &pb.DeleteAccountResponse{}, nil
}
