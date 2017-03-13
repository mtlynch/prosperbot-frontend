package account

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

const (
	accountInformationSerializedA = `{"Value":{"AvailableCashBalance":100,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue":0,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-01-28T15:35:04.000000022Z"}`
	accountInformationSerializedB = `{"Value":{"AvailableCashBalance":125.5,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue":0,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:28:15.000000022Z"}`
	accountInformationSerializedC = `{"Value":{"AvailableCashBalance":125.5,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue":0,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:29:15.000000022Z"}`
	accountInformationSerializedD = `{"Value":{"AvailableCashBalance":95.25,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue":0,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:30:15.000000022Z"}`
	badJSON                       = "{{mock bad JSON"
)

func TestCashBalanceHistory(t *testing.T) {
	var tests = []struct {
		accountList []string
		lrangeErr   error
		wantRecords []AccountAttributeRecord
		wantSuccess bool
		msg         string
	}{
		{
			accountList: []string{
				accountInformationSerializedB,
				accountInformationSerializedA,
			},
			lrangeErr:   errors.New("mock LRange error"),
			wantSuccess: false,
			msg:         "should return error when redis call fails",
		},
		{
			accountList: []string{},
			wantRecords: []AccountAttributeRecord{},
			wantSuccess: true,
			msg:         "should return empty list when there is no account history",
		},
		{
			accountList: []string{
				accountInformationSerializedB,
			},
			wantRecords: []AccountAttributeRecord{
				{
					Value:     125.50,
					Timestamp: time.Date(2016, 2, 14, 12, 28, 15, 22, time.UTC),
				},
			},
			wantSuccess: true,
			msg:         "should return valid record when single valid record exists",
		},
		{
			accountList: []string{
				accountInformationSerializedD,
				accountInformationSerializedC,
				accountInformationSerializedB,
				accountInformationSerializedA,
			},
			wantRecords: []AccountAttributeRecord{
				{
					Value:     100.0,
					Timestamp: time.Date(2016, 1, 28, 15, 35, 4, 22, time.UTC),
				},
				{
					Value:     125.50,
					Timestamp: time.Date(2016, 2, 14, 12, 28, 15, 22, time.UTC),
				},
				{
					Value:     95.25,
					Timestamp: time.Date(2016, 2, 14, 12, 30, 15, 22, time.UTC),
				},
			},
			wantSuccess: true,
			msg:         "should return valid record when valid records exist",
		},
		{
			accountList: []string{badJSON},
			wantSuccess: false,
			msg:         "malformed JSON in redis should cause error",
		},
	}
	for _, tt := range tests {
		aim := accountInfoManager{
			redis: &mockRedisListReader{
				List: tt.accountList,
			},
		}
		gotRecords, gotErr := aim.CashBalanceHistory()
		if gotErr != nil && tt.wantSuccess {
			t.Errorf("%s: unexpected error from CashBalanceHistory, got: %v, want: nil", tt.msg, gotErr)
		}
		if !tt.wantSuccess {
			continue
		}
		if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
			t.Errorf("%s: unexpected records from CashBalanceHistory, got: %v, want: %v", tt.msg, gotRecords, tt.wantRecords)
		}
	}
}
