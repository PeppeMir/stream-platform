package filter

import "time"

type SearchMoviesFilter struct {
	Id          int64
	Title       string
	ReleaseDate time.Time
	Genre       string
}
