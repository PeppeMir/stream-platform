package converters

import (
	"stream-platform/models/database"
	"stream-platform/models/dto"
)

func FromFieldsToUserDb(email string, password string) database.User {
	return database.User{Email: email, Password: password}
}

func FromFieldsToUserDto(id int64, email string) dto.UserDto {
	return dto.UserDto{Id: id, Email: email}
}

func FromUserDbToUserDto(dbUser *database.User) dto.UserDto {
	return dto.UserDto{Id: dbUser.Id, Email: dbUser.Email}
}
