// Copyright Key Inside Co., Ltd. 2018 All Rights Reserved.

package main

import (
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
	if t.After(*c.ExpiryTime) {
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
