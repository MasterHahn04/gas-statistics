package storage

type Storage interface {
	Connect() error
	GetIDs() ([]string, error)
	StoreResponse(body string) error
	StorePrices(id string, status string, e5 float32, e10 float32, diesel float32) error
	Close() error
}
