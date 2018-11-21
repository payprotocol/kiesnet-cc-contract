// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"fmt"
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

// QueryAllContracts _
const QueryAllContracts = `{
	"selector": {
		"sign.signer":"%s",
		"ccid": "%s"
	},
	"use_index": ["contract_list", "contract_all"]
}`

// CreateQueryAllContracts _
func CreateQueryAllContracts(kid string, ccid string) string {
	return fmt.Sprintf(QueryAllContracts, kid, ccid)
}

// QueryAwaitContracts _
const QueryAwaitContracts = `{
	"selector": {
		"sign.signer": "%s",
		"ccid": "%s",
		"expiry_time": {
			"$gt": "%s"
		},
		"$and":[
			{
				"sign.approved_time":{
					"$exists": false
				}
			},{
				"sign.disapproved_time":{
					"$exists": false
				}
			},{
				"executed_time": {
					"$exists": false
				}
			},
			{
				"canceled_time": {
					"$exists": false
				}
			}
		]
	}, 
	"use_index":["contract_list", "contract_awaiter"]
}`

// CreateQueryAwaitContracts _
func CreateQueryAwaitContracts(kid string, ccid string, t string) string {
	return fmt.Sprintf(QueryAwaitContracts, kid, ccid, t)
}

// QueryFinContracts _
const QueryFinContracts = `{
	"selector": {
		"sign.signer": "%s",
		"ccid": "%s",
		"$or":[
			{
				"sign.approved_time":{
					"$exists": true
				}
			},
			{
				"sign.disapproved_time":{
					"$exists": true
				}
			},
			{
				"executed_time":{
					"$exists": true
				}
			},
			{
				"canceled_time":{
					"$exists": true
				}
			}
		]
	}, 
	"use_index":["contract_list", "contract_fin"]
}`

// CreateQueryFinContracts _
func CreateQueryFinContracts(kid string, ccid string) string {
	return fmt.Sprintf(QueryFinContracts, kid, ccid)
}
