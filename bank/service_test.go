bank/mocks/ package bank

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"example.com/banking/db"
	"example.com/banking/db/mocks"
	uuidgen "github.com/pborman/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_bankService_CreateAccount(t *testing.T) {
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
			name:    "positiveTest",
			args:    args{ctx: context.TODO(), accReq: CreateAccountRequest{Email: "abc@gmail.com", PhoneNumber: "1234567899"}},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("CreateAccount", context.TODO(), mock.AnythingOfType("db.User"), mock.AnythingOfType("db.Account")).Return(nil)
			},
		},
		// negative test
		{
			name:    "negativeTest",
			args:    args{ctx: context.TODO(), accReq: CreateAccountRequest{Email: "abc@gmail.com", PhoneNumber: "1234567899"}},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("CreateAccount", context.TODO(), mock.AnythingOfType("db.User"), mock.AnythingOfType("db.Account")).Return(errors.New("mocked error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			logger := zapLogger.Sugar()

			s := mocks.NewStorer(t)

			b := NewBankService(s, logger)
			tt.prepare(tt.args, s)

			gotAccRes, err := b.CreateAccount(tt.args.ctx, tt.args.accReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("bankService.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if reflect.TypeOf(gotAccRes) != reflect.TypeOf(CreateAccountResponse{}) {
				t.Errorf("Incorrect response type")
				return
			}
		})
	}
}

func Test_bankService_GetAccountList(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantAccounts []db.Account
		prepare      func(args, *mocks.Storer)
	}{
		//positive test
		{
			name:         "positiveTest",
			args:         args{ctx: context.TODO()},
			wantErr:      false,
			wantAccounts: []db.Account{},
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountList", context.TODO()).Return([]db.Account{}, nil)
			},
		},
		//negative test
		{
			name:         "negativeTest",
			args:         args{ctx: context.TODO()},
			wantErr:      true,
			wantAccounts: nil,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountList", context.TODO()).Return(nil, errors.New("My error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			zapLogger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			logger := zapLogger.Sugar()

			s := mocks.NewStorer(t)

			b := NewBankService(s, logger)
			tt.prepare(tt.args, s)

			accounts, err := b.GetAccountList(tt.args.ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("bankService.GetAccountList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(accounts, tt.wantAccounts) {
				t.Errorf("Incorrect response type")
				return
			}
		})
	}
}

func Test_bankService_GetAccountDetails(t *testing.T) {
	type args struct {
		ctx   context.Context
		accId string
	}
	tests := []struct {
		name    string
		args    args
		wantAcc db.Account
		wantErr bool
		prepare func(args, *mocks.Storer)
	}{
		// positive test
		{
			name:    "positiveTest",
			args:    args{context.TODO(), uuidgen.New()},
			wantAcc: db.Account{},
			wantErr: false,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountDetails", a.ctx, a.accId).Return(db.Account{}, nil)
			},
		},
		// negatice test
		{
			name:    "negativeTest",
			args:    args{context.TODO(), uuidgen.New()},
			wantAcc: db.Account{},
			wantErr: true,
			prepare: func(a args, s *mocks.Storer) {
				s.On("GetAccountDetails", a.ctx, a.accId).Return(db.Account{}, errors.New("my error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			zapLogger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			logger := zapLogger.Sugar()

			s := mocks.NewStorer(t)

			b := NewBankService(s, logger)
			tt.prepare(tt.args, s)

			gotAcc, err := b.GetAccountDetails(tt.args.ctx, tt.args.accId)
			if (err != nil) != tt.wantErr {
				t.Errorf("bankService.GetAccountDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAcc, tt.wantAcc) {
				t.Errorf("bankService.GetAccountDetails() = %v, want %v", gotAcc, tt.wantAcc)
			}
		})
	}
}
