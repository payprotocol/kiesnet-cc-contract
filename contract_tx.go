// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/key-inside/kiesnet-ccpkg/kid"
	"github.com/key-inside/kiesnet-ccpkg/stringset"
	"github.com/pkg/errors"
)

// params[0] : contract ID
func contractApprove(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return shim.Error(err.Error())
	}
	contract, err = cb.ApproveContract(contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	if contract.ExecutedTime != nil {
		// execute contract
		if err = invokeExecuteContract(stub, contract); err != nil {
			return shim.Error("failed to approve the contract: " + err.Error())
		}
	}

	data, err := json.Marshal(contract)
	if err != nil {
		return shim.Error("failed to marshal the contract")
	}
	return shim.Success(data)
}

// params[0] : contract ID
func contractCancel(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}
	// TODO
	return shim.Success([]byte("cancel"))
}

// params[0] : document (JSON string)
// params[1] : expiry (duration represented by int64 seconds, multi-sig only)
// params[2:] : signers' KID (exclude invoker, max 127)
func contractCreate(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	// ccid, err := ccid.GetID(stub)
	// if err != nil || "kiesnet-contract" == ccid || "kiesnet-cc-contract" == ccid {
	// 	return shim.Error("invalid access")
	// }
	// TODO: ... fix bug
	ccid := "kiesnet-cc-token"

	if len(params) < 3 {
		return shim.Error("incorrect number of parameters. expecting 3+")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	signers := stringset.New(kid)
	signers.AppendSlice(params[2:])

	if signers.Size() < 2 {
		return shim.Error("not enough signers")
	} else if signers.Size() > 128 {
		return shim.Error("too many signers")
	}

	expiry, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		expiry = 0
	}

	document := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.CreateContracts(kid, ccid, document, signers, expiry)
	if err != nil {
		logger.Debug(err.Error())
		return shim.Error("failed to create contracts")
	}

	data, err := json.Marshal(contract)
	if err != nil {
		return shim.Error("failed to marshal the contract")
	}
	return shim.Success(data)
}

// params[0] : contract ID
func contractDisapprove(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return shim.Error(err.Error())
	}
	contract, err = cb.DisapproveContract(contract)
	if err != nil {
		return shim.Error(err.Error())
	}

	// cancel contract
	if err = invokeCancelContract(stub, contract); err != nil {
		return shim.Error("failed to cancel the contract: " + err.Error())
	}

	data, err := json.Marshal(contract)
	if err != nil {
		return shim.Error("failed to marshal the contract")
	}
	return shim.Success(data)
}

// params[0] : contract ID
func contractGet(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters. expecting 1")
	}

	// authentication
	kid, err := kid.GetID(stub, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	id := params[0]

	cb := NewContractStub(stub)
	contract, err := cb.GetContract(id, kid)
	if err != nil {
		return shim.Error(err.Error())
	}

	data, err := json.Marshal(contract)
	if err != nil {
		return shim.Error("failed to marshal the contract")
	}
	return shim.Success(data)
}

// params[0] : option
// params[1] : bookmark
func contractList(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	/*
		TODO:
		2) ccid handling
			ASIS : exact case search
	*/
	if len(params) < 2 {
		return shim.Error("incorrect number of parameters. expecting 2")
	}
	option := params[0]
	b := params[1]
	cb := NewContractStub(stub)
	res, err := cb.GetContractList(option, b)
	if nil != err {
		return shim.Error(err.Error())
	}

	return response(res)
}

// helpers

func invokeCallback(stub shim.ChaincodeStubInterface, ccid string, args [][]byte) error {
	res := stub.InvokeChaincode(ccid, args, "")
	if res.GetStatus() == 200 {
		return nil
	}
	return errors.New(res.GetMessage())
}

func invokeExecuteContract(stub shim.ChaincodeStubInterface, contract *Contract) error {
	args := [][]byte{[]byte("contract/execute"), []byte(contract.DOCTYPEID), []byte(contract.Document)}
	return invokeCallback(stub, contract.CCID, args)
}

func invokeCancelContract(stub shim.ChaincodeStubInterface, contract *Contract) error {
	args := [][]byte{[]byte("contract/cancel"), []byte(contract.DOCTYPEID), []byte(contract.Document)}
	return invokeCallback(stub, contract.CCID, args)
}
