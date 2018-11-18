// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import "time"

// Sign represents signer's action. (approve or disapprove)
type Sign struct {
	Signer          string     `json:"signer"`
	ApprovedTime    *time.Time `json:"approved_time,omitempty"`
	DisapprovedTime *time.Time `json:"disapproved_time,omitempty"`
}
