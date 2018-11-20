// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// QueryContractsByID _
const QueryContractsByID = `{
	"selector": {
		"@contract": "%s"
	},
	"use_index": ["contract", "id"]
}`

// CreateQueryContractsByID _
func CreateQueryContractsByID(id string) string {
	return fmt.Sprintf(QueryContractsByID, id)
}

const queryContractsBySigner = `{
	"selector":{
		"sign.signer":"%s"
	},
	"use_index":["contract","signer"]
}`

// CreateQueryContractsBySigner _
func CreateQueryContractsBySigner(kid string) string {
	return fmt.Sprintf(queryContractsBySigner, kid)
}

const queryContractsSigned = `{
	"selector":{
		"sign.signer":"%s",
		"$or":[
			{"sign.approved_time":{"$exists":true}},
			{"sign.disapproved_time":{"$exists":true}}
		]
	},
	"use_index":["contract","signed"]
}`

// CreateQueryContractsSigned _
func CreateQueryContractsSigned(kid string) string {
	return fmt.Sprintf(queryContractsSigned, kid)
}

const queryContractsUnsigned = `{
	"selector": {
		"sign.signer": "%s",
		"sign.approved_time": {"$exists": false},
		"sign.disapproved_time": {"$exists": false},
		"expiry_time": {"$gt": %s},
		"executed_time": {"$exists": false},
		"canceled_time": {"$exists": false}
	},
	"use_index":["contract","unsigned"]
 }`

// CreateQueryContractsUnsigned _
func CreateQueryContractsUnsigned(kid string, time *time.Time) (string, error) {
	timestring, err := json.Marshal(time)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(queryContractsUnsigned, kid, timestring), nil
}
