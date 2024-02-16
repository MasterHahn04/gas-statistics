package request

type Request interface {
	MakeRequest([]string) (*PricesRespond, string, error)
}
