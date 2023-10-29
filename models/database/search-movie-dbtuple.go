package database

type SearchMovieTuple struct {
	Movie       Movie
	User        User
	CastMembers []CastMember
}
