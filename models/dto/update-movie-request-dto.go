package dto

import "time"

type UpdateMovieRequestDto struct {
	Id            int64     `json:"id" validate:"required"`
	Title         string    `json:"title" validate:"required"`
	Release_Date  time.Time `json:"releaseDate"`
	Genre         string    `json:"genre"`
	Synopsis      string    `json:"synopsis"`
	CastMemberIds []int64   `json:"castMemberIds" validate:"required,min=1"`
}
