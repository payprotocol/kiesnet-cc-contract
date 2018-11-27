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

// QueryAwaitContracts _
const QueryAwaitContracts = `{
	"selector": {
		"sign.signer": "%s",
		"ccid": "%s",
		"created_time":{"$lt":"%[3]s"},
		"finished_time": {
			"$gt": "%[3]s"
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
	}, %s
}`

// ConditionAwaitUrgentContracts _
const ConditionAwaitUrgentContracts = `
"sort": [ {"sign.signer":"desc"},{"ccid":"desc"},{"finished_time":"desc"}], "use_index": ["contract_list", "contract_await_urgency"]`

// CreateQueryAwaitUrgentContracts _
func CreateQueryAwaitUrgentContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryAwaitContracts, kid, ccid, t, ConditionAwaitUrgentContracts)
}

// ConditionAwaitOldestContracts _
const ConditionAwaitOldestContracts = `
"sort": [ "sign.signer","ccid","created_time" ], "use_index": ["contract_list", "contract_await_oldest"]`

// CreateQueryAwaitOldestContracts _
func CreateQueryAwaitOldestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryAwaitContracts, kid, ccid, t, ConditionAwaitOldestContracts)
}

// QueryOngoingContracts _
const QueryOngoingContracts = `{
	"selector": {
		"sign.signer": "%s",
		"ccid": "%s",
		"finished_time": {
			"$gt": "%s"
		},
		"$and":[
			{
				"executed_time": {
					"$exists": false
				}
			},
			{
				"canceled_time": {
					"$exists": false
				}
			}
		],
		"$or":[
			{
				"sign.approved_time":{
					"$exists": true
				}
			},{
				"sign.disapproved_time":{
					"$exists": true
				}
			}
		]
	}, %s
}`

// ConditionOngoingBriskContracts _
const ConditionOngoingBriskContracts = `
"sort": [
	{ "sign.signer": "desc"},
    { "ccid": "desc" },
	{ "updated_time": "desc" }
], "use_index": ["contract_list", "contract_ongoing_brisk"]`

// CreateQueryOngoingBriskContracts _
func CreateQueryOngoingBriskContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryOngoingContracts, kid, ccid, t, ConditionOngoingBriskContracts)
}

// ConditionOngoingOldestContracts _
const ConditionOngoingOldestContracts = `
"sort": [ "sign.signer","ccid","created_time" ], "use_index": ["contract_list", "contract_ongoing_oldest"]`

// CreateQueryOngoingOldestContracts _
func CreateQueryOngoingOldestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryOngoingContracts, kid, ccid, t, ConditionOngoingOldestContracts)
}

// QueryFinContracts _
const QueryFinContracts = `{
	"selector": {
	   "sign.signer": "%s",
	   "ccid": "%s",
	   "finished_time":{
		   "$lte":"%s"
	   }
	},%s
}`

// ConditionFinLatestContracts _
const ConditionFinLatestContracts = `
"sort": [
		{"sign.signer": "desc"},
        {"ccid": "desc"},
		{"finished_time": "desc"}
   ],
	"use_index": [
	   "contract_list",
	   "contract_fin_latest"
]`

// CreateQueryFinLatestContracts _
func CreateQueryFinLatestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryFinContracts, kid, ccid, t, ConditionFinLatestContracts)
}

// ConditionFinOldestContracts _
const ConditionFinOldestContracts = `
"sort": [ "sign.signer","ccid","created_time" ],
	"use_index": [
	   "contract_list",
	   "contract_fin_oldest"
]`

// CreateQueryFinOldestContracts _
func CreateQueryFinOldestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryFinContracts, kid, ccid, t, ConditionFinOldestContracts)
}
