package notes

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mtlynch/gofn-prosper/prosper"
)

type mockRedisLRangeReaderByPattern struct {
	KeysPatternGot string
	MatchingKeys   []string
	KeysErr        error
	KeysGot        []string
	LRangeErrs     []error
	lists          map[string][]string
}

func (r *mockRedisLRangeReaderByPattern) Keys(pattern string) ([]string, error) {
	r.KeysPatternGot = pattern
	return r.MatchingKeys, r.KeysErr
}

func (r *mockRedisLRangeReaderByPattern) LRange(key string, start int64, stop int64) ([]string, error) {
	r.KeysGot = append(r.KeysGot, key)
	var err error
	err, r.LRangeErrs = r.LRangeErrs[0], r.LRangeErrs[1:]
	if err != nil {
		return []string{}, err
	}
	list := r.lists[key]
	end := stop + 1
	if end == -1 || end > int64(len(list)) {
		end = int64(len(list))
	}
	return list[start:end], nil
}

func (lr mockRedisLRangeReaderByPattern) Quit() (string, error) {
	return "", nil
}

const (
	noteIdA         = "note:a"
	noteASerialized = `{"Note":{"AgeInMonths":0,"AmountBorrowed":0,"BorrowerRate":0,"DaysPastDue":0,"DebtSaleProceedsReceivedProRataShare":0,"InterestPaidProRataShare":0,"IsSold":false,"LateFeesPaidProRataShare":0,"ListingNumber":0,"LoanNoteID":"a","LoanNumber":0,"NextPaymentDueAmountProRataShare":0,"NextPaymentDueDate":"0001-01-01T00:00:00Z","NoteDefaultReasonDescription":"","NoteDefaultReason":null,"NoteOwnershipAmount":0,"NoteSaleFeesPaid":0,"NoteSaleGrossAmountReceived":0,"NoteStatusDescription":"","NoteStatus":0,"OriginationDate":"0001-01-01T00:00:00Z","PrincipalBalanceProRataShare":0,"PrincipalPaidProRataShare":0,"ProsperFeesPaidProRataShare":0,"Rating":0,"ServiceFeesPaidProRataShare":0,"Term":0},"Timestamp":"2016-03-04T23:19:22.000000022Z"}`
	noteIdB         = "note:b"
	noteBSerialized = `{"Note":{"AgeInMonths":0,"AmountBorrowed":0,"BorrowerRate":0,"DaysPastDue":0,"DebtSaleProceedsReceivedProRataShare":0,"InterestPaidProRataShare":0,"IsSold":false,"LateFeesPaidProRataShare":0,"ListingNumber":0,"LoanNoteID":"b","LoanNumber":0,"NextPaymentDueAmountProRataShare":0,"NextPaymentDueDate":"0001-01-01T00:00:00Z","NoteDefaultReasonDescription":"","NoteDefaultReason":null,"NoteOwnershipAmount":0,"NoteSaleFeesPaid":0,"NoteSaleGrossAmountReceived":0,"NoteStatusDescription":"","NoteStatus":0,"OriginationDate":"0001-01-01T00:00:00Z","PrincipalBalanceProRataShare":0,"PrincipalPaidProRataShare":0,"ProsperFeesPaidProRataShare":0,"Rating":0,"ServiceFeesPaidProRataShare":0,"Term":0},"Timestamp":"2016-03-04T23:19:22.000000022Z"}`
)

var (
	mockKeysErr   = errors.New("mock Keys error")
	mockLRangeErr = errors.New("mock LRange error")
	zeroTime      = time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)

	noteA = prosper.Note{
		LoanNoteID:         "a",
		NextPaymentDueDate: zeroTime,
		OriginationDate:    zeroTime,
	}
	noteB = prosper.Note{
		LoanNoteID:         "b",
		NextPaymentDueDate: zeroTime,
		OriginationDate:    zeroTime,
	}
)

func TestNotes(t *testing.T) {
	var tests = []struct {
		matchingKeys []string
		keysErr      error
		lrangeErrs   []error
		lists        map[string][]string
		wantSuccess  bool
		wantNotes    []prosper.Note
		msg          string
	}{
		{
			keysErr:     mockKeysErr,
			wantSuccess: false,
			msg:         "should return error when Keys call fails",
		},
		{
			matchingKeys: []string{},
			wantSuccess:  false,
			msg:          "should return error when Keys returns empty list",
		},
		{
			matchingKeys: []string{noteIdA},
			lrangeErrs:   []error{mockLRangeErr},
			wantSuccess:  false,
			msg:          "should return error when LRange returns error",
		},
		{
			matchingKeys: []string{noteIdA},
			lrangeErrs:   []error{nil},
			lists: map[string][]string{
				noteIdA: {},
			},
			wantSuccess: false,
			msg:         "should return error when LRange returns empty list",
		},
		{
			matchingKeys: []string{noteIdA},
			lrangeErrs:   []error{nil},
			lists: map[string][]string{
				noteIdA: {noteASerialized},
			},
			wantSuccess: true,
			wantNotes:   []prosper.Note{noteA},
			msg:         "should succeed when keys returns single note",
		},
		{
			matchingKeys: []string{noteIdA, noteIdB},
			lrangeErrs:   []error{nil, nil},
			lists: map[string][]string{
				noteIdA: {noteASerialized},
				noteIdB: {noteBSerialized},
			},
			wantSuccess: true,
			wantNotes:   []prosper.Note{noteA, noteB},
			msg:         "should succeed when keys returns single note",
		},
	}
	for _, tt := range tests {
		nr := noteReader{
			redis: &mockRedisLRangeReaderByPattern{
				MatchingKeys: tt.matchingKeys,
				KeysErr:      tt.keysErr,
				LRangeErrs:   tt.lrangeErrs,
				lists:        tt.lists,
			},
		}
		gotNotes, gotErr := nr.notes()
		if gotErr != nil && tt.wantSuccess {
			t.Errorf("%s: unexpected error from notes, got: %v, want: nil", tt.msg, gotErr)
		}
		if !tt.wantSuccess {
			continue
		}
		if !reflect.DeepEqual(gotNotes, tt.wantNotes) {
			t.Errorf("%s: unexpected records from notes, got: %#v, want: %#v", tt.msg, gotNotes, tt.wantNotes)
		}
	}
}
