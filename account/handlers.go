package account

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	accountInfoFunc func(m accountInfoManager) (interface{}, error)

	Handlers struct {
		AccountInfoManager accountInfoManager
	}
)

func NewHandlers() (Handlers, error) {
	m, err := NewAccountInfoManager()
	if err != nil {
		return Handlers{}, err
	}
	return Handlers{AccountInfoManager: m}, nil
}

func (h Handlers) accountAttributeHandler(f accountInfoFunc) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		b, err := f(h.AccountInfoManager)
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

func (h Handlers) CashBalanceHistoryHandler() http.Handler {
	return h.accountAttributeHandler(
		func(m accountInfoManager) (interface{}, error) {
			return m.CashBalanceHistory()
		})
}

func (h Handlers) AccountValueHistoryHandler() http.Handler {
	return h.accountAttributeHandler(
		func(m accountInfoManager) (interface{}, error) {
			return m.AccountValueHistory()
		})
}

func (h *Handlers) Close() {
	h.AccountInfoManager.Close()
}
