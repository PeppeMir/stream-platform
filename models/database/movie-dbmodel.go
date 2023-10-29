package database

import "time"

type Movie struct {
	Id            int64
	Title         string
	Release_Date  time.Time
	Genre         string
	Synopsis      string
	CreateUser_Id int64
}
