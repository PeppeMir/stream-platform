package dto

import (
	"time"
)

type MovieDto struct {
	Id           int64           `json:"id"`
	Title        string          `json:"title" validate:"required"`
	Release_Date time.Time       `json:"releaseDate"`
	Genre        string          `json:"genre"`
	Synopsis     string          `json:"synopsis"`
	CreateUser   UserDto         `json:"createUser"`
	CastMembers  []CastMemberDto `json:"castMembers"`
}
