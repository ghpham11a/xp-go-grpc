package main

import (
	"log"
	"net/http"
	"os"

	pb "xp-go-grpc-client/proto"

	"github.com/gorilla/mux"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	AccountsClient pb.AccountsServiceClient
}

func main() {

	serverDNS := os.Getenv("GO_GRPC_SERVER_DNS")

	// Set up a connection to the server.
	conn, err := grpc.NewClient(serverDNS, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new Account client.
	ac := pb.NewAccountsServiceClient(conn)

	app := &App{AccountsClient: ac}

	// Create a new router
	r := mux.NewRouter()

	r.HandleFunc("/accounts", app.CreateAccount).Methods("POST")
	r.HandleFunc("/accounts/{accountId}", app.ReadAccount).Methods("GET")
	r.HandleFunc("/accounts/{accountId}", app.UpdateAccount).Methods("PUT")
	r.HandleFunc("/accounts/{accountId}", app.DeleteAccount).Methods("DELETE")

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
