package repositories

import (
	"log/slog"
	"stream-platform/databases"
	"stream-platform/models/database"
	"strings"
)

func CountCastMembersByIds(ids []int64) (int, error) {
	placeholders, idsIntfc := prepareIdsPlaceholders(ids)

	var count int
	err := databases.DB.QueryRow("SELECT COUNT(*) FROM castmembers WHERE Id IN "+placeholders, idsIntfc...).Scan(&count)
	if err != nil {
		slog.Error("Error while counting cast members in the database", err)
	}

	return count, err
}

func FindAllCastMembers(ids []int64) ([]*database.CastMember, error) {
	placeholders, idsIntfc := prepareIdsPlaceholders(ids)

	rows, err := databases.DB.Query("SELECT * FROM castmembers WHERE Id IN "+placeholders, idsIntfc...)
	if err != nil {
		slog.Error("Error while finding cast members in the database", err)
		return nil, err
	}

	defer rows.Close()

	var results []*database.CastMember

	for rows.Next() {
		var dbCastMember database.CastMember
		err := rows.Scan(&dbCastMember.Id,
			&dbCastMember.Name, &dbCastMember.Surname, &dbCastMember.Age)
		if err != nil {
			slog.Error("Error while scanning cast members SELECT result")
			return nil, err
		}

		results = append(results, &dbCastMember)
	}

	return results, nil
}

func prepareIdsPlaceholders(ids []int64) (string, []interface{}) {
	placeholders := "(?" + strings.Repeat(",?", len(ids)-1) + ")"

	idsIntfc := make([]interface{}, len(ids))
	for i := range ids {
		idsIntfc[i] = ids[i]
	}

	return placeholders, idsIntfc
}
