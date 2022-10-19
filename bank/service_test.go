package bank

import (
	"context"
	"errors"
	"testing"

	uuidgen "github.com/pborman/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"example.com/banking/app"
	"example.com/banking/db"
	"example.com/banking/db/mocks"
)

func init() {
	app.InitLogger()
}

type BankServiceTestSuite struct {
	suite.Suite
	logger      *zap.SugaredLogger
	storer      *mocks.Storer
	bankService Service
}

func (bsts *BankServiceTestSuite) SetupSuite() {
	bsts.T().Logf("SetupSuite - Creating the logger instance")
	bsts.logger = app.GetLogger()
}

func (bsts *BankServiceTestSuite) SetupTest() {
	bsts.T().Logf("SetupTest - Creating the mock db instance and the bank service")

	bsts.storer = mocks.NewStorer(bsts.T())
	bsts.bankService = NewBankService(bsts.storer, bsts.logger)
}

func TestBankServiceTestSuite(t *testing.T) {
	suite.Run(t, &BankServiceTestSuite{})
}

func (bsts *BankServiceTestSuite) Test_bankService_CreateAccount() {
	type args struct {
		ctx    context.Context
		accReq CreateAccountRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// positive test
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
				accReq: CreateAccountRequest{
					Email:       "abc@gmail.com",
					PhoneNumber: "1234567899",
				},
			},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("CreateAccount", context.TODO(), mock.AnythingOfType("db.User"), mock.AnythingOfType("db.Account")).Return(nil).Once()
			},
		},
		// negative test
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
				accReq: CreateAccountRequest{
					Email:       "abc@gmail.com",
					PhoneNumber: "1234567899",
				},
			},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("CreateAccount", context.TODO(), mock.AnythingOfType("db.User"), mock.AnythingOfType("db.Account")).Return(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		bsts.T().Run(tt.name, func(t *testing.T) {

			tt.prepare(tt.args, bsts.storer)

			gotAccRes, err := bsts.bankService.CreateAccount(tt.args.ctx, tt.args.accReq)

			if tt.wantErr {
				// bsts.
				bsts.ErrorContains(err, "mocked error")
			} else {
				bsts.ErrorIs(err, nil)
			}

			bsts.IsType(CreateAccountResponse{}, gotAccRes)
		})
	}
}

func (bsts *BankServiceTestSuite) Test_bankService_GetAccountList() {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantAccounts []db.UserAccountDetails
		prepare      func(args, *mocks.Storer)
	}{
		//positive test
		{
			name: "positiveTest",
			args: args{
				ctx: context.TODO(),
			},
			wantErr:      false,
			wantAccounts: []db.UserAccountDetails{},
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountList", context.TODO()).Return([]db.UserAccountDetails{}, nil).Once()
			},
		},
		//negative test
		{
			name: "negativeTest",
			args: args{
				ctx: context.TODO(),
			},
			wantErr:      true,
			wantAccounts: nil,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountList", context.TODO()).Return(nil, errors.New("mocked error"))
			},
		},
	}

	for _, tt := range tests {
		bsts.T().Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, bsts.storer)

			accounts, err := bsts.bankService.GetAccountList(tt.args.ctx)

			if tt.wantErr {
				// bsts.
				bsts.ErrorContains(err, "mocked error")
				bsts.Nil(accounts)
			} else {
				bsts.ErrorIs(err, nil)
				bsts.Equal(tt.wantAccounts, accounts)
			}
		})
	}
}

func (bsts *BankServiceTestSuite) Test_bankService_GetAccountDetails() {
	type args struct {
		ctx    context.Context
		accId  string
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantAcc db.UserAccountDetails
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// positive test
		{
			name:    "positiveTest",
			args:    args{context.TODO(), uuidgen.New(), "1"},
			wantAcc: db.UserAccountDetails{},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountDetails", a.ctx, a.accId, a.userID).Return(db.UserAccountDetails{}, nil)
			},
		},
		// negatice test
		{
			name:    "negativeTest",
			args:    args{context.TODO(), uuidgen.New(), "2"},
			wantAcc: db.UserAccountDetails{},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountDetails", a.ctx, a.accId, a.userID).Return(db.UserAccountDetails{}, errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		bsts.T().Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args, bsts.storer)

			gotAcc, err := bsts.bankService.GetAccountDetails(tt.args.ctx, tt.args.accId, tt.args.userID)

			if tt.wantErr {
				// bsts.
				bsts.ErrorContains(err, "mocked error")
			} else {
				bsts.ErrorIs(err, nil)
			}
			bsts.Equal(tt.wantAcc, gotAcc)
		})
	}
}
