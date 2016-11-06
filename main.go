package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/mtlynch/prosperbot-frontend/account"
	"github.com/mtlynch/prosperbot-frontend/notes"
)

func serveSingle(pattern string, filename string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested static file: %v\n", r.URL.Path)
		http.ServeFile(w, r, filename)
	})
}

func main() {
	p := flag.Int("port", 8082, "port on which to listen for web requests")
	flag.Parse()
	log.Printf("starting up dashboard on port %d", *p)

	http.Handle("/cashBalanceHistory", account.CashBalanceHistoryHandler())
	http.Handle("/accountValueHistory", account.AccountValueHistoryHandler())
	http.Handle("/notes.json", notes.NotesHandler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *p), nil))
}
