package request

type Request interface {
	MakeRequest([]string) error
}
