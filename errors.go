// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import "fmt"

// NotExistedContractError _
type NotExistedContractError struct {
	id     string
	signer string
}

// Error implements error interface
func (e NotExistedContractError) Error() string {
	return fmt.Sprintf("the contract '%s' for the signer '%s' is not exists", e.id, e.signer)
}
