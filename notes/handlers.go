package notes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func NotesHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		reader, err := NewNoteReader()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get note information: %v", err)))
			return
		}
		defer reader.Close()
		n, err := reader.notes()
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get note information: %v", err)))
			return
		}
		response, err := json.Marshal(n)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("failed to get note information: %v", err)))
			return
		}
		w.Write(response)
	}
	return http.HandlerFunc(fn)
}
