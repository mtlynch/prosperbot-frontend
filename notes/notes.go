package notes

import (
	"encoding/json"
	"errors"

	"github.com/mtlynch/gofn-prosper/prosper"
	"github.com/mtlynch/prosperbot/redis"
)

type (
	redisLRangeReaderByPattern interface {
		Keys(pattern string) ([]string, error)
		LRange(key string, start int64, stop int64) ([]string, error)
		Quit() (string, error)
	}

	noteReader struct {
		redis redisLRangeReaderByPattern
	}
)

func NewNoteReader() (noteReader, error) {
	r, err := redis.New()
	if err != nil {
		return noteReader{}, err
	}
	return noteReader{r}, nil
}

var ErrNoNotes = errors.New("no notes found")

func (nr noteReader) notes() ([]prosper.Note, error) {
	noteKeys, err := nr.redis.Keys(redis.KeyPrefixNote + "*")
	if err != nil {
		return []prosper.Note{}, err
	}
	if len(noteKeys) < 1 {
		return []prosper.Note{}, ErrNoNotes
	}
	notes := []prosper.Note{}
	for _, noteKey := range noteKeys {
		noteHistory, err := nr.redis.LRange(noteKey, 0, 0)
		if err != nil {
			return []prosper.Note{}, err
		}
		if len(noteHistory) < 1 {
			return []prosper.Note{}, errors.New("no note history found for " + noteKey)
		}
		noteSerialized := noteHistory[0]
		var noteDeserialized redis.NoteRecord
		err = json.Unmarshal([]byte(noteSerialized), &noteDeserialized)
		if err != nil {
			return []prosper.Note{}, err
		}
		notes = append(notes, noteDeserialized.Note)
	}
	return notes, nil
}

func (nr noteReader) Close() {
	nr.redis.Quit()
}
