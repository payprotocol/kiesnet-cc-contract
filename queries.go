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
	}, %s
}`

// ConditionAwaitUrgentContracts _
const ConditionAwaitUrgentContracts = `
"sort": [
	{ "expiry_time": "asc" },
	{ "sign.signer": "asc"	},
	{ "ccid": "asc" }
], "use_index": ["contract_list", "contract_awaiter_urgent"]`

// CreateQueryAwaitUrgentContracts _
func CreateQueryAwaitUrgentContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryAwaitContracts, kid, ccid, t, ConditionAwaitUrgentContracts)
}

// ConditionAwaitOldestContracts _
const ConditionAwaitOldestContracts = `
"sort": [
	{ "created_time": "asc"	},
	{ "expiry_time":"asc" },
	{ "sign.signer": "asc"},
    { "ccid": "asc" }
], "use_index": ["contract_list", "contract_awaiter_oldest"]`

// CreateQueryAwaitOldestContracts _
func CreateQueryAwaitOldestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryAwaitContracts, kid, ccid, t, ConditionAwaitOldestContracts)
}

// QueryOngoingContracts _
const QueryOngoingContracts = `{
	"selector": {
		"sign.signer": "%s",
		"ccid": "%s",
		"expiry_time": {
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
	{ "updated_time": "desc" },
	{ "expiry_time":"desc" },
	{ "sign.signer": "desc"},
    { "ccid": "desc" }
], "use_index": ["contract_list", "contract_ongoing_brisk"]`

// CreateQueryOngoingBriskContracts _
func CreateQueryOngoingBriskContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryOngoingContracts, kid, ccid, t, ConditionOngoingBriskContracts)
}

// ConditionOngoingOldestContracts _
const ConditionOngoingOldestContracts = `
"sort": [
	{ "created_time": "asc"	},
	{ "expiry_time":"asc" },
	{ "sign.signer": "asc"},
    { "ccid": "asc" }
], "use_index": ["contract_list", "contract_ongoing_oldest"]`

// CreateQueryOngoingOldestContracts _
func CreateQueryOngoingOldestContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryOngoingContracts, kid, ccid, t, ConditionOngoingOldestContracts)
}

// QueryFinContracts _
const QueryFinContracts = `{
	"selector": {
	   "sign.signer": "%s",
	   "ccid": "%s",
	   "$or": [
		  {
			 "sign.approved_time": {
				"$exists": true
			 }
		  },
		  {
			 "sign.disapproved_time": {
				"$exists": true
			 }
		  },
		  {
			 "executed_time": {
				"$exists": true
			 }
		  },
		  {
			 "canceled_time": {
				"$exists": true
			 }
		  },
		  {
			 "$and": [
				{
				   "sign.approved_time": {
					  "$exists": false
				   }
				},
				{
				   "sign.disapproved_time": {
					  "$exists": false
				   }
				},
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
			 ]
		  }
	   ],
	   "expiry_time": {
		  "$lte": "%s"
	   }
	},
	"sort": [
		{"created_time": "desc"},
		{"expiry_time": "desc"},
        {"sign.signer": "desc"},
        {"ccid": "desc"}
   ],
	"use_index": [
	   "contract_list",
	   "contract_fin"
	]
 }`

// CreateQueryFinContracts _
func CreateQueryFinContracts(kid, ccid, t string) string {
	return fmt.Sprintf(QueryFinContracts, kid, ccid, t)
}
