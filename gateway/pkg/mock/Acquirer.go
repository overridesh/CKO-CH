// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	acquirer "github.com/overridesh/checkout-challenger/pkg/service/acquirer"
	mock "github.com/stretchr/testify/mock"
)

// Acquirer is an autogenerated mock type for the Acquirer type
type Acquirer struct {
	mock.Mock
}

// Purchase provides a mock function with given fields: _a0
func (_m *Acquirer) Purchase(_a0 *acquirer.PaymentRequest) (*acquirer.PaymentResponse, error) {
	ret := _m.Called(_a0)

	var r0 *acquirer.PaymentResponse
	if rf, ok := ret.Get(0).(func(*acquirer.PaymentRequest) *acquirer.PaymentResponse); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*acquirer.PaymentResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*acquirer.PaymentRequest) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}