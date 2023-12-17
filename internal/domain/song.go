package domain

import (
	"time"
)

type Song struct {
	Title    string
	Artist   string
	Album    string
	Duration time.Duration
	// SongData stores mp3 file as binary
	SongData []byte
	// CoverData stores png file as binary
	CoverData []byte
}
