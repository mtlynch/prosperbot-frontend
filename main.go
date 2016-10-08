package main

import (
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
	log.Println("starting up dashboard")

	http.Handle("/cashBalanceHistory", account.CashBalanceHistoryHandler())
	http.Handle("/accountValueHistory", account.AccountValueHistoryHandler())
	http.Handle("/notes.json", notes.NotesHandler())
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static/"))))
	serveSingle("/", "./static/dashboard.html")
	serveSingle("/notes", "./static/notes.html")

	log.Fatal(http.ListenAndServe(":8082", nil))
}
