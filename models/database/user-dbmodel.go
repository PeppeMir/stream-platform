package database

import "time"

type User struct {
	Id                 int64
	Email              string
	Password           string
	Creation_Date_Time time.Time
}
