package db

const (
	DBNAME = "hotel-reservation"
	DBURI  = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
