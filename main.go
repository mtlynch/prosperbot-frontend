package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/prosperbot-frontend/account"
	"github.com/mtlynch/prosperbot-frontend/notes"
)

func main() {
	p := flag.Int("port", 8082, "port on which to listen for web requests")
	flag.Parse()
	log.Printf("starting up dashboard on port %d", *p)

	h, err := account.NewHandlers()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize handlers: %v", err))
	}
	defer h.Close()
	http.Handle("/cashBalanceHistory", h.CashBalanceHistoryHandler())
	http.Handle("/accountValueHistory", h.AccountValueHistoryHandler())
	http.Handle("/notes.json", notes.NotesHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *p), nil))
}
