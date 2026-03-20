package main

import (
	"log"
	"os"

	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/dadata"
	"studentgit.kata.academy/ar.konovalov202_gmail.com/rpc/rpcserver"
)

func main() {
	apiKey := os.Getenv("DADATA_API_KEY")
	secretKey := os.Getenv("DADATA_SECRET_KEY")

	dadataClient, err := dadata.NewClient(apiKey, secretKey)
	if err != nil {
		log.Fatalf("Failed to create DaData client: %v", err)
	}

	rpcserver.StartRpcServer(dadataClient)
}
