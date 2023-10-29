package converters

import (
	"stream-platform/models/database"
	"stream-platform/models/dto"
)

func FromMovieDbToMovieDto(dbMovie *database.Movie, dbUser *database.User, dbCastMembers *[]database.CastMember) *dto.MovieDto {
	movieDto := dto.MovieDto{Id: dbMovie.Id, Title: dbMovie.Title, Release_Date: dbMovie.Release_Date,
		Genre: dbMovie.Genre, Synopsis: dbMovie.Synopsis,
		CreateUser: dto.UserDto{
			Id:    dbUser.Id,
			Email: dbUser.Email,
		},
		CastMembers: []dto.CastMemberDto{}}

	for _, dbMember := range *dbCastMembers {
		memberDto := dto.CastMemberDto{
			Id:      dbMember.Id,
			Name:    dbMember.Name,
			Surname: dbMember.Surname,
			Age:     dbMember.Age,
		}

		movieDto.CastMembers = append(movieDto.CastMembers, memberDto)
	}

	return &movieDto
}
