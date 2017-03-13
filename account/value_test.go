package account

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

const (
	accountInformationSerializedE = `{"Value":{"AvailableCashBalance":0,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue": 34502.21,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-01-28T15:35:04.000000022Z"}`
	accountInformationSerializedF = `{"Value":{"AvailableCashBalance":0,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue": 34568.43,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:28:15.000000022Z"}`
	accountInformationSerializedG = `{"Value":{"AvailableCashBalance":0,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue": 34568.43,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:29:15.000000022Z"}`
	accountInformationSerializedH = `{"Value":{"AvailableCashBalance":0,"TotalPrincipalReceivedOnActiveNotes":0,"OutstandingPrincipalOnActiveNotes":0,"LastWithdrawAmount":0,"LastDepositAmount":0,"LastDepositDate":"0001-01-01T00:00:00Z","PendingInvestmentsPrimaryMarket":0,"PendingInvestmentsSecondaryMarket":0,"PendingQuickInvestOrders":0,"TotalAmountInvestedOnActiveNotes":0,"TotalAccountValue": 34575.98,"InflightGross":0,"LastWithdrawDate":"0001-01-01T00:00:00Z"},"Timestamp":"2016-02-14T12:30:15.000000022Z"}`
)

func TestAccountValueHistory(t *testing.T) {
	var tests = []struct {
		accountList []string
		lrangeErr   error
		wantRecords []AccountAttributeRecord
		wantSuccess bool
		msg         string
	}{
		{
			accountList: []string{
				accountInformationSerializedF,
				accountInformationSerializedE,
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
				accountInformationSerializedF,
			},
			wantRecords: []AccountAttributeRecord{
				{
					Value:     34568.43,
					Timestamp: time.Date(2016, 2, 14, 12, 28, 15, 22, time.UTC),
				},
			},
			wantSuccess: true,
			msg:         "should return valid record when single valid record exists",
		},
		{
			accountList: []string{
				accountInformationSerializedH,
				accountInformationSerializedG,
				accountInformationSerializedF,
				accountInformationSerializedE,
			},
			wantRecords: []AccountAttributeRecord{
				{
					Value:     34502.21,
					Timestamp: time.Date(2016, 1, 28, 15, 35, 4, 22, time.UTC),
				},
				{
					Value:     34568.43,
					Timestamp: time.Date(2016, 2, 14, 12, 28, 15, 22, time.UTC),
				},
				{
					Value:     34575.98,
					Timestamp: time.Date(2016, 2, 14, 12, 30, 15, 22, time.UTC),
				},
			},
			wantSuccess: true,
			msg:         "should return valid records when valid records exist",
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
		gotRecords, gotErr := aim.AccountValueHistory()
		if gotErr != nil && tt.wantSuccess {
			t.Errorf("%s: unexpected error from AccountValueHistory, got: %v, want: nil", tt.msg, gotErr)
		}
		if !tt.wantSuccess {
			continue
		}
		if !reflect.DeepEqual(gotRecords, tt.wantRecords) {
			t.Errorf("%s: unexpected records from AccountValueHistory, got: %v, want: %v", tt.msg, gotRecords, tt.wantRecords)
		}
	}
}
