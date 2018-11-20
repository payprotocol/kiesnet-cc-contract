// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/json"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/key-inside/kiesnet-ccpkg/ccid"
	"github.com/key-inside/kiesnet-ccpkg/kid"
	"github.com/key-inside/kiesnet-ccpkg/stringset"
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

	if contract.SignersCount == contract.ApprovedCount {
		// TODO: invoke Execute!
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

	// TODO: invoke Cancel!

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

// params[0] : stateCriteria (all, unsigned, signed)
// unsigned : 서멍해야 하는거 (서명 못하는 상태는 뺌)
// signed : 서명한 것 (상태가 어떻든 보여줌)
// params[1] : bookmark
//TODO ccid : 없으면 전체
//XXX lifetimeCriteria (all, pending, expired, executed, canceled)
//XXX creator : 의미없음
//TODO sort : created_time, update_time, executed_time,
func contractList(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	// authentication
	kid, err := kid.GetID(stub, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	stateCriteria := "all"
	if len(params) > 0 {
		stateCriteria = params[0]
	}
	bookmark := ""
	if len(params) > 1 {
		bookmark = params[1]
	}

	cb := NewContractStub(stub)
	res, err := cb.GetQueryContractsResult(kid, stateCriteria, bookmark)
	if err != nil {
		logger.Debug(err.Error())
		return shim.Error("failed to get account addresses list")
	}

	data, err := json.Marshal(res)
	if err != nil {
		logger.Debug(err.Error())
		return shim.Error("failed to marshal account addresses list")
	}
	return shim.Success(data)
}

// params[0] : document (JSON string)
// params[1] : expiry (duration represented by int64 seconds, multi-sig only)
// params[2:] : signers' KID (exclude invoker, max 127)
func contractNew(stub shim.ChaincodeStubInterface, params []string) peer.Response {
	ccid, err := ccid.GetID(stub)
	if err != nil || "kiesnet-contract" == ccid || "kiesnet-cc-contract" == ccid {
		return shim.Error("invalid access")
	}

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
