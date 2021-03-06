// Code generated by MockGen. DO NOT EDIT.
// Source: currency_parsing.go

// Package mock_pkg is a generated GoMock package.
package mock_pkg

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCurrencyCalculator is a mock of CurrencyCalculator interface.
type MockCurrencyCalculator struct {
	ctrl     *gomock.Controller
	recorder *MockCurrencyCalculatorMockRecorder
}

// MockCurrencyCalculatorMockRecorder is the mock recorder for MockCurrencyCalculator.
type MockCurrencyCalculatorMockRecorder struct {
	mock *MockCurrencyCalculator
}

// NewMockCurrencyCalculator creates a new mock instance.
func NewMockCurrencyCalculator(ctrl *gomock.Controller) *MockCurrencyCalculator {
	mock := &MockCurrencyCalculator{ctrl: ctrl}
	mock.recorder = &MockCurrencyCalculatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCurrencyCalculator) EXPECT() *MockCurrencyCalculatorMockRecorder {
	return m.recorder
}

// ConvertRubTo mocks base method.
func (m *MockCurrencyCalculator) ConvertRubTo(currency string, sum float32) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConvertRubTo", currency, sum)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ConvertRubTo indicates an expected call of ConvertRubTo.
func (mr *MockCurrencyCalculatorMockRecorder) ConvertRubTo(currency, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConvertRubTo", reflect.TypeOf((*MockCurrencyCalculator)(nil).ConvertRubTo), currency, sum)
}

// UpdateRates mocks base method.
func (m *MockCurrencyCalculator) UpdateRates() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateRates")
}

// UpdateRates indicates an expected call of UpdateRates.
func (mr *MockCurrencyCalculatorMockRecorder) UpdateRates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRates", reflect.TypeOf((*MockCurrencyCalculator)(nil).UpdateRates))
}
