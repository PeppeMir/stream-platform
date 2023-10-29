package dto

import "time"

type CreateMovieRequestDto struct {
	Title         string    `json:"title" validate:"required"`
	Release_Date  time.Time `json:"releaseDate"`
	Genre         string    `json:"genre"`
	Synopsis      string    `json:"synopsis"`
	CastMemberIds []int64   `json:"castMemberIds" validate:"required,min=1"`
}
