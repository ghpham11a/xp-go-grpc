package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	// Import the generated protobuf package
	pb "xp-go-grpc-server/proto"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

// App holds application-wide dependencies, including the database connection pool.
type App struct {
	DB *pgxpool.Pool

	pb.UnimplementedAccountsServiceServer
}

func main() {
	// Listen on a TCP port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	host := os.Getenv("POSTGRES_HOST")
	database := os.Getenv("POSTGRES_DATABASE")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")

	testHost := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, database)
	fmt.Println("Test Host:", testHost)

	// Construct the PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, database)

	// Create a connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	app := &App{DB: pool}

	// Create a gRPC server object
	s := grpc.NewServer()

	// Test the database connection
	err = testDBConnection(pool)
	if err != nil {
		log.Fatalf("Database connection test failed: %v", err)
	}
	fmt.Println("Successfully connected to PostgreSQL!")

	// Register our service implementation with the gRPC server
	pb.RegisterAccountsServiceServer(s, app)

	fmt.Println("Server listening at port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// testDBConnection runs a simple query to check if the database connection is working.
func testDBConnection(pool *pgxpool.Pool) error {
	var currentDB string
	err := pool.QueryRow(context.Background(), "SELECT current_database()").Scan(&currentDB)
	if err != nil {
		return fmt.Errorf("failed to get current database: %w", err)
	}
	fmt.Println("Connected to Database:", currentDB)
	return nil
}
