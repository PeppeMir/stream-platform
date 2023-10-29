package repositories

import (
	"errors"
	"log/slog"
	"stream-platform/customerrors"
	"stream-platform/databases"
	"stream-platform/models/database"
	"stream-platform/models/filter"
	"strings"
)

func InsertMovie(dbMovie *database.Movie, castMemberIds *[]int64) (int64, error) {
	// Begin a new transaction
	tx, txErr := databases.DB.Begin()
	if txErr != nil {
		slog.Error("Error while creating transaction context", txErr)
		return -1, errors.New("")
	}

	// Insert movie
	sqlMovies := "INSERT INTO movies (Title, Release_Date, Genre, Synopsis, CreateUser_Id) VALUES (?, ?, ?, ?, ?);"
	res, insertErr := tx.Exec(sqlMovies, dbMovie.Title, dbMovie.Release_Date, dbMovie.Genre, dbMovie.Synopsis, dbMovie.CreateUser_Id)
	if insertErr != nil {
		slog.Error("Error while inserting the movie", insertErr)
		tx.Rollback()
		return -1, insertErr
	}

	// Stick with just inserted movie identifier
	lastInsertId, _ := res.LastInsertId()

	// Insert all the (movieId, castMemberId) pairs
	moviesCastMembersSql := "INSERT INTO movies_castmembers (Movie_Id, CastMember_Id) VALUES (?, ?);"
	for _, castMemberId := range *castMemberIds {
		_, err := tx.Exec(moviesCastMembersSql, lastInsertId, castMemberId)
		if err != nil {
			slog.Error("Error while inserting a movie-castmember record", "movieId", lastInsertId, "castMemberId", castMemberId)
			tx.Rollback()
			return -1, err
		}
	}

	// Commit transaction
	txCommitErr := tx.Commit()
	if txCommitErr != nil {
		slog.Error("Error while commiting the create movie transaction", txCommitErr)
		return -1, txCommitErr
	}

	return res.LastInsertId()
}

func UpdateMovie(dbMovie *database.Movie, castMemberIdsToCreate *[]int64, castMemberIdsToDelete *[]int64) error {
	// Begin a new transaction
	tx, txErr := databases.DB.Begin()
	if txErr != nil {
		slog.Error("Error while creating transaction context", txErr)
		return errors.New("")
	}

	// Insert movie
	updateMovieSql := "UPDATE movies SET Title=?, Release_Date=?, Genre=?, Synopsis=? WHERE Id = ?;"
	_, updateErr := tx.Exec(updateMovieSql, dbMovie.Title, dbMovie.Release_Date, dbMovie.Genre, dbMovie.Synopsis, dbMovie.Id)
	if updateErr != nil {
		slog.Error("Error while updating the movie", updateErr)
		tx.Rollback()
		return updateErr
	}

	// Insert all the new (movieId, castMemberId) pairs
	insertCastMembersSql := "INSERT INTO movies_castmembers (Movie_Id, CastMember_Id) VALUES (?, ?);"
	for _, castMemberId := range *castMemberIdsToCreate {
		_, err := tx.Exec(insertCastMembersSql, dbMovie.Id, castMemberId)
		if err != nil {
			slog.Error("Error while inserting a movie-castmember record", "movieId", dbMovie.Id, "castMemberId", castMemberId)
			tx.Rollback()
			return err
		}
	}

	// Delete all the obsolete (movieId, castMemberId) pairs
	deleteCastMemberSql := "DELETE FROM movies_castmembers WHERE Movie_Id=? AND CastMember_Id=?;"
	for _, castMemberId := range *castMemberIdsToDelete {
		_, err := tx.Exec(deleteCastMemberSql, dbMovie.Id, castMemberId)
		if err != nil {
			slog.Error("Error while deleting a movie-castmember record", "movieId", dbMovie.Id, "castMemberId", castMemberId)
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	txCommitErr := tx.Commit()
	if txCommitErr != nil {
		slog.Error("Error while commiting the update movie transaction", txCommitErr)
		return txCommitErr
	}

	return nil
}

func GetMovie(id int64) (*database.Movie, *database.User, *[]database.CastMember, error) {
	filters := filter.SearchMoviesFilter{Id: id}
	tuples, err := GetMovies(&filters)
	if err != nil {
		return nil, nil, nil, err
	}

	tuple, isPresent := (*tuples)[id]
	if !isPresent {
		return nil, nil, nil, customerrors.ErrMovieNotFound
	}

	return &tuple.Movie, &tuple.User, &tuple.CastMembers, nil
}

func GetMovies(filters *filter.SearchMoviesFilter) (*map[int64]*database.SearchMovieTuple, error) {
	sql := `SELECT m.Id, m.Title, m.Release_Date, m.Genre, m.Synopsis, m.CreateUser_Id,
				   u.Id, u.Email,
				   c.Id, c.Name, c.Surname, c.Age 
			FROM movies m 
				INNER JOIN users u 
					ON m.CreateUser_Id = u.Id
				INNER JOIN movies_castmembers mc
					ON mc.Movie_Id = m.Id
				INNER JOIN castmembers c
					ON c.Id = mc.CastMember_Id`

	sql, params := buildGetMoviesWhere(filters, sql)

	rows, err := databases.DB.Query(sql, params...)
	if err != nil {
		slog.Error("Error while retrieving movies from the database", err)
		return nil, err
	}

	defer rows.Close()

	results := make(map[int64]*database.SearchMovieTuple)

	for rows.Next() {
		var dbMovie = database.Movie{}
		var dbUser = database.User{}
		var dbCastMember database.CastMember
		err := rows.Scan(&dbMovie.Id, &dbMovie.Title, &dbMovie.Release_Date,
			&dbMovie.Genre, &dbMovie.Synopsis, &dbMovie.CreateUser_Id,
			&dbUser.Id, &dbUser.Email,
			&dbCastMember.Id, &dbCastMember.Name, &dbCastMember.Surname, &dbCastMember.Age)
		if err != nil {
			slog.Error("Error while retrieving cast members from the database")
			return nil, err
		}

		tuple, isPresent := results[dbMovie.Id]
		if !isPresent {
			results[dbMovie.Id] = &database.SearchMovieTuple{
				Movie:       dbMovie,
				User:        dbUser,
				CastMembers: []database.CastMember{},
			}

			tuple = results[dbMovie.Id]
		}

		tuple.CastMembers = append(tuple.CastMembers, dbCastMember)
	}

	return &results, nil
}

func DeleteMovie(id int64, userId int64) (int64, error) {
	// Begin a new transaction
	tx, txErr := databases.DB.Begin()
	if txErr != nil {
		slog.Error("Error while creating transaction context", txErr)
		return -1, errors.New("")
	}

	mainSql := "DELETE FROM movies WHERE Id = ? AND CreateUser_Id = ?;"
	res, mainDeleteErr := tx.Exec(mainSql, id, userId)

	if mainDeleteErr != nil {
		slog.Error("Error while deleting movie from the database", "id", id, mainDeleteErr)
		tx.Rollback()
		return -1, mainDeleteErr
	}

	numRows, _ := res.RowsAffected()
	if numRows > 0 {
		sql := "DELETE FROM movies_castmembers WHERE Movie_Id = ?;"
		_, deleteErr := tx.Exec(sql, id)

		if deleteErr != nil {
			slog.Error("Error while deleting movie cast-members from the database", "movieId", id, deleteErr)
			tx.Rollback()
			return -1, deleteErr
		}
	}

	// Commit transaction
	txCommitErr := tx.Commit()
	if txCommitErr != nil {
		slog.Error("Error while commiting the delete movie transaction", txCommitErr)
		return -1, txCommitErr
	}

	return res.RowsAffected()
}

func buildGetMoviesWhere(filters *filter.SearchMoviesFilter, sql string) (string, []any) {
	if filters.Id != 0 || filters.Title != "" || filters.Genre != "" || !filters.ReleaseDate.IsZero() {
		sql = sql + " WHERE "
	}

	var params []any

	if filters.Id != 0 {
		sql = sql + "m.Id = ? AND "
		params = append(params, filters.Id)
	}

	if filters.Title != "" {
		sql = sql + "m.Title LIKE ? AND "
		params = append(params, "%"+filters.Title+"%")
	}

	if filters.Genre != "" {
		sql = sql + "m.Genre LIKE ? AND "
		params = append(params, "%"+filters.Genre+"%")
	}

	if !filters.ReleaseDate.IsZero() {
		sql = sql + "m.Release_Date = ?"
		params = append(params, filters.ReleaseDate)
	}

	if strings.HasSuffix(sql, "AND ") {
		idx := strings.LastIndex(sql, "AND")
		sql = sql[:idx] + strings.Replace(sql[idx:], "AND", "", 1)
	}
	return sql, params
}
