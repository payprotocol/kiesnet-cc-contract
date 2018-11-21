// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("kiesnet-contract")

// Chaincode _
type Chaincode struct {
}

// Init implements shim.Chaincode interface.
func (cc *Chaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke implements shim.Chaincode interface.
func (cc *Chaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, params := stub.GetFunctionAndParameters()
	if txFn := routes[fn]; txFn != nil {
		return txFn(stub, params)
	}
	return shim.Error("unknown function: '" + fn + "'")
}

// TxFunc _
type TxFunc func(shim.ChaincodeStubInterface, []string) peer.Response

// routes is the map of invoke functions
var routes = map[string]TxFunc{
	"approve":    contractApprove,
	"cancel":     contractCancel,
	"disapprove": contractDisapprove,
	"get":        contractGet,
	"list":       contractList,
	"new":        contractNew,
	"ver":        ver,
}

func ver(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	return shim.Success([]byte("Kiesnet Contract v1.0 created by Key Inside Co., Ltd."))
}

func response(payload Payload) peer.Response {
	data, err := payload.MarshalPayload()
	if err != nil {
		logger.Debug(err.Error())
		return shim.Error("failed to marshal payload")
	}
	return shim.Success(data)
}

func main() {
	if err := shim.Start(new(Chaincode)); err != nil {
		logger.Criticalf("failed to start chaincode: %s", err)
	}
}
