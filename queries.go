// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import "fmt"

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

// QueryAllContractsByKID _
const QueryAllContractsByKID = `{
	"selector": {
		"sign.signer":"%s"
	}
}`

//CreateQueryAllContractsByKID _
func CreateQueryAllContractsByKID(kid string) string {
	return fmt.Sprintf(QueryAllContractsByKID, kid)
}

// QueryContractsNeedSign _
const QueryContractsNeedSign = `{
	"selector": {
		"sign.signer":"%s",
		"$exists":{
			"sign.ApprovedTime": false
		}
	}
}`

//CreateQueryContractsNeedSign _
func CreateQueryContractsNeedSign(kid string) string {
	return fmt.Sprintf(QueryContractsNeedSign, kid)
}
