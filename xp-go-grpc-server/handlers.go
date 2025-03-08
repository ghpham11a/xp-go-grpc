package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Import the generated protobuf package
	pb "xp-go-grpc-server/proto"
)

func (s *server) CreateAccountRPC(ctx context.Context, in *pb.Account) (*pb.Account, error) {
	// Safely access the date_of_birth, which is of type *date.Date
	dob := in.GetDateOfBirth()
	// If dob is non-nil, break out year/month/day
	var dobStr string
	if dob != nil {
		dobStr = fmt.Sprintf("%04d-%02d-%02d", dob.GetYear(), dob.GetMonth(), dob.GetDay())
	} else {
		dobStr = "(nil)"
	}

	// Safely access the created_at, which is a *timestamp.Timestamp
	createdAt := in.GetCreatedAt()
	var createdAtStr string
	if createdAt != nil {
		// Convert to time.Time for readable output
		createdAtStr = createdAt.AsTime().Format(time.RFC3339)
	} else {
		createdAtStr = "(nil)"
	}

	log.Printf(
		"Received Account:\n"+
			"ID: %d\n"+
			"Email: %s\n"+
			"DateOfBirth: %s\n"+
			"AccountNumber: %s\n"+
			"Balance: %s\n"+
			"CreatedAt: %s\n",
		in.GetId(),
		in.GetEmail(),
		dobStr,
		in.GetAccountNumber(),
		in.GetBalance(),
		createdAtStr,
	)

	// Return something; for example, echo back the account
	return &pb.Account{
		Id:            in.GetId(),
		Email:         in.GetEmail(),
		DateOfBirth:   in.GetDateOfBirth(),
		AccountNumber: in.GetAccountNumber(),
		Balance:       in.GetBalance(),
		CreatedAt:     in.GetCreatedAt(),
	}, nil
}

func (s *server) GetAccountRPC(ctx context.Context, in *pb.Account) (*pb.Account, error) {
	log.Printf("Received request for account number=%s", in.GetAccountNumber())
	return &pb.Account{AccountNumber: in.GetAccountNumber()}, nil
}
