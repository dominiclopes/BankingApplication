// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	db "example.com/banking/db"

	mock "github.com/stretchr/testify/mock"
)

// Storer is an autogenerated mock type for the Storer type
type Storer struct {
	mock.Mock
}

type Storer_Expecter struct {
	mock *mock.Mock
}

func (_m *Storer) EXPECT() *Storer_Expecter {
	return &Storer_Expecter{mock: &_m.Mock}
}

// AddTransaction provides a mock function with given fields: ctx, t
func (_m *Storer) AddTransaction(ctx context.Context, t db.Transaction) error {
	ret := _m.Called(ctx, t)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, db.Transaction) error); ok {
		r0 = rf(ctx, t)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storer_AddTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AddTransaction'
type Storer_AddTransaction_Call struct {
	*mock.Call
}

// AddTransaction is a helper method to define mock.On call
//  - ctx context.Context
//  - t db.Transaction
func (_e *Storer_Expecter) AddTransaction(ctx interface{}, t interface{}) *Storer_AddTransaction_Call {
	return &Storer_AddTransaction_Call{Call: _e.mock.On("AddTransaction", ctx, t)}
}

func (_c *Storer_AddTransaction_Call) Run(run func(ctx context.Context, t db.Transaction)) *Storer_AddTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.Transaction))
	})
	return _c
}

func (_c *Storer_AddTransaction_Call) Return(err error) *Storer_AddTransaction_Call {
	_c.Call.Return(err)
	return _c
}

// CreateAccount provides a mock function with given fields: ctx, u, acc
func (_m *Storer) CreateAccount(ctx context.Context, u db.User, acc db.Account) error {
	ret := _m.Called(ctx, u, acc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, db.User, db.Account) error); ok {
		r0 = rf(ctx, u, acc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storer_CreateAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAccount'
type Storer_CreateAccount_Call struct {
	*mock.Call
}

// CreateAccount is a helper method to define mock.On call
//  - ctx context.Context
//  - u db.User
//  - acc db.Account
func (_e *Storer_Expecter) CreateAccount(ctx interface{}, u interface{}, acc interface{}) *Storer_CreateAccount_Call {
	return &Storer_CreateAccount_Call{Call: _e.mock.On("CreateAccount", ctx, u, acc)}
}

func (_c *Storer_CreateAccount_Call) Run(run func(ctx context.Context, u db.User, acc db.Account)) *Storer_CreateAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(db.User), args[2].(db.Account))
	})
	return _c
}

func (_c *Storer_CreateAccount_Call) Return(err error) *Storer_CreateAccount_Call {
	_c.Call.Return(err)
	return _c
}

// DeleteAccount provides a mock function with given fields: ctx, accID
func (_m *Storer) DeleteAccount(ctx context.Context, accID string) error {
	ret := _m.Called(ctx, accID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, accID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storer_DeleteAccount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAccount'
type Storer_DeleteAccount_Call struct {
	*mock.Call
}

// DeleteAccount is a helper method to define mock.On call
//  - ctx context.Context
//  - accID string
func (_e *Storer_Expecter) DeleteAccount(ctx interface{}, accID interface{}) *Storer_DeleteAccount_Call {
	return &Storer_DeleteAccount_Call{Call: _e.mock.On("DeleteAccount", ctx, accID)}
}

func (_c *Storer_DeleteAccount_Call) Run(run func(ctx context.Context, accID string)) *Storer_DeleteAccount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storer_DeleteAccount_Call) Return(err error) *Storer_DeleteAccount_Call {
	_c.Call.Return(err)
	return _c
}

// DepositAmount provides a mock function with given fields: ctx, accID, amount
func (_m *Storer) DepositAmount(ctx context.Context, accID string, amount float32) error {
	ret := _m.Called(ctx, accID, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, float32) error); ok {
		r0 = rf(ctx, accID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storer_DepositAmount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DepositAmount'
type Storer_DepositAmount_Call struct {
	*mock.Call
}

// DepositAmount is a helper method to define mock.On call
//  - ctx context.Context
//  - accID string
//  - amount float32
func (_e *Storer_Expecter) DepositAmount(ctx interface{}, accID interface{}, amount interface{}) *Storer_DepositAmount_Call {
	return &Storer_DepositAmount_Call{Call: _e.mock.On("DepositAmount", ctx, accID, amount)}
}

func (_c *Storer_DepositAmount_Call) Run(run func(ctx context.Context, accID string, amount float32)) *Storer_DepositAmount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(float32))
	})
	return _c
}

func (_c *Storer_DepositAmount_Call) Return(err error) *Storer_DepositAmount_Call {
	_c.Call.Return(err)
	return _c
}

// GetAccountDetails provides a mock function with given fields: ctx, accID
func (_m *Storer) GetAccountDetails(ctx context.Context, accID string) (db.Account, error) {
	ret := _m.Called(ctx, accID)

	var r0 db.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) db.Account); ok {
		r0 = rf(ctx, accID)
	} else {
		r0 = ret.Get(0).(db.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storer_GetAccountDetails_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccountDetails'
type Storer_GetAccountDetails_Call struct {
	*mock.Call
}

// GetAccountDetails is a helper method to define mock.On call
//  - ctx context.Context
//  - accID string
func (_e *Storer_Expecter) GetAccountDetails(ctx interface{}, accID interface{}) *Storer_GetAccountDetails_Call {
	return &Storer_GetAccountDetails_Call{Call: _e.mock.On("GetAccountDetails", ctx, accID)}
}

func (_c *Storer_GetAccountDetails_Call) Run(run func(ctx context.Context, accID string)) *Storer_GetAccountDetails_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storer_GetAccountDetails_Call) Return(acc db.Account, err error) *Storer_GetAccountDetails_Call {
	_c.Call.Return(acc, err)
	return _c
}

// GetAccountList provides a mock function with given fields: ctx
func (_m *Storer) GetAccountList(ctx context.Context) ([]db.Account, error) {
	ret := _m.Called(ctx)

	var r0 []db.Account
	if rf, ok := ret.Get(0).(func(context.Context) []db.Account); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Account)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storer_GetAccountList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAccountList'
type Storer_GetAccountList_Call struct {
	*mock.Call
}

// GetAccountList is a helper method to define mock.On call
//  - ctx context.Context
func (_e *Storer_Expecter) GetAccountList(ctx interface{}) *Storer_GetAccountList_Call {
	return &Storer_GetAccountList_Call{Call: _e.mock.On("GetAccountList", ctx)}
}

func (_c *Storer_GetAccountList_Call) Run(run func(ctx context.Context)) *Storer_GetAccountList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Storer_GetAccountList_Call) Return(accounts []db.Account, err error) *Storer_GetAccountList_Call {
	_c.Call.Return(accounts, err)
	return _c
}

// GetTransactions provides a mock function with given fields: ctx, accID
func (_m *Storer) GetTransactions(ctx context.Context, accID string) ([]db.Transaction, error) {
	ret := _m.Called(ctx, accID)

	var r0 []db.Transaction
	if rf, ok := ret.Get(0).(func(context.Context, string) []db.Transaction); ok {
		r0 = rf(ctx, accID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]db.Transaction)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, accID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storer_GetTransactions_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetTransactions'
type Storer_GetTransactions_Call struct {
	*mock.Call
}

// GetTransactions is a helper method to define mock.On call
//  - ctx context.Context
//  - accID string
func (_e *Storer_Expecter) GetTransactions(ctx interface{}, accID interface{}) *Storer_GetTransactions_Call {
	return &Storer_GetTransactions_Call{Call: _e.mock.On("GetTransactions", ctx, accID)}
}

func (_c *Storer_GetTransactions_Call) Run(run func(ctx context.Context, accID string)) *Storer_GetTransactions_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *Storer_GetTransactions_Call) Return(transactions []db.Transaction, err error) *Storer_GetTransactions_Call {
	_c.Call.Return(transactions, err)
	return _c
}

// GetUserByEmailAndPassword provides a mock function with given fields: ctx, email, password
func (_m *Storer) GetUserByEmailAndPassword(ctx context.Context, email string, password string) (db.User, error) {
	ret := _m.Called(ctx, email, password)

	var r0 db.User
	if rf, ok := ret.Get(0).(func(context.Context, string, string) db.User); ok {
		r0 = rf(ctx, email, password)
	} else {
		r0 = ret.Get(0).(db.User)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, email, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Storer_GetUserByEmailAndPassword_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmailAndPassword'
type Storer_GetUserByEmailAndPassword_Call struct {
	*mock.Call
}

// GetUserByEmailAndPassword is a helper method to define mock.On call
//  - ctx context.Context
//  - email string
//  - password string
func (_e *Storer_Expecter) GetUserByEmailAndPassword(ctx interface{}, email interface{}, password interface{}) *Storer_GetUserByEmailAndPassword_Call {
	return &Storer_GetUserByEmailAndPassword_Call{Call: _e.mock.On("GetUserByEmailAndPassword", ctx, email, password)}
}

func (_c *Storer_GetUserByEmailAndPassword_Call) Run(run func(ctx context.Context, email string, password string)) *Storer_GetUserByEmailAndPassword_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *Storer_GetUserByEmailAndPassword_Call) Return(u db.User, err error) *Storer_GetUserByEmailAndPassword_Call {
	_c.Call.Return(u, err)
	return _c
}

// WithdrawAmount provides a mock function with given fields: ctx, accID, amount
func (_m *Storer) WithdrawAmount(ctx context.Context, accID string, amount float32) error {
	ret := _m.Called(ctx, accID, amount)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, float32) error); ok {
		r0 = rf(ctx, accID, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Storer_WithdrawAmount_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WithdrawAmount'
type Storer_WithdrawAmount_Call struct {
	*mock.Call
}

// WithdrawAmount is a helper method to define mock.On call
//  - ctx context.Context
//  - accID string
//  - amount float32
func (_e *Storer_Expecter) WithdrawAmount(ctx interface{}, accID interface{}, amount interface{}) *Storer_WithdrawAmount_Call {
	return &Storer_WithdrawAmount_Call{Call: _e.mock.On("WithdrawAmount", ctx, accID, amount)}
}

func (_c *Storer_WithdrawAmount_Call) Run(run func(ctx context.Context, accID string, amount float32)) *Storer_WithdrawAmount_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(float32))
	})
	return _c
}

func (_c *Storer_WithdrawAmount_Call) Return(err error) *Storer_WithdrawAmount_Call {
	_c.Call.Return(err)
	return _c
}

type mockConstructorTestingTNewStorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewStorer creates a new instance of Storer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStorer(t mockConstructorTestingTNewStorer) *Storer {
	mock := &Storer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}