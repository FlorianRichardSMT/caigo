package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/FlorianRichardSMT/caigo/types"
)

// AddDeclareTransactionOutput provides the output for AddDeclareTransaction.
type AddDeclareTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ClassHash       string `json:"class_hash"`
}

// AddDeployTransactionOutput provides the output for AddDeployTransaction.
type AddDeployTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
	ContractAddress string `json:"contract_address"`
}

// AddInvokeTransactionOutput provides the output for AddInvokeTransaction.
type AddInvokeTransactionOutput struct {
	TransactionHash string `json:"transaction_hash"`
}

// AddInvokeTransaction estimates the fee for a given StarkNet transaction.
func (sc *Client) AddInvokeTransaction(ctx context.Context, functionCall types.FunctionCall, signature []string, maxFee, version string) (*AddInvokeTransactionOutput, error) {
	var output AddInvokeTransactionOutput
	if err := sc.do(ctx, "starknet_addInvokeTransaction", &output, functionCall, signature, maxFee, version); err != nil {
		return nil, err
	}
	return &output, nil
}

// AddDeclareTransaction submits a new class declaration transaction.
func (sc *Client) AddDeclareTransaction(ctx context.Context, contractDefinition types.ContractClass, version string) (*AddDeclareTransactionOutput, error) {
	program, ok := contractDefinition.Program.(string)
	if !ok {
		data, err := json.Marshal(contractDefinition.Program)
		if err != nil {
			return nil, err
		}
		// TODO: change Program from contractDefinition to have a type that can handle
		// compressed and uncompressed data.
		program, err = encodeProgram(data)
		if err != nil {
			return nil, err
		}
	}
	contractDefinition.Program = program

	var result AddDeclareTransactionOutput
	if err := sc.do(ctx, "starknet_addDeclareTransaction", &result, contractDefinition, version); err != nil {
		return nil, err
	}
	return &result, nil
}

// AddDeployTransaction allows to declare a class and instantiate the
// associated contract in one command. This function will be deprecated and
// replaced by AddDeclareTransaction to declare a class, followed by
// AddInvokeTransaction to instantiate the contract. For now, it remains the only
// way to deploy an account without being charged for it.
func (sc *Client) AddDeployTransaction(ctx context.Context, contractAddressSalt string, constructorCallData []string, contractDefinition types.ContractClass) (*AddDeployTransactionOutput, error) {
	program, ok := contractDefinition.Program.(string)
	if !ok {
		data, err := json.Marshal(contractDefinition.Program)
		if err != nil {
			return nil, err
		}
		// TODO: change Program from contractDefinition to have a type that can handle
		// compressed and uncompressed data.
		program, err = encodeProgram(data)
		if err != nil {
			return nil, err
		}
	}
	contractDefinition.Program = program

	var result AddDeployTransactionOutput
	if err := sc.do(ctx, "starknet_addDeployTransaction", &result, contractAddressSalt, constructorCallData, contractDefinition); err != nil {
		return nil, err
	}
	return &result, nil
}

func encodeProgram(content []byte) (string, error) {
	buf := bytes.NewBuffer(nil)
	gzipContent := gzip.NewWriter(buf)
	_, err := gzipContent.Write(content)
	if err != nil {
		return "", err
	}
	gzipContent.Close()
	program := base64.StdEncoding.EncodeToString(buf.Bytes())
	return program, nil
}
