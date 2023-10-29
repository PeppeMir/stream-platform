package dto

type CastMemberDto struct {
	Id      int64  `json:"id"`
	Name    string `json:"name" validate:"required"`
	Surname string `json:"surname" validate:"required"`
	Age     int16  `json:"age"`
}
