package account

import (
	"github.com/mtlynch/prosperbot/redis"
)

func accountValueAttribute(r redis.AccountRecord) float64 {
	return r.Value.TotalAccountValue
}

func (aim accountInfoManager) AccountValueHistory() ([]AccountAttributeRecord, error) {
	return aim.history(accountValueAttribute)
}
