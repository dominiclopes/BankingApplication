package bank

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
)

type PingResponse struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) Validate() (err error) {
	// Verify if the email and password is provided
	if r.Email == "" || r.Password == "" {
		return fmt.Errorf("err - email address and password must be provided")
	}
	// Verify if the email is valid
	if _, err = mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("err - invalid email address")
	}

	// Trim the email for any trainling whitespaces
	r.Email = strings.Trim(r.Email, " ")
	return
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateAccountRequest struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func (r *CreateAccountRequest) Validate() (err error) {
	// Verify if the email and password is provided
	if r.Email == "" || r.PhoneNumber == "" {
		return fmt.Errorf("err - email address and password must be provided")
	}

	// Verify if the email is valid
	if _, err = mail.ParseAddress(r.Email); err != nil {
		return fmt.Errorf("err - invalid email address")
	}

	re := regexp.MustCompile(`^\d{10}$`)
	if !re.MatchString(r.PhoneNumber) {
		return fmt.Errorf("err - invalid phone number, must contain 10 digits")
	}

	// Trim the email for any trainling whitespaces
	r.Email = strings.Trim(r.Email, " ")
	return
}

type CreateAccountResponse struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID string `json:"account_id"`
}

type DepositWithdrawAmountRequest struct {
	Amount float32 `json:"amount"`
}

type GetTransactionDetailsRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func (r *GetTransactionDetailsRequest) Validate() (err error) {
	exp_date_format := "2006-01-02"

	// Validate the format of the start date
	startDateTime, err := time.Parse(exp_date_format, r.StartDate)
	if err != nil {
		return fmt.Errorf("err - invalid start date format")
	}

	// Validate the format of the end date
	endDateTime, err := time.Parse(exp_date_format, r.EndDate)
	if err != nil {
		return fmt.Errorf("err - invalid end date format")
	}

	// Verify if the start date is smaller than the end date
	if startDateTime == endDateTime || startDateTime.After(endDateTime) {
		return fmt.Errorf("err - start date must be less than end date")
	}

	// Verify if the difference between the start and end is between 1-30 days
	exp_diff_btw_dates := 30
	diffStartAndEndDate := endDateTime.Sub(startDateTime)
	if diffStartAndEndDate.Hours()/24 > float64(exp_diff_btw_dates) {
		return fmt.Errorf("err - difference between the start date and end date must be less than or equal to %v days", exp_diff_btw_dates)
	}

	return
}
