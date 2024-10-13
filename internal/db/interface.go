package db

type DB interface {
	GetUsers() ([]User, error)
	GetEnableBinanceUsers() ([]User, error)
}
