// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/key-inside/kiesnet-ccpkg/ccid"
	"github.com/key-inside/kiesnet-ccpkg/kid"
	"github.com/key-inside/kiesnet-ccpkg/stringset"
	"github.com/key-inside/kiesnet-ccpkg/txtime"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
)

// ContractStub _
type ContractStub struct {
	stub shim.ChaincodeStubInterface
}

// NewContractStub _
func NewContractStub(stub shim.ChaincodeStubInterface) *ContractStub {
	return &ContractStub{stub}
}

// CreateKey _
func (cb *ContractStub) CreateKey(id, signer string) string {
	return fmt.Sprintf("CTR_%s_%s", id, signer)
}

// CreateHash _
func (cb *ContractStub) CreateHash(text string) string {
	h := make([]byte, 32)
	sha3.ShakeSum256(h, []byte(text))
	return hex.EncodeToString(h)
}

// CreateContracts _
func (cb *ContractStub) CreateContracts(creator, ccid, document string, signers stringset.Set, expiry int64) (*Contract, error) {
	ts, err := txtime.GetTime(cb.stub)
	if err != nil {
		return nil, err
	}

	scount := signers.Size()
	var expTime time.Time
	if expiry > 0 {
		expTime = ts.Add(time.Second * time.Duration(expiry))
	} else { // default 15 days
		expTime = ts.AddDate(0, 0, 15)
	}

	id := cb.CreateHash(creator + cb.stub.GetTxID())
	// check id collision
	query := CreateQueryContractsByID(id)
	iter, err := cb.stub.GetQueryResult(query)
	if err != nil {
		return nil, err
	}
	defer iter.Close()
	if iter.HasNext() {
		return nil, errors.New("contract ID collided")
	}

	var _contract *Contract // creator's contract (for return)

	for signer := range signers {
		sign := &Sign{
			Signer: signer,
		}
		contract := &Contract{
			DOCTYPEID:     id,
			Creator:       creator,
			SignersCount:  scount,
			ApprovedCount: 1,
			CCID:          ccid,
			Document:      document,
			CreatedTime:   ts,
			UpdatedTime:   ts,
			ExpiryTime:    &expTime,
			Sign:          sign,
		}
		if creator == signer {
			sign.ApprovedTime = ts
			_contract = contract
		}
		if err = cb.PutContract(contract); err != nil {
			return nil, err
		}
	}

	return _contract, nil
}

// GetContract _
func (cb *ContractStub) GetContract(id, signer string) (*Contract, error) {
	data, err := cb.stub.GetState(cb.CreateKey(id, signer))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the contract state")
	}
	if data != nil {
		contract := &Contract{}
		if err = json.Unmarshal(data, contract); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal the contract")
		}
		return contract, nil
	}
	return nil, NotExistedContractError{id, signer}
}

// PutContract _
func (cb *ContractStub) PutContract(contract *Contract) error {
	data, err := json.Marshal(contract)
	if err != nil {
		return errors.Wrap(err, "failed to marshal the contract")
	}
	key := cb.CreateKey(contract.DOCTYPEID, contract.Sign.Signer)
	if err = cb.stub.PutState(key, data); err != nil {
		return errors.Wrap(err, "failed to put the contract state")
	}
	return nil
}

// ApproveContract _
func (cb *ContractStub) ApproveContract(contract *Contract) (*Contract, error) {
	ts, err := txtime.GetTime(cb.stub)
	if err != nil {
		return nil, err
	}

	if err = contract.AssertSignable(*ts); err != nil {
		return nil, err
	}

	contract.Sign.ApprovedTime = ts
	contract.UpdatedTime = ts
	contract.ApprovedCount++
	if contract.SignersCount == contract.ApprovedCount {
		contract.ExecutedTime = ts
	}

	// update all other signers
	if err = cb.UpdateContracts(contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// DisapproveContract _
func (cb *ContractStub) DisapproveContract(contract *Contract) (*Contract, error) {
	ts, err := txtime.GetTime(cb.stub)
	if err != nil {
		return nil, err
	}

	if err = contract.AssertSignable(*ts); err != nil {
		return nil, err
	}

	contract.Sign.DisapprovedTime = ts
	contract.UpdatedTime = ts
	contract.CanceledTime = ts

	// update all other signers
	if err = cb.UpdateContracts(contract); err != nil {
		return nil, err
	}

	return contract, nil
}

// UpdateContracts updates contracts with values of updater
func (cb *ContractStub) UpdateContracts(updater *Contract) error {
	query := CreateQueryContractsByID(updater.DOCTYPEID)
	iter, err := cb.stub.GetQueryResult(query)
	if err != nil {
		return errors.Wrap(err, "failed to query contracts")
	}
	defer iter.Close()

	_copy := *updater // copy
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return errors.Wrap(err, "failed to get a contract")
		}
		updatee := &Contract{}
		if err = json.Unmarshal(kv.Value, updatee); err != nil {
			return errors.Wrap(err, "failed to unmarshal the contract")
		}
		if updatee.Sign.Signer != updater.Sign.Signer {
			_copy.Sign = updatee.Sign // switch signer
		} else {
			_copy.Sign = updater.Sign
		}
		if err = cb.PutContract(&_copy); err != nil {
			return errors.Wrap(err, "failed to update a contract")
		}
	}

	return nil
}

// ContractListFetchSize _
const ContractListFetchSize = 20

// GetContractList  _
func (cb *ContractStub) GetContractList(option, bookmark string) (*QueryResult, error) {
	kid, err := kid.GetID(cb.stub, false)
	if nil != err {
		return nil, err
	}
	ccid, err := ccid.GetID(cb.stub)
	if nil != err {
		return nil, err
	}
	query := ""
	ts, err := txtime.GetTime(cb.stub)
	if nil != err {
		return nil, err
	}
	t := ts.Format(time.RFC3339)
	switch option {
	case "await.urgency":
		query = CreateQueryAwaitUrgentContracts(kid, ccid, t)
	case "await.oldest":
		query = CreateQueryAwaitOldestContracts(kid, ccid, t)
	case "ongoing.brisk":
		query = CreateQueryOngoingBriskContracts(kid, ccid, t)
	case "ongoing.oldest":
		query = CreateQueryOngoingOldestContracts(kid, ccid, t)
	case "fin":
		query = CreateQueryFinContracts(kid, ccid, t)
	default:
		query = CreateQueryAwaitUrgentContracts(kid, ccid, t)
	}
	fmt.Println(query)
	// Bookmark - little bit different..too long
	iter, meta, err := cb.stub.GetQueryResultWithPagination(query, ContractListFetchSize, bookmark)
	if nil != err {
		return nil, err
	}
	// Issue... Handling custom error or normal case
	if meta.GetFetchedRecordsCount() == 0 {
		return nil, NoFetchRecordsCountError{}
	}
	defer iter.Close()
	result, err := NewQueryResult(meta, iter)
	if nil != err {
		return nil, err
	}
	return result, nil
}
