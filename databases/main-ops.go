package databases

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"stream-platform/customerrors"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	var err error

	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PWD"), os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DB")))
	if err != nil {
		slog.Error(customerrors.ErrCannotConnectToDB.Error(), err)
		panic(err.Error())
	}

	if err = DB.Ping(); err != nil {
		slog.Error(customerrors.ErrCannotConnectToDB.Error(), err)
		panic(err.Error())
	}

	slog.Info("Successfully connected to SQL database!")
}
