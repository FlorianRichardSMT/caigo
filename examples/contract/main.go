package main

import (
	"context"
	"fmt"
	"time"

	"github.com/FlorianRichardSMT/caigo/gateway"
	"github.com/FlorianRichardSMT/caigo/types"
)

// Start Devnet:
// 	- starknet-devnet
var (
	name         string = "local"
	contract     string = "../../rpc/tests/counter.json"
	maxPoll      int    = 5
	pollInterval int    = 5
)

func main() {
	// init starknet gateway client
	gw := gateway.NewClient(gateway.WithChain(name))

	// will fail w/o new seed
	deployResponse, err := gw.Deploy(context.Background(), contract, types.DeployRequest{
		ContractAddressSalt: fmt.Sprintf("0x%x", time.Now().UnixNano()),
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deployment Response: \n\t%+v\n\n", deployResponse)

	// poll until the desired transaction status
	n, receipt, err := gw.PollTx(context.Background(), deployResponse.TransactionHash, types.ACCEPTED_ON_L2, pollInterval, maxPoll)
	if err != nil {
		fmt.Println("Transaction Failure: ", receipt.Status)
		panic(err.Error())
	}
	fmt.Printf("Poll %dsec %dx \n\ttransaction(%s) status: %s\n\n", n*pollInterval, n, deployResponse.TransactionHash, receipt.Status)

	// fetch transaction data
	tx, err := gw.Transaction(context.Background(), gateway.TransactionOptions{TransactionHash: deployResponse.TransactionHash})
	if err != nil {
		panic(err.Error())
	}

	// call StarkNet contract
	callResp, err := gw.Call(context.Background(), types.FunctionCall{
		ContractAddress:    tx.Transaction.ContractAddress,
		EntryPointSelector: "get_rand",
	}, "")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Counter is currently at: ", callResp[0])
}
