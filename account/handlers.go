package account

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type accountInfoFunc func(m accountInfoManager) (interface{}, error)

func accountAttributeHandler(f accountInfoFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cbm, err := NewAccountInfoManager()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get account information: %v", err)))
			return
		}
		defer cbm.Close()
		b, err := f(cbm)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get account information: %v", err)))
			return
		}
		response, err := json.Marshal(b)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get account information: %v", err)))
			return
		}
		w.Write(response)
	}
	return http.HandlerFunc(fn)
}

func CashBalanceHistoryHandler() http.Handler {
	return accountAttributeHandler(
		func(m accountInfoManager) (interface{}, error) {
			return m.CashBalanceHistory()
		})
}

func AccountValueHistoryHandler() http.Handler {
	return accountAttributeHandler(
		func(m accountInfoManager) (interface{}, error) {
			return m.AccountValueHistory()
		})
}
