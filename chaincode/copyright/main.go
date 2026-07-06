package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"log"
)

func main() {
	cc, err := contractapi.NewChaincode(&CopyrightContract{})
	if err != nil {
		log.Panicf("Error creating copyright chaincode: %v", err)
	}
	if err := cc.Start(); err != nil {
		log.Panicf("Error starting copyright chaincode: %v", err)
	}
}
