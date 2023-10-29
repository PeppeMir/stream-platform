package repositories

import (
	"log/slog"
	"stream-platform/databases"
	"stream-platform/models/database"
)

func InsertUser(dbUser *database.User) (int64, error) {
	sql := "INSERT INTO users (Email, Password) VALUES (?, ?);"
	res, err := databases.DB.Exec(sql, dbUser.Email, dbUser.Password)

	if err != nil {
		slog.Error("Error while inserting user in the database", err)
		return -1, err
	}

	return res.LastInsertId()
}

func FindUserByEmail(email string) (*database.User, error) {
	var dbUser database.User
	err := databases.DB.QueryRow("SELECT * FROM users WHERE Email = ?", email).Scan(&dbUser.Id,
		&dbUser.Email, &dbUser.Password, &dbUser.Creation_Date_Time)
	if err != nil {
		slog.Error("Error while finding user by email from the database", err)
		return nil, err
	}

	return &dbUser, err
}
