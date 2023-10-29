package converters

import (
	"stream-platform/models/database"
	"stream-platform/models/dto"
)

func FromCastMemberDbsToCastMemberDtos(dbCastMembers []*database.CastMember) []dto.CastMemberDto {
	var castMembersDtos []dto.CastMemberDto

	for _, dbMember := range dbCastMembers {
		memberDto := dto.CastMemberDto{
			Id:      dbMember.Id,
			Name:    dbMember.Name,
			Surname: dbMember.Surname,
			Age:     dbMember.Age,
		}

		castMembersDtos = append(castMembersDtos, memberDto)
	}

	return castMembersDtos
}
