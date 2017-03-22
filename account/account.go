package account

import (
	"encoding/json"
	"time"

	"github.com/mtlynch/prosperbot/redis"
)

type (
	redisListReader interface {
		LRange(key string, start int64, stop int64) ([]string, error)
		Quit() (string, error)
	}

	accountInfoManager struct {
		redis redisListReader
	}

	accountAttributeFunc func(r redis.AccountRecord) float64

	// AccountAttributeRecord represents a timestamped value-based property of the
	// user's account (e.g. total value, cash balance).
	AccountAttributeRecord struct {
		Value     float64
		Timestamp time.Time
	}
)

// NewAccountInfoManager creates a new accountInfoManager.
func NewAccountInfoManager() (accountInfoManager, error) {
	r, err := redis.New()
	if err != nil {
		return accountInfoManager{}, err
	}
	return accountInfoManager{r}, nil
}

func accountRecordToAccountAttributeRecord(r redis.AccountRecord, f accountAttributeFunc) AccountAttributeRecord {
	return AccountAttributeRecord{
		Value:     f(r),
		Timestamp: r.Timestamp,
	}
}

func accountRecordsToAccountAttributeRecords(accountRecords []redis.AccountRecord, f accountAttributeFunc) []AccountAttributeRecord {
	cashBalanceRecords := make([]AccountAttributeRecord, len(accountRecords))
	for i, r := range accountRecords {
		cashBalanceRecords[i] = accountRecordToAccountAttributeRecord(r, f)
	}
	return cashBalanceRecords
}

// collapseAccountAttributeRecords collapses the records in a slice to remove
// sequences where the value does not change between adjacent records.
func collapseAccountAttributeRecords(r []AccountAttributeRecord) []AccountAttributeRecord {
	if len(r) <= 1 {
		return r
	}
	prev := r[len(r)-1]
	collapsed := []AccountAttributeRecord{prev}
	for i := len(r) - 2; i >= 0; i-- {
		if prev.Value != r[i].Value {
			collapsed = append(collapsed, r[i])
		}
		prev = r[i]
	}
	return collapsed
}

func (aim accountInfoManager) history(f accountAttributeFunc) ([]AccountAttributeRecord, error) {
	accountRecordsSerialized, err := aim.redis.LRange(redis.KeyAccountInformation, 0, -1)
	if err != nil {
		return []AccountAttributeRecord{}, err
	}
	accountRecords := make([]redis.AccountRecord, len(accountRecordsSerialized))
	for i, v := range accountRecordsSerialized {
		err = json.Unmarshal([]byte(v), &accountRecords[i])
		if err != nil {
			return []AccountAttributeRecord{}, err
		}
	}
	records := accountRecordsToAccountAttributeRecords(accountRecords, f)
	return collapseAccountAttributeRecords(records), nil
}

func (aim accountInfoManager) Close() {
	aim.redis.Quit()
}
