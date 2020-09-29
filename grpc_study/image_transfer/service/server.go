package service

type Server interface {
	Listen() (err error)
	Close()
}