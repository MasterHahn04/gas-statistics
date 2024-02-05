package storage

type Storage interface {
	Connect() error
	GetIDs() ([]string, error)
}
