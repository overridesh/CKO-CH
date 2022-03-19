package model

import (
	"strings"
)

type Status string

const (
	Approved Status = "approved"
	Failed   Status = "failed"
	Pending  Status = "pending"
)

func NewStatus(raw string) Status {
	if strings.EqualFold(raw, Approved.String()) {
		return Approved
	}
	if strings.EqualFold(raw, Failed.String()) {
		return Failed
	}

	return Pending
}

func (s Status) IsApproved() bool {
	return s == Approved
}

func (s Status) String() string {
	return string(s)
}
