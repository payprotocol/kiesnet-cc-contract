// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Contract represents the contract
type Contract struct {
	DOCTYPEID     string     `json:"@contract"`
	Creator       string     `json:"creator"`
	SignersCount  int        `json:"signers_count"`
	ApprovedCount int        `json:"approved_count"`
	CCID          string     `json:"ccid"`
	Document      string     `json:"document"`
	CreatedTime   *time.Time `json:"created_time,omitempty"`
	UpdatedTime   *time.Time `json:"updated_time,omitempty"`
	ExpiryTime    *time.Time `json:"expiry_time,omitempty"`
	ExecutedTime  *time.Time `json:"executed_time,omitempty"`
	CanceledTime  *time.Time `json:"canceled_time,omitempty"`
	FinishedTime  *time.Time `json:"finished_time,omitempty"`
	Sign          *Sign      `json:"sign"`
}

// AssertSignable _
func (c *Contract) AssertSignable(t time.Time) error {
	if c.ExecutedTime != nil {
		return errors.New("already executed")
	}
	if c.CanceledTime != nil {
		return errors.New("already canceled")
	}
	if c.ExpiryTime != nil && !c.ExpiryTime.After(t) { // t == expiry => expired
		return errors.New("already expired")
	}
	if c.Sign.ApprovedTime != nil {
		return errors.New("already approved")
	}
	if c.Sign.DisapprovedTime != nil {
		return errors.New("already dispproved")
	}
	return nil
}

// MarshalPayload _
func (c *Contract) MarshalPayload() ([]byte, error) {
	return json.Marshal(c)
}

// // CombinedTimeSet represents ...
// type CombinedTimeSet struct {
// 	// finishedtime + createtime =>
// 	/*
// 				   c			f			e			cts
// 				18880101	20011111	20011111
// 	   a		20001010	20001021	20001030	20001010 21			11	10 + 21 21
// 	   b		20001011	20001025	20001031	20001011 25			14  11 + 25 36
// 	   c		20001012	20001023	20001032	20001012 23			11	12 + 23 35
// 	   d		20001013	20001030	20001033    20001013 30			17	13 + 30 43
// 	   e		20001014	20001026	20001026(x) 20001014 26	12		12	14 + 26 40 -12 28
// 					  14.5	      15          26          15                14.5+15 29.5
// 	   f		20001015	20001024	20001024(x) 20001015 24	9       9	15 + 24 39 -9  30
// 	   g		20001025	20001025	20001035	20001025 25 0		0	25 + 25 50

// 	 			20001025	20001040	20001040	20001025 40 15			25 + 40 65

// 	   h		20001030	20001040	20001040	20001030 40 *			30 + 40 70

// 							20001031				20001031 31				31 + 31 62
// 													20001025 20001031
// 	   a-b-c-d     a-c-b-d

// 	*/

// }
