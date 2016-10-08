package account

import (
	"github.com/mtlynch/prosperbot/redis"
)

func cashBalanceAttribute(r redis.AccountRecord) float64 {
	return r.Value.AvailableCashBalance
}

func (aim accountInfoManager) CashBalanceHistory() ([]AccountAttributeRecord, error) {
	return aim.history(cashBalanceAttribute)
}
